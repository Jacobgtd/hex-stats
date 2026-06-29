package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/interal/monitoring"
	"github.com/Jacobgtd/hex-stats/backend/internal/common"

	"github.com/rs/zerolog"
)

type GithubUser struct {
	Login string `json:"login"`
}

type GithubClient struct {
	logger zerolog.Logger
	config *GithubClientConfig
}

func NewGithubClient(logger zerolog.Logger, config *GithubClientConfig) *GithubClient {
	return &GithubClient{
		config: config,
		logger: logger,
	}
}

// IsAdmin checks if the GitHub user associated with the token is an admin.
// Returns the username and nil if user is admin, empty string and StatusError if not or on failure.
func (g *GithubClient) IsAdmin(ctx context.Context, ghtoken string) (string, *common.StatusError) {
	g.logger.Debug().Msg("checking if GitHub user is admin")

	var user GithubUser
	request := monitoring.NewHttpRequest(g.config.Url, "user", http.MethodGet)
	statusCode, err := request.WithCtx(ctx).WithBearerToken(ghtoken).WithTimeout(time.Second * 5).WithExpectedFailureCode().Do(&user)
	if err != nil {
		return "", &common.StatusError{
			Code:  statusCode,
			Error: err,
		}
	}

	g.logger.Debug().Str("login", user.Login).Msg("fetched GitHub user")

	for _, admin := range g.config.Admins {
		if user.Login == admin {
			g.logger.Info().Str("login", user.Login).Msg("user is admin")
			return user.Login, nil
		}
	}

	g.logger.Warn().Str("login", user.Login).Msg("user is not admin")
	return "", &common.StatusError{Code: http.StatusForbidden, Error: fmt.Errorf("user is not admin")}
}
