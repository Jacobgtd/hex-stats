package server

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/Jacobgtd/hex-stats/backend/internal/ca"
	"github.com/Jacobgtd/hex-stats/backend/internal/db"
	"github.com/Jacobgtd/hex-stats/backend/internal/github"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ServerClients struct {
	GithubClient *github.GithubClient
	CAClient     *ca.CAClient
	DBClient     *db.DBClient
}

type Server struct {
	engine  *gin.Engine
	logger  zerolog.Logger
	config  *ServerConfig
	clients *ServerClients
}

func recoveryMiddleware(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Any("panic", err).
					Str("path", c.Request.URL.Path).
					Msg("panic recovered")

				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}

func NewServer(logger zerolog.Logger, config *ServerConfig, clients *ServerClients) *Server {

	e := gin.New()
	e.Use(loggerMiddleware(logger))
	e.Use(recoveryMiddleware(logger))

	e.GET("/health", health)

	caAuthGroup := e.Group("/api/v1")
	caAuthGroup.Use(caAuthMiddleware(logger, clients.CAClient))
	adminAuthGroup := e.Group("/api/v1")
	adminAuthGroup.Use(adminAuthMiddleware(logger, clients.GithubClient))
	noAuthGroup := e.Group("/api/v1")

	server := &Server{
		logger:  logger,
		engine:  e,
		config:  config,
		clients: clients,
	}

	noAuthGroup.POST("/devices/:deviceId/certificate", server.generateDeviceCertificate)
	caAuthGroup.GET("/devices/:deviceId/certificate/verify", server.verifyDeviceAuthHandler)
	adminAuthGroup.POST("/devices", server.newDeviceHandler)

	return server
}

func (s *Server) Run() error {
	s.logger.Info().Uint("port", s.config.port).Msg("starting server")

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(int(s.config.port)),
		Handler: s.engine,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
			ClientCAs:  nil,
		},
	}
	return srv.ListenAndServeTLS(s.config.crtPath, s.config.keyPath)
}
