package alphavantage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/monitoring"
)

type AlphavantageClient struct {
	config     *AlphavantageClientConfig
	httpClient *monitoring.HTTPClient
}

func NewAlphavantageClient(config *AlphavantageClientConfig, httpClient *monitoring.HTTPClient) *AlphavantageClient {
	return &AlphavantageClient{
		config:     config,
		httpClient: httpClient,
	}
}

func nextMidnightUTC(t time.Time) time.Time {
	t = t.UTC()
	year, month, day := t.Date()

	return time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC)
}

func (c *AlphavantageClient) getTickerData(ctx context.Context, ticker string) (*GlobalQuoteResponse, error) {

	validUntil := nextMidnightUTC(time.Now())

	params := url.Values{}
	params.Set("function", "GLOBAL_QUOTE")
	params.Set("symbol", ticker)
	params.Set("apikey", c.config.APIKey)

	resp, err := c.httpClient.NewHTTPRequest(c.config.URL, "query", http.MethodGet).
		WithQueryParams(params).
		WithTimeout(10 * time.Second).
		WithCacheExpiry(validUntil).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	quote := &GlobalQuoteResponse{}

	err = resp.Unmarshal(quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (c *AlphavantageClient) GetTickerPrice(ctx context.Context, ticker string) (float64, error) {
	quote, err := c.getTickerData(ctx, ticker)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(quote.GlobalQuote.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func ParsePercentToFloat(s string) (float64, error) {
	// Clean whitespace and remove the "%" symbol if it exists
	cleaned := strings.TrimSpace(s)
	cleaned = strings.TrimSuffix(cleaned, "%")
	cleaned = strings.TrimSpace(cleaned) // Handle any spaces left before the %

	// Parse the remaining string as a 64-bit float
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse percentage string '%s': %w", s, err)
	}

	return val, nil
}

func (c *AlphavantageClient) GetTickerChangePercent(ctx context.Context, ticker string) (float64, error) {
	quote, err := c.getTickerData(ctx, ticker)
	if err != nil {
		return 0, err
	}

	changePercent, err := ParsePercentToFloat(quote.GlobalQuote.ChangePercent)
	if err != nil {
		return 0, err
	}

	return changePercent, nil
}
