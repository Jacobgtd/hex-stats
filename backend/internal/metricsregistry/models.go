package metricsregistry

type MetricType string

const (
	MetricTypeGauge MetricType = "gauge"
)

type DisplayFormat string

const (
	DisplayFormatUSD     DisplayFormat = "usd"
	DisplayFormatPercent DisplayFormat = "percent"
	DisplayFormatSeconds DisplayFormat = "seconds"
)

type GenericMetricResponse struct {
	RefreshAt uint64 `json:"refresh_at"`
}

type GaugeMetricResponse struct {
	GenericMetricResponse
	Value float64 `json:"value"`
}
