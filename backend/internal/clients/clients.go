package clients

import (
	"github.com/Jacobgtd/hex-stats/backend/internal/alphavantage"
	"github.com/Jacobgtd/hex-stats/backend/internal/auth"
	"github.com/Jacobgtd/hex-stats/backend/internal/db"
	"github.com/Jacobgtd/hex-stats/backend/internal/github"
)

type Clients struct {
	GithubClient       *github.GithubClient
	AuthClient         *auth.AuthClient
	DBClient           *db.DBClient
	AlphavantageClient *alphavantage.AlphavantageClient
}
