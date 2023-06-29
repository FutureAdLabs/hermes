package hermes

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

type CloudWatchWriter struct {
	logGroupName  string
	logStreamName string
	service       *cloudwatchlogs.CloudWatchLogs
}

// Generate the current timestamp as an int64
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (cw *CloudWatchWriter) Write(p []byte) (n int, err error) {
	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(cw.logGroupName),
		LogStreamName: aws.String(cw.logStreamName),
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   aws.String(string(p)),
				Timestamp: aws.Int64(makeTimestamp()),
			},
		},
	}

	_, err = cw.service.PutLogEvents(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == cloudwatchlogs.ErrCodeResourceNotFoundException {
				fmt.Println("A CloudWatch LogGroup does not exist with this name. Verify that it exists or that you are not initialising the Logger with the wrong name")
				panic(err.Error())
			} else {
				panic(err.Error())
			}
		} else {
			fmt.Println("Failed to cast error to awserr.Error:", err)
		}
		fmt.Println(err.Error())
		return 0, err
	}

	return len(p), nil
}

// Initialise the Hermes and its logging components.
// svcName is used as the logGroupName.
func Init(svcName string) {

	fmt.Println("Init logging")
	env := os.Getenv("ENV")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if env == "dev" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1"),
		})
		if err != nil {
			panic(err)
		}

		logGroup := fmt.Sprintf("%s-%s", svcName, env)
		podName := os.Getenv("HOSTNAME")
		service := cloudwatchlogs.New(sess)

		// Create the Log Stream if it doesn't exist
		_, err = service.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroup),
			LogStreamName: aws.String(podName),
		})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Ignore error if the Log Stream already exists
				if awsErr.Code() != cloudwatchlogs.ErrCodeResourceAlreadyExistsException {
					panic(err)
				}
			} else {
				panic(err)
			}
		}

		writer := &CloudWatchWriter{
			logGroupName:  logGroup,
			logStreamName: podName,
			service:       service,
		}

		logger = zerolog.New(writer).With().Timestamp().Logger()
	}
}

// Public function to retrieve the initialised logger
func Logger() zerolog.Logger {
	return logger
}
