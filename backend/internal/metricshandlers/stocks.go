package metricshandlers

import (
	"net/http"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/clients"
	"github.com/Jacobgtd/hex-stats/backend/internal/metricsregistry"
	"github.com/gin-gonic/gin"
)

func getStockPriceHandler(clients *clients.Clients) func(*gin.Context) {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		// Fetch stock price for the given symbol using clients
		price, err := clients.AlphavantageClient.GetTickerPrice(c.Request.Context(), symbol)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock price"})
			return
		}

		resp := metricsregistry.GaugeMetricResponse{
			GenericMetricResponse: metricsregistry.GenericMetricResponse{
				RefreshAt: uint64(time.Now().Unix()),
			},
			Value: price,
		}

		c.JSON(http.StatusOK, resp)

	}
}

func getStockPercentChangeHandler(clients *clients.Clients) func(*gin.Context) {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		// Fetch stock percent change for the given symbol using clients
		percentChange, err := clients.AlphavantageClient.GetTickerChangePercent(c.Request.Context(), symbol)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock percent change"})
			return
		}

		resp := metricsregistry.GaugeMetricResponse{
			GenericMetricResponse: metricsregistry.GenericMetricResponse{
				RefreshAt: uint64(time.Now().Unix()),
			},
			Value: percentChange,
		}

		c.JSON(http.StatusOK, resp)

	}
}
