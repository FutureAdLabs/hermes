# Hermes

Hermes is Utopia's main logging utility. It handles local environments by logging to standard output or AWS CloudWatch Logs for our deployed environments.

## How to use

First import and initialise Hermes.

```go
import (
  "github.com/FutureAdLabs/hermes"
)

...

func main() {
...
  // Initialise it with your service name
  hermes.Init("sso")
...
}
```

Once done, you can get a logger from Hermes anywhere!


```go
import (
  "github.com/FutureAdLabs/hermes"
)
func myFunc() {
  logger := hermes.Logger()

  logger.Info().Msg("beep boop this is a log")
}
```

The logging library Hermes uses is zerolog. To know more about how to use it https://github.com/rs/zerolog.

## Troubleshooting
- Hermes does NOT create LogGroups on CloudWatch (the LogGroup name to use is selected as ${YOUR_SERVICE}-${ENVIRONMENT} (e.g. sso-testing).
- Hermes DOES create LogStreams inside the LogGroup. It uses the HOSTNAME env variable to decide the stream name. Each Pod in Kubernetes has its own hostname.
- Hermes differentiates between environments based on the ENV environment variable. For local development ENV=dev, otherwise it should match the deployed env as ENV=testing or ENV=production
