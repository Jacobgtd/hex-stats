package metricsapi

type MetricType string

const (
	MetricTypeGauge MetricType = "gauge"
)

type genericMetricResponse struct {
	RefreshAt uint64 `json:"refresh_at"`
}

type gaugeMetricResponse struct {
	genericMetricResponse
	Value         float64 `json:"value"`
	DisplayFormat string  `json:"display_format"`
}
