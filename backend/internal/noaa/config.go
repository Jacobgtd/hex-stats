package noaa

import "github.com/Jacobgtd/hex-stats/backend/internal/configpack"

type NOAAClientConfig struct {
	Url string
}

func LoadNOAAClientConfig() (*NOAAClientConfig, error) {
	err := configpack.Load("noaa.config")
	if err != nil {
		return nil, err
	}

	url, err := configpack.String("URL")
	return &NOAAClientConfig{
		Url: url,
	}, nil
}
