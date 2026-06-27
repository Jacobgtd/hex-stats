package main

import (
	"os"

	"github.com/Jacobgtd/hex-stats/backend/internal/authn"
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

	authnConfig, err := authn.LoadAuthnConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load authn config")
	}
	authnClient := authn.NewAuthnClient(logger, authnConfig)

	clients := &server.ServerClients{
		GithubClient: ghClient,
		AuthnClient:  authnClient,
	}

	srv := server.NewServer(logger, config, clients)
	if err := srv.Run(); err != nil {
		logger.Fatal().Err(err).Msg("failed to run server")
	}

}
