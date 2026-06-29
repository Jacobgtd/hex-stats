package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"
)

type HttpRequest struct {
	ctx                     context.Context
	baseUrl                 string
	path                    string
	method                  string
	headers                 http.Header
	body                    interface{}
	queryParams             url.Values
	timeout                 *time.Duration
	expectedSuccessCode     []int
	expectedFailureCode     []int
	useExpectedFailureCodes bool
}

func NewHttpRequest(baseUrl, path, method string) HttpRequest {
	return HttpRequest{
		baseUrl:                 baseUrl,
		path:                    path,
		method:                  method,
		useExpectedFailureCodes: false,
	}
}

func (r *HttpRequest) WithCtx(ctx context.Context) *HttpRequest {
	r.ctx = ctx
	return r
}

func (r *HttpRequest) WithHeaders(h http.Header) *HttpRequest {
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

func (r *HttpRequest) WithBearerToken(token string) *HttpRequest {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
}

func (r *HttpRequest) WithBody(h http.Header) *HttpRequest {
	r.headers = h
	return r
}

func (r *HttpRequest) WithQueryParams(p url.Values) *HttpRequest {
	r.queryParams = p
	return r
}

func (r *HttpRequest) WithTimeout(t time.Duration) *HttpRequest {
	var timeout time.Duration
	timeout = t
	r.timeout = &timeout
	return r
}

func (r *HttpRequest) WithExpectedSuccessCode(codes ...int) *HttpRequest {
	r.expectedSuccessCode = append(r.expectedSuccessCode, codes...)
	return r
}

func (r *HttpRequest) WithExpectedFailureCode(codes ...int) *HttpRequest {
	r.expectedFailureCode = append(r.expectedFailureCode, codes...)
	r.useExpectedFailureCodes = true
	return r
}

func (r *HttpRequest) Do(res interface{}) (int, error) {
	// Make request

	if r.ctx == nil {
		newCtx := context.Background()
		r.ctx = newCtx
	}

	if r.timeout != nil {
		newCtx, cancel := context.WithTimeout(r.ctx, *r.timeout)
		defer cancel()
		r.ctx = newCtx
	}

	url := fmt.Sprintf("%s/%s", r.baseUrl, r.path)
	if r.queryParams != nil {
		url = fmt.Sprintf("%s?%s", url, r.queryParams.Encode())
	}

	req, err := http.NewRequestWithContext(r.ctx, r.method, fmt.Sprintf("%s/%s", r.baseUrl, r.path), nil)
	req.Header = r.headers
	if err != nil {
		return http.StatusInternalServerError, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if r.expectedSuccessCode == nil {
		r.expectedSuccessCode = []int{http.StatusOK}
	}
	responseCodeIsSuccess := slices.Contains(r.expectedSuccessCode, resp.StatusCode)

	responseCodeIsFailure := !responseCodeIsSuccess
	if r.useExpectedFailureCodes {
		responseCodeIsFailure = slices.Contains(r.expectedFailureCode, resp.StatusCode)
	}

	if !responseCodeIsSuccess && !responseCodeIsFailure {
		return http.StatusInternalServerError, fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if responseCodeIsFailure {
		return resp.StatusCode, nil
	}

	if res != nil {
		err = json.Unmarshal(body, res)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode, nil
}
