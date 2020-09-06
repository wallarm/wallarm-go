package wallarm

import (
	"net/http"
	"sync"
	"time"
)

// APIURL is API host domain of the Wallarm Cloud
var apiURL = "https://api.wallarm.com"

// Option is a functional option for configuring the API client
type Option func(*API) error

// RetryPolicy specifies number of retries and min/max retry delays
// This config is used when the client exponentially backs off after errored requests
type RetryPolicy struct {
	MaxRetries    int
	MinRetryDelay time.Duration
	MaxRetryDelay time.Duration
}

// Logger defines the interface this library needs to use logging
// This is a subset of the methods implemented in the log package
type Logger interface {
	Printf(format string, v ...interface{})
}

// API holds the configuration for the current API client. A client should not
// be modified concurrently.
type API struct {
	baseURL, UserAgent string
	ClientID           int
	headers            http.Header
	httpClient         *http.Client
	retryPolicy        RetryPolicy
	logger             Logger
	*sync.Mutex
}

// HTTPClient accepts a custom *http.Client for making API calls.
func HTTPClient(client *http.Client) Option {
	return func(api *API) error {
		api.httpClient = client
		return nil
	}
}

// Headers allows you to set custom HTTP headers when making API calls (e.g. for
// satisfying HTTP proxies, or for debugging).
func Headers(headers http.Header) Option {
	return func(api *API) error {
		api.headers = headers
		return nil
	}
}

// UserAgent allows to set custome User-Agent header.
func UserAgent(userAgent string) Option {
	return func(api *API) error {
		api.UserAgent = userAgent
		return nil
	}
}

// UsingRetryPolicy applies a non-default number of retries and min/max retry delays
// This will be used when the client exponentially backs off after errored requests
func UsingRetryPolicy(maxRetries int, minRetryDelaySecs int, maxRetryDelaySecs int) Option {
	// seconds is very granular for a minimum delay - but this is only in case of failure
	return func(api *API) error {
		api.retryPolicy = RetryPolicy{
			MaxRetries:    maxRetries,
			MinRetryDelay: time.Duration(minRetryDelaySecs) * time.Second,
			MaxRetryDelay: time.Duration(maxRetryDelaySecs) * time.Second,
		}
		return nil
	}
}

// UsingLogger can be set if you want to get log output from this API instance
// By default no log output is emitted
func UsingLogger(logger Logger) Option {
	return func(api *API) error {
		api.logger = logger
		return nil
	}
}

// parseOptions parses the supplied options functions and returns a configured
// *API instance.
func (api *API) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(api)
		if err != nil {
			return err
		}
	}

	return nil
}
