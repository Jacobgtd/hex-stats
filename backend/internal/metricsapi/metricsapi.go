package metricsapi

import (
	"net/http"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/clients"
	"github.com/gin-gonic/gin"
)

func RegisterMetricsRoutes(rg *gin.RouterGroup, clients *clients.Clients) {
	rg.GET("/stocks/:symbol/price", getStockPriceHandler(clients))
	rg.GET("/stocks/:symbol/change-percent", getStockPercentChangeHandler(clients))
}

func getStockPriceHandler(clients *clients.Clients) func(*gin.Context) {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		// Fetch stock price for the given symbol using clients
		price, err := clients.AlphavantageClient.GetTickerPrice(c.Request.Context(), symbol)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock price"})
			return
		}

		resp := gaugeMetricResponse{
			genericMetricResponse: genericMetricResponse{
				RefreshAt: uint64(time.Now().Unix()),
			},
			Value:         price,
			DisplayFormat: "currency",
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

		resp := gaugeMetricResponse{
			genericMetricResponse: genericMetricResponse{
				RefreshAt: uint64(time.Now().Unix()),
			},
			Value:         percentChange,
			DisplayFormat: "percent",
		}

		c.JSON(http.StatusOK, resp)

	}
}
