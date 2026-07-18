package monitoring

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/cache"
)

type HTTPClient struct {
	client http.Client
	cache  cache.Cache
}

func NewHTTPClient(client http.Client, cache cache.Cache) *HTTPClient {
	return &HTTPClient{
		client: client,
		cache:  cache,
	}
}

type HTTPRequest struct {
	HTTPClient            *HTTPClient
	baseURL               string
	path                  string
	method                string
	headers               http.Header
	body                  interface{}
	queryParams           url.Values
	timeout               *time.Duration
	possibleResponseCodes []int
	useCache              bool
	cachettl              *time.Duration
	cacheExpiry           *time.Time
}

func (c *HTTPClient) NewHTTPRequest(baseUrl, path, method string) *HTTPRequest {
	return &HTTPRequest{
		HTTPClient: c,
		baseURL:    baseUrl,
		path:       path,
		method:     method,
		useCache:   false,
	}
}

func (r *HTTPRequest) WithHeaders(h http.Header) *HTTPRequest {
	if r.headers == nil {
		r.headers = h
		return r
	}

	for header, values := range h {
		for _, value := range values {
			r.headers.Add(header, value)
		}
	}

	return r
}

func (r *HTTPRequest) WithBearerToken(token string) *HTTPRequest {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
}

func (r *HTTPRequest) WithBody(h http.Header) *HTTPRequest {
	r.headers = h
	return r
}

func (r *HTTPRequest) WithCacheTTL(ttl time.Duration) *HTTPRequest {
	r.useCache = true
	r.cachettl = &ttl
	return r
}

func (r *HTTPRequest) WithCacheExpiry(expiry time.Time) *HTTPRequest {
	r.useCache = true
	r.cacheExpiry = &expiry
	return r
}

func (r *HTTPRequest) WithQueryParams(p url.Values) *HTTPRequest {
	r.queryParams = p
	return r
}

func (r *HTTPRequest) WithTimeout(t time.Duration) *HTTPRequest {
	var timeout time.Duration
	timeout = t
	r.timeout = &timeout
	return r
}

func (r *HTTPRequest) WithPossibleResponseCodes(codes ...int) *HTTPRequest {
	r.possibleResponseCodes = append(r.possibleResponseCodes, codes...)
	return r
}

type httpResult struct {
	StatusCode int
	Body       []byte
}
type HttpResponse struct {
	result httpResult
}

func (r *HttpResponse) StatusCode() int {
	return r.result.StatusCode
}

func (r *HttpResponse) Unmarshal(res interface{}) error {
	err := json.Unmarshal(r.result.Body, res)
	return err
}

func (r *HTTPRequest) Do(ctx context.Context) (*HttpResponse, error) {
	// Make request
	requestHash := r.hash()

	if ctx == nil {
		ctx = context.Background()
	}

	if r.timeout != nil {
		newCtx, cancel := context.WithTimeout(ctx, *r.timeout)
		defer cancel()
		ctx = newCtx
	}

	if r.useCache {
		resp := &httpResult{}
		err := r.HTTPClient.cache.Get("http-client", requestHash, resp).Do(ctx)
		if err == nil {
			return &HttpResponse{result: *resp}, nil
		}
	}

	url := fmt.Sprintf("%s/%s", r.baseURL, r.path)
	if r.queryParams != nil {
		url = fmt.Sprintf("%s?%s", url, r.queryParams.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, r.method, url, nil)
	req.Header = r.headers
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if r.possibleResponseCodes == nil {
		r.possibleResponseCodes = []int{http.StatusOK}
	}

	if !slices.Contains(r.possibleResponseCodes, resp.StatusCode) {
		return nil, fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := httpResult{
		StatusCode: resp.StatusCode,
		Body:       body,
	}

	if r.useCache {
		var err error
		if r.cacheExpiry != nil {
			err = r.HTTPClient.cache.Set("http-client", requestHash, result).WithExpiry(*r.cacheExpiry).Do(ctx)
		} else {
			err = r.HTTPClient.cache.Set("http-client", requestHash, result).WithTTL(*r.cachettl).Do(ctx)
		}
		if err != nil {
			// Log the error but don't fail the request
		}
	}

	return &HttpResponse{
		result: result,
	}, nil
}

func (r *HTTPRequest) hash() string {
	hash := hashObjects(r.baseURL, r.path, r.method, r.queryParams, r.headers, r.body)
	return fmt.Sprintf("%x", hash)
}

func hashObjects(objects ...any) [32]byte {
	data, err := json.Marshal(objects)
	if err != nil {
		panic(err)
	}

	return sha256.Sum256(data)
}
