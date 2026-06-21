package github

import "github.com/Jacobgtd/hex-stats/backend/internal/configpack"

type GithubClientConfig struct {
	Url    string
	Admins []string
}

func LoadGithubClientConfig() (*GithubClientConfig, error) {

	err := configpack.Load("github.config")
	if err != nil {
		return nil, err
	}

	url, err := configpack.String("URL")
	if err != nil {
		return nil, err
	}

	admins, err := configpack.StringSliceFromEnv("ADMINS")
	if err != nil {
		return nil, err
	}

	return &GithubClientConfig{
		Url:    url,
		Admins: admins,
	}, nil
}
