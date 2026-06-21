package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	// Create request to GitHub API
	req, err := http.NewRequestWithContext(ctx, "GET", g.config.Url+"/user", nil)
	if err != nil {
		g.logger.Error().Err(err).Msg("failed to create GitHub API request")
		return "", &common.StatusError{Code: http.StatusInternalServerError, Error: err}
	}

	// Add bearer token header
	req.Header.Set("Authorization", "Bearer "+ghtoken)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.logger.Error().Err(err).Msg("failed to fetch GitHub user")
		return "", &common.StatusError{Code: http.StatusInternalServerError, Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
		g.logger.Error().Err(err).Int("status", resp.StatusCode).Msg("GitHub API error")
		return "", &common.StatusError{Code: resp.StatusCode, Error: err}
	}

	// Parse response
	var user GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		g.logger.Error().Err(err).Msg("failed to parse GitHub user response")
		return "", &common.StatusError{Code: http.StatusInternalServerError, Error: err}
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
