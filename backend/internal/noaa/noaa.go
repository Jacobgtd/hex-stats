package noaa

import "github.com/rs/zerolog"

type NoaaClient struct {
	config NoaaClientConfig
	logger zerolog.Logger
}

func NewNoaaClient(logger zerolog.Logger, config NoaaClientConfig) NoaaClient {
	return NoaaClient{
		config: config,
		logger: logger,
	}
}
