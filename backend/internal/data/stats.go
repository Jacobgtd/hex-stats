package data

import "github.com/rs/zerolog"

type DataClient struct {
	clients DataClientClients
	logger  zerolog.Logger
}

func NewDataClient(logger zerolog.Logger, clients DataClientClients) *DataClient {
	return &DataClient{
		clients: clients,
		logger:  logger,
	}
}
