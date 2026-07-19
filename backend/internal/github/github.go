package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/common"
	"github.com/Jacobgtd/hex-stats/backend/internal/monitoring"

	"github.com/rs/zerolog"
)

type GithubUser struct {
	Login string `json:"login"`
}

type GithubClient struct {
	logger     zerolog.Logger
	config     *GithubClientConfig
	httpClient *monitoring.HTTPClient
}

func NewGithubClient(logger zerolog.Logger, config *GithubClientConfig, httpClient *monitoring.HTTPClient) *GithubClient {
	return &GithubClient{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

// IsAdmin checks if the GitHub user associated with the token is an admin.
// Returns the username and nil if user is admin, empty string and StatusError if not or on failure.
func (g *GithubClient) IsAdmin(ctx context.Context, ghtoken string) (string, *common.StatusError) {
	g.logger.Debug().Msg("checking if GitHub user is admin")

	var user GithubUser
	resp, err := g.httpClient.NewHTTPRequest(g.config.Url, "user", http.MethodGet).
		WithBearerToken(ghtoken).
		WithTimeout(5*time.Second).
		WithPossibleResponseCodes(http.StatusOK, http.StatusUnauthorized, http.StatusForbidden).
		WithCacheTTL(24 * time.Hour).
		Do(ctx)

	if err != nil {
		return "", &common.StatusError{
			Code:  http.StatusFailedDependency,
			Error: err,
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &common.StatusError{
			Code:  http.StatusUnauthorized,
			Error: fmt.Errorf("unauthorized"),
		}
	}

	if err := resp.Unmarshal(&user); err != nil {
		return "", &common.StatusError{
			Code:  http.StatusUnauthorized,
			Error: fmt.Errorf("failed to unmarshal user data"),
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
