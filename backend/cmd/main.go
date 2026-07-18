package main

import (
	"net/http"
	"os"

	"github.com/Jacobgtd/hex-stats/backend/internal/alphavantage"
	"github.com/Jacobgtd/hex-stats/backend/internal/auth"
	"github.com/Jacobgtd/hex-stats/backend/internal/cache"
	"github.com/Jacobgtd/hex-stats/backend/internal/clients"
	"github.com/Jacobgtd/hex-stats/backend/internal/db"
	"github.com/Jacobgtd/hex-stats/backend/internal/github"
	"github.com/Jacobgtd/hex-stats/backend/internal/monitoring"
	"github.com/Jacobgtd/hex-stats/backend/internal/server"
	"github.com/rs/zerolog"
)

func main() {

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Caller().Logger()

	config, err := server.LoadServerConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load server config")
	}

	// Initialize Cache
	cacheConfig, err := cache.LoadLocalCacheConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load cache config")
	}
	cache := cache.NewLocalCache(cacheConfig)

	// Initialize HTTP client
	httpClient := monitoring.NewHTTPClient(http.Client{}, cache)

	// Initialize Alphavantage clients
	avConfig, err := alphavantage.LoadAlphavantageClientConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load Alphavantage config")
	}
	avClient := alphavantage.NewAlphavantageClient(avConfig, httpClient)

	// Initialize GithubClient
	ghConfig, err := github.LoadGithubClientConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load GitHub client config")
	}

	ghClient := github.NewGithubClient(logger, ghConfig, httpClient)

	// Initialize AuthClient
	authConfig, err := auth.LoadAuthConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load auth config")
	}
	authClient := auth.NewAuthClient(logger, authConfig)

	// Initialize DBClient
	dbConfig, err := db.LoadDBConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load DB config")
	}

	dbClient, err := db.NewDBClient(logger, dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize DB client")
	}

	clients := &clients.Clients{
		GithubClient:       ghClient,
		AuthClient:         authClient,
		DBClient:           dbClient,
		AlphavantageClient: avClient,
	}

	srv := server.NewServer(logger, config, clients)
	if err := srv.Run(); err != nil {
		logger.Fatal().Err(err).Msg("failed to run server")
	}

}
