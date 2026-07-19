package metricshandlers

import (
	"context"
	"fmt"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/clients"
	"github.com/Jacobgtd/hex-stats/backend/internal/metricsregistry"
	"github.com/Jacobgtd/hex-stats/backend/internal/noaa"
	"github.com/gin-gonic/gin"
)

func getTideHandler(getTideFunc func(ctx context.Context, stationID string) (*noaa.TidePrediction, error)) func(*gin.Context) {
	return func(c *gin.Context) {

		tide, err := getTideFunc(c.Request.Context(), c.Param("stationID"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, metricsregistry.GaugeMetricResponse{
			GenericMetricResponse: metricsregistry.GenericMetricResponse{
				RefreshAt: uint64(time.Now().Add(5 * time.Minute).Unix()),
			},
			Value: float64(absInt(tide.Time.Unix() - time.Now().Unix())),
		})
	}
}

func getLowTideHandler(clients *clients.Clients) func(*gin.Context) {
	return getTideHandler(clients.NOAAClient.GetNearestLowTide)
}
func getHighTideHandler(clients *clients.Clients) func(*gin.Context) {
	return getTideHandler(clients.NOAAClient.GetNearestHighTide)
}

func getStationsFunc(clients *clients.Clients) func(ctx context.Context) ([]metricsregistry.MetricsParamEntry, error) {
	return func(ctx context.Context) ([]metricsregistry.MetricsParamEntry, error) {
		stations, err := clients.NOAAClient.GetStations(ctx)
		if err != nil {
			return nil, err
		}
		var result []metricsregistry.MetricsParamEntry
		for _, station := range stations {
			result = append(result, metricsregistry.MetricsParamEntry{
				DisplayName: fmt.Sprintf("%s - %s", station.State, station.Name),
				Value:       station.ID,
			})
		}
		return result, nil
	}
}
