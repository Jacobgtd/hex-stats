package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/rs/zerolog"
)

type DBClient struct {
	config *DBConfig
	logger zerolog.Logger
	db     *sql.DB
}

func NewDBClient(logger zerolog.Logger, config *DBConfig) (*DBClient, error) {

	db, err := sql.Open("pgx", config.getConnStr())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open database connection")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Failed to ping database")
		return nil, err
	}

	return &DBClient{
		config: config,
		logger: logger,
		db:     db,
	}, nil
}

func (d *DBClient) Close() error {
	return d.db.Close()
}
