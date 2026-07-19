package clients

import (
	"github.com/Jacobgtd/hex-stats/backend/internal/alphavantage"
	"github.com/Jacobgtd/hex-stats/backend/internal/auth"
	"github.com/Jacobgtd/hex-stats/backend/internal/db"
	"github.com/Jacobgtd/hex-stats/backend/internal/github"
	"github.com/Jacobgtd/hex-stats/backend/internal/noaa"
)

type Clients struct {
	GithubClient       *github.GithubClient
	AuthClient         *auth.AuthClient
	DBClient           *db.DBClient
	AlphavantageClient *alphavantage.AlphavantageClient
	NOAAClient         *noaa.NOAAClient
}
