package noaa

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/monitoring"
	"github.com/rs/zerolog"
)

const stationsPath = "mdapi/prod/webapi/stations.json"
const dataGetterPath = "api/prod/datagetter"

type NOAAClient struct {
	config     *NOAAClientConfig
	logger     zerolog.Logger
	httpClient *monitoring.HTTPClient
}

func NewNOAAClient(logger zerolog.Logger, config *NOAAClientConfig, httpClient *monitoring.HTTPClient) *NOAAClient {
	return &NOAAClient{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

func (n *NOAAClient) GetStations(ctx context.Context) ([]Station, error) {

	resp, err := n.httpClient.NewHTTPRequest(n.config.Url, stationsPath, http.MethodGet).
		WithTimeout(5 * time.Second).
		WithCacheTTL(24 * time.Hour).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var stationsResp StationsResponse
	err = resp.Unmarshal(&stationsResp)
	if err != nil {
		return nil, err
	}
	return stationsResp.Stations, nil
}

func validateDateOnly(t time.Time) error {
	if t.Hour() != 0 ||
		t.Minute() != 0 ||
		t.Second() != 0 ||
		t.Nanosecond() != 0 {
		return fmt.Errorf("date must not include a time component")
	}

	return nil
}

func (n *NOAAClient) getTides(ctx context.Context, stationID string, begin, end time.Time) (*TidePredictionsResponse, error) {

	params := url.Values{}
	params.Set("product", "predictions")
	params.Set("application", "hex-stats")
	params.Set("begin_date", begin.Format("20060102"))
	params.Set("end_date", end.Format("20060102"))
	params.Set("datum", "MLLW")
	params.Set("station", stationID)
	params.Set("time_zone", "gmt")
	params.Set("units", "english")
	params.Set("interval", "hilo")
	params.Set("format", "json")

	resp, err := n.httpClient.NewHTTPRequest(n.config.Url, dataGetterPath, http.MethodGet).
		WithQueryParams(params).
		WithTimeout(5 * time.Second).
		WithCacheTTL(24 * time.Hour).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var tideResp TidePredictionsResponse
	err = resp.Unmarshal(&tideResp)
	if err != nil {
		return nil, err
	}

	return &tideResp, nil
}

func (n *NOAAClient) getNearestTide(ctx context.Context, stationID string, tide TideType) (*TidePrediction, error) {

	now := time.Now()

	resp, err := n.getTides(ctx, stationID, now.Add(-24*time.Hour), now.Add(24*time.Hour))
	if err != nil {
		return nil, err
	}

	var nearestTide *TidePrediction
	var nearestTideDuration time.Duration

	for _, prediction := range resp.Predictions {

		timeDiffTide := prediction.Time.Sub(now).Abs()

		if prediction.Type == tide && (nearestTide == nil || timeDiffTide < nearestTideDuration) {
			nearestTide = &prediction
			nearestTideDuration = timeDiffTide

		}
	}

	if nearestTide == nil {
		return nil, fmt.Errorf("no %s tide found", tide)
	}
	return nearestTide, nil

}

func (n *NOAAClient) GetNearestHighTide(ctx context.Context, stationID string) (*TidePrediction, error) {
	return n.getNearestTide(ctx, stationID, HighTide)
}

func (n *NOAAClient) GetNearestLowTide(ctx context.Context, stationID string) (*TidePrediction, error) {
	return n.getNearestTide(ctx, stationID, LowTide)
}
