package metricshandlers

import (
	"github.com/Jacobgtd/hex-stats/backend/internal/clients"
	"github.com/Jacobgtd/hex-stats/backend/internal/metricsregistry"
	"github.com/gin-gonic/gin"
)

func RegisterMetricsRoutes(rg *gin.RouterGroup, clients *clients.Clients) error {

	metricsRegistry := metricsregistry.NewMetricsRegistry(rg)

	// Stocks
	stockSymbolParam := metricsregistry.NewMetricsParam(
		"symbol",
		"Ticker Symbol",
		"stock ticker symbol, ex: AAPL",
	)

	stockPriceEntry, err := metricsregistry.NewMetricsRegistryEntry(
		"Stock Price",
		"Get a stock price over the last 24 hours",
		"/stocks/:symbol/price",
		metricsregistry.MetricTypeGauge,
		[]metricsregistry.MetricsParam{*stockSymbolParam},
		metricsregistry.DisplayFormatUSD,
		getStockPriceHandler(clients),
	)

	if err != nil {
		return err
	}

	stockChangeEntry, err := metricsregistry.NewMetricsRegistryEntry(
		"Stock Change",
		"Get the change in stock price over the last 24 hours",
		"/stocks/:symbol/change-percent",
		metricsregistry.MetricTypeGauge,
		[]metricsregistry.MetricsParam{*stockSymbolParam},
		metricsregistry.DisplayFormatPercent,
		getStockPercentChangeHandler(clients),
	)
	if err != nil {
		return err
	}

	metricsRegistry.RegisterEntry(
		*stockPriceEntry,
		*stockChangeEntry,
	)

	// Tides

	stationParam := metricsregistry.NewMetricsParam(
		"stationID",
		"Station ID",
		"ID of the tide station, ex: 9447130",
	)

	metricsRegistry.RegisterParamOptions(*stationParam, getStationsFunc(clients))

	lowTidesEntry, err := metricsregistry.NewMetricsRegistryEntry(
		"Low Tide",
		"Get the next low tide for a given station",
		"/tides/:stationID/low",
		metricsregistry.MetricTypeGauge,
		[]metricsregistry.MetricsParam{*stationParam},
		metricsregistry.DisplayFormatSeconds,
		getLowTideHandler(clients),
	)
	if err != nil {
		return err
	}

	highTidesEntry, err := metricsregistry.NewMetricsRegistryEntry(
		"High Tide",
		"Get the next high tide for a given station",
		"/tides/:stationID/high",
		metricsregistry.MetricTypeGauge,
		[]metricsregistry.MetricsParam{*stationParam},
		metricsregistry.DisplayFormatSeconds,
		getHighTideHandler(clients),
	)
	if err != nil {
		return err
	}

	metricsRegistry.RegisterEntry(
		*lowTidesEntry,
		*highTidesEntry,
	)

	return nil
}
