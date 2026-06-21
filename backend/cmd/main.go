package main

import (
	"os"

	"github.com/Jacobgtd/hex-stats/backend/internal/ca"
	"github.com/Jacobgtd/hex-stats/backend/internal/github"
	"github.com/Jacobgtd/hex-stats/backend/internal/server"
	"github.com/rs/zerolog"
)

func main() {

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Caller().Logger()

	config, err := server.LoadServerConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load server config")
	}

	// Initialize GithubClient
	ghConfig, err := github.LoadGithubClientConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load GitHub client config")
	}

	ghClient := github.NewGithubClient(logger, ghConfig)

	//Initialize CAClient
	caConfig, err := ca.LoadCAConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load CA config")
	}

	caClient := ca.NewCAClient(logger, caConfig)

	clients := &server.ServerClients{
		GithubClient: ghClient,
		CAClient:     caClient,
	}

	srv := server.NewServer(logger, config, clients)
	if err := srv.Run(); err != nil {
		logger.Fatal().Err(err).Msg("failed to run server")
	}

}
