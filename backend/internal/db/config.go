package db

import (
	"fmt"

	"github.com/Jacobgtd/hex-stats/backend/internal/configpack"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSL      bool
}

func LoadDBConfig() (*DBConfig, error) {

	err := configpack.Load("db.config")
	if err != nil {
		return nil, err
	}

	host, err := configpack.String("HOST")
	if err != nil {
		return nil, err
	}
	port, err := configpack.Int("PORT")
	if err != nil {
		return nil, err
	}
	user, err := configpack.String("USER")
	if err != nil {
		return nil, err
	}
	password, err := configpack.String("PASSWORD")
	if err != nil {
		return nil, err
	}
	dbName, err := configpack.String("NAME")
	if err != nil {
		return nil, err
	}
	ssl, err := configpack.Bool("SSL")
	if err != nil {
		return nil, err
	}

	return &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     dbName,
		SSL:      ssl,
	}, nil
}

func (c *DBConfig) getConnStr() string {
	sslMode := "disable"
	if c.SSL {
		sslMode = "require"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Name, sslMode)
}
