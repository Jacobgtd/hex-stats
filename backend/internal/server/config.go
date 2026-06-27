package server

import "github.com/Jacobgtd/hex-stats/backend/internal/configpack"

type ServerConfig struct {
	port    uint
	crt     string
	crtPath string
	key     string
	keyPath string
}

func LoadServerConfig() (*ServerConfig, error) {
	port := configpack.IntOrDefault("PORT", 8080)
	crt, err := configpack.LoadFile("server.crt")
	if err != nil {
		return nil, err
	}
	crtPath, err := configpack.GetPath("server.crt")
	if err != nil {
		return nil, err
	}
	key, err := configpack.LoadFile("server.key")
	if err != nil {
		return nil, err
	}
	keyPath, err := configpack.GetPath("server.key")
	if err != nil {
		return nil, err
	}
	return &ServerConfig{
		port:    uint(port),
		crt:     crt,
		crtPath: crtPath,
		key:     key,
		keyPath: keyPath,
	}, nil
}
