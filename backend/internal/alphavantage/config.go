package alphavantage

import "github.com/Jacobgtd/hex-stats/backend/internal/configpack"

type AlphavantageClientConfig struct {
	APIKey string
	URL    string
}

func LoadAlphavantageClientConfig() (*AlphavantageClientConfig, error) {
	configpack.Load("alphavantage.config")

	apiKey, err := configpack.String("API_KEY")
	if err != nil {
		return nil, err
	}

	URL, err := configpack.String("URL")
	if err != nil {
		return nil, err
	}

	return &AlphavantageClientConfig{
		APIKey: apiKey,
		URL:    URL,
	}, nil
}
