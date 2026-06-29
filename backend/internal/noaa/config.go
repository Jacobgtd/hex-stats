package noaa

import "github.com/Jacobgtd/hex-stats/backend/internal/configpack"

type NoaaClientConfig struct {
	Url string
}

func LoadNoaaConfig() (*NoaaClientConfig, error) {
	err := configpack.Load("noaa.config")
	if err != nil {
		return nil, err
	}

	url, err := configpack.String("URL")
	return &NoaaClientConfig{
		Url: url,
	}, nil
}
