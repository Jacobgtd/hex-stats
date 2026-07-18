package server

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/Jacobgtd/hex-stats/backend/internal/auth"
	"github.com/Jacobgtd/hex-stats/backend/internal/clients"
	"github.com/Jacobgtd/hex-stats/backend/internal/metricsapi"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	engine  *gin.Engine
	logger  zerolog.Logger
	config  *ServerConfig
	clients *clients.Clients
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

func NewServer(logger zerolog.Logger, config *ServerConfig, clients *clients.Clients) *Server {

	e := gin.New()
	e.Use(loggerMiddleware(logger))
	e.Use(recoveryMiddleware(logger))

	server := &Server{
		logger:  logger,
		engine:  e,
		config:  config,
		clients: clients,
	}

	e.GET("/health", Health)

	apiGroup := e.Group("/api/v1")

	apiGroup.POST("/auth/github", server.authGithub)
	apiGroup.POST("/auth/device", server.authDevice)
	apiGroup.GET("/auth", server.authMiddleware(logger, auth.PermissionsDefault), server.checkAuth)
	apiGroup.POST("/device", server.authMiddleware(logger, auth.PermissionsAdmin), server.newDeviceHandler)

	metricsApiGroup := apiGroup.Group("/metrics")
	metricsApiGroup.Use(server.authMiddleware(logger, auth.PermissionsDefault))
	metricsapi.RegisterMetricsRoutes(metricsApiGroup, clients)

	return server
}

func (s *Server) Run() error {
	s.logger.Info().Uint("port", s.config.port).Msg("starting server")

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(int(s.config.port)),
		Handler: s.engine,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}
	return srv.ListenAndServeTLS(s.config.crtPath, s.config.keyPath)
}
