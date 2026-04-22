// Package wallarm implements the Wallarm v2 API.
package wallarm

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// gzipReadCloser wraps a gzip.Reader and closes both the reader and the underlying body.
type gzipReadCloser struct {
	reader     *gzip.Reader
	underlying io.ReadCloser
}

func (g *gzipReadCloser) Read(p []byte) (int, error) { return g.reader.Read(p) }
func (g *gzipReadCloser) Close() error {
	g.reader.Close()
	return g.underlying.Close()
}

// ErrExistingResource is returned when resource was created other than Terrafom ways - directly via the API.
var ErrExistingResource = errors.New("This resource has already been created earlier")

// ErrInvalidCredentials is raised when not all the credentials are presented.
var ErrInvalidCredentials = errors.New("Credentials are not set. Specify Token or Pair of Secret and UUID")

// APIError represents an error response from the Wallarm API.
// Use errors.As to check for specific status codes:
//
//	var apiErr *wallarm.APIError
//	if errors.As(err, &apiErr) && apiErr.StatusCode == 404 { ... }
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP Status: %d, Body: %s", e.StatusCode, e.Body)
}

// NewAPIError creates an APIError with the given status code and body.
func NewAPIError(statusCode int, body string) *APIError {
	return &APIError{StatusCode: statusCode, Body: body}
}

// New creates a new Wallarm API client.
func New(opts ...Option) (API, error) {

	api, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func newClient(opts ...Option) (API, error) {
	defaultUserAgent := "Wallarm-go/" + Version

	api := &api{
		baseURL: apiURL,
		headers: make(http.Header),
		retryPolicy: RetryPolicy{
			MaxRetries:    12,
			MinRetryDelay: time.Duration(1) * time.Second,
			MaxRetryDelay: time.Duration(30) * time.Second,
		},
		Mutex: &sync.Mutex{},
	}

	if err := api.parseOptions(opts...); err != nil {
		return nil, errors.Wrap(err, "options parsing failed")
	}

	// Fall back to http.DefaultClient if the package user does not provide
	// their own.
	if api.httpClient == nil {
		api.httpClient = http.DefaultClient
		api.UserAgent = defaultUserAgent
	}

	return api, nil
}

// makeRequest makes an HTTP request and returns the body as a byte slice,
// closing it before returning. params will be serialized to JSON or string for GET query.
func (api *api) makeRequest(method, uri, reqType string, params interface{}, headers map[string]string) ([]byte, error) {
	return api.makeRequestContext(context.TODO(), method, uri, reqType, params, headers)
}

func (api *api) makeRequestContext(ctx context.Context, method, uri, reqType string, params interface{}, headers map[string]string) ([]byte, error) {
	// Replace nil with a JSON object if needed
	var (
		jsonBody []byte
		err      error
		resp     *http.Response
		respErr  error
		reqBody  io.Reader
		respBody []byte
	)

	if params != nil {
		if _, ok := params.(string); ok {
			jsonBody = nil
		} else if paramBytes, ok := params.([]byte); ok {
			jsonBody = paramBytes
		} else {
			jsonBody, err = json.Marshal(params)
			if err != nil {
				return nil, err
			}
		}
	} else {
		jsonBody = nil
	}

	var lastStatusCode int

	for i := 0; i <= api.retryPolicy.MaxRetries; i++ {
		if jsonBody != nil {
			reqBody = bytes.NewReader(jsonBody)
		}

		if i > 0 {
			// Use status-specific retry delays, capped by MaxRetryDelay.
			var sleepDuration time.Duration
			switch {
			case lastStatusCode == http.StatusLocked: // 423
				sleepDuration = 5 * time.Second
			case lastStatusCode >= 500:
				sleepDuration = 10 * time.Second
			default:
				sleepDuration = time.Duration(math.Pow(2, float64(i-1)) * float64(api.retryPolicy.MinRetryDelay))
			}
			if sleepDuration > api.retryPolicy.MaxRetryDelay {
				sleepDuration = api.retryPolicy.MaxRetryDelay
			}
			log.Printf("[DEBUG] Retrying request (attempt %d/%d) after %s due to HTTP %d",
				i+1, api.retryPolicy.MaxRetries+1, sleepDuration, lastStatusCode)
			time.Sleep(sleepDuration)
		}

		if query, ok := params.(string); ok {
			q := strings.NewReader(query)
			resp, err = api.request(ctx, method, uri, reqType, reqBody, q, headers)
			respErr = errors.Wrap(err, "could not make a request with get query")
		} else {
			resp, err = api.request(ctx, method, uri, reqType, reqBody, nil, headers)
			respErr = errors.Wrap(err, "could not make a request with JSON body")
		}

		if respErr != nil {
			lastStatusCode = 0
			continue
		}

		lastStatusCode = resp.StatusCode

		// Retry on: 423 (locked/rules being updated), 429 (rate limit), 5xx (server error).
		if resp.StatusCode == http.StatusLocked ||
			resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= 500 {
			respBody, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				respErr = errors.Wrap(err, "could not read response body")
			}
			continue
		}

		respBody, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, errors.Wrap(err, "could not read response body")
		}
		break
	}

	if respErr != nil {
		return nil, respErr
	}

	specificResourceProcessing := []string{"user"}
	switch {
	case resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices:
	case resp.StatusCode == http.StatusBadRequest && (reqType == "node" || reqType == "app" || reqType == "client") && string(respBody) == `{"status":400,"body":"Already exists"}`:
		return nil, errors.Wrap(ErrExistingResource, NewAPIError(resp.StatusCode, string(respBody)).Error())
	case resp.StatusCode == http.StatusConflict && Contains(specificResourceProcessing, reqType):
		return nil, errors.Wrap(ErrExistingResource, NewAPIError(resp.StatusCode, string(respBody)).Error())
	default:
		return nil, NewAPIError(resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (api *api) request(ctx context.Context, method, uri, reqType string, reqBody, query io.Reader, headers map[string]string) (*http.Response, error) {
	api.Lock()
	defer api.Unlock()

	req, err := http.NewRequestWithContext(ctx, method, api.baseURL+uri, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request creation failed")
	}

	for k, vs := range api.headers {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Accept-Encoding", "gzip")
	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	if req.Header.Get("Content-Type") == "" &&
		(Contains([]string{"POST", "PUT"}, method) && reqType != "userdetails") ||
		(method == "DELETE" && reqType == "ip_rules") {
		req.Header.Set("Content-Type", "application/json")
	} else if method == "GET" {
		req.Header.Del("Content-Type")
	}

	if method == "DELETE" && reqType != "ip_rules" {
		req.Header.Del("Content-Type")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if query != nil {
		q, err := ioutil.ReadAll(query)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = string(q)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request failed")
	}

	// Decompress gzip responses (we set Accept-Encoding: gzip explicitly).
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, errors.Wrap(err, "gzip decompression failed")
		}
		resp.Body = &gzipReadCloser{reader: reader, underlying: resp.Body}
		resp.Header.Del("Content-Encoding")
		resp.ContentLength = -1
	}

	return resp, nil
}
