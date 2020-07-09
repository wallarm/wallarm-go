package wallarm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// New creates a new Wallarm API client.
func New(apiURL, uuid, secret string, opts ...Option) (*API, error) {
	if uuid == "" || secret == "" {
		return nil, errors.New("Credentials are not set")
	}

	api, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	api.baseURL = apiURL
	api.apiUUID = uuid
	api.apiSecret = secret
	api.headers.Add("X-WallarmAPI-UUID", uuid)
	api.headers.Add("X-WallarmAPI-Secret", secret)

	return api, nil
}

func newClient(opts ...Option) (*API, error) {
	silentLogger := log.New(ioutil.Discard, "", log.LstdFlags)

	api := &API{
		baseURL: apiURL,
		headers: make(http.Header),
		retryPolicy: RetryPolicy{
			MaxRetries:    3,
			MinRetryDelay: time.Duration(1) * time.Second,
			MaxRetryDelay: time.Duration(30) * time.Second,
		},
		logger: silentLogger,
	}

	err := api.parseOptions(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "options parsing failed")
	}

	// Fall back to http.DefaultClient if the package user does not provide
	// their own.
	if api.httpClient == nil {
		api.httpClient = http.DefaultClient
	}

	return api, nil
}

// makeRequest makes a HTTP request and returns the body as a byte slice,
// closing it before returning. params will be serialized to JSON.
func (api *API) makeRequest(method, uri, reqType string, params interface{}) ([]byte, error) {
	return api.makeRequestContext(context.TODO(), method, uri, reqType, params)
}

func (api *API) makeRequestContext(ctx context.Context, method, uri, reqType string, params interface{}) ([]byte, error) {
	// Replace nil with a JSON object if needed
	var jsonBody []byte
	var err error

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

	var resp *http.Response
	var respErr error
	var reqBody io.Reader
	var respBody []byte
	for i := 0; i <= api.retryPolicy.MaxRetries; i++ {
		if jsonBody != nil {
			reqBody = bytes.NewReader(jsonBody)
		}

		if i > 0 {
			// expect the backoff introduced here on errored requests to dominate the effect of rate limiting
			// don't need a random component here as the rate limiter should do something similar
			// nb time duration could truncate an arbitrary float. Since our inputs are all ints, we should be ok
			sleepDuration := time.Duration(math.Pow(2, float64(i-1)) * float64(api.retryPolicy.MinRetryDelay))

			if sleepDuration > api.retryPolicy.MaxRetryDelay {
				sleepDuration = api.retryPolicy.MaxRetryDelay
			}
			// useful to do some simple logging here, maybe introduce levels later
			api.logger.Printf("Sleeping %s before retry attempt number %d for request %s %s", sleepDuration.String(), i, method, uri)
			time.Sleep(sleepDuration)

		}

		if query, ok := params.(string); ok {
			q := strings.NewReader(query)
			resp, respErr = api.request(ctx, method, uri, reqType, reqBody, q)
		} else {
			resp, respErr = api.request(ctx, method, uri, reqType, reqBody, nil)
		}

		// retry if the server is rate limiting us or if it failed
		// assumes server operations are rolled back on failure
		if respErr != nil || resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			if respErr == nil {
				respBody, err = ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				respErr = errors.Wrap(err, "could not read response body")

				api.logger.Printf("Request: %s %s got an error response %d: %s\n", method, uri, resp.StatusCode,
					strings.Replace(strings.Replace(string(respBody), "\n", "", -1), "\t", "", -1))
			} else {
				api.logger.Printf("Error performing request: %s %s : %s \n", method, uri, respErr.Error())
			}
			continue
		} else {
			respBody, err = ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				return nil, errors.Wrap(err, "could not read response body")
			}
			break
		}
	}
	if respErr != nil {
		return nil, respErr
	}

	specificResourceProcessing := []string{"scanner", "user"}

	switch {
	case resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices:
	case resp.StatusCode == http.StatusUnauthorized:
		return nil, errors.Errorf("Status code: %d, Body: %s", resp.StatusCode, respBody)
	case resp.StatusCode == http.StatusForbidden:
		return nil, errors.Errorf("Status code: %d, Body: %s", resp.StatusCode, respBody)
	case resp.StatusCode == http.StatusServiceUnavailable,
		resp.StatusCode == http.StatusBadGateway,
		resp.StatusCode == http.StatusGatewayTimeout,
		resp.StatusCode == 522,
		resp.StatusCode == 523,
		resp.StatusCode == 524:
		return nil, errors.Errorf("Status code: %d, Body: %s", resp.StatusCode, respBody)
	case resp.StatusCode == http.StatusBadRequest && (reqType == "node" || reqType == "app") && string(respBody) == `{"status":400,"body":"Already exists"}`:
		err := &ExistingResourceError{Status: resp.StatusCode, Body: string(respBody)}
		return nil, err
	case resp.StatusCode == http.StatusConflict && Contains(specificResourceProcessing, reqType):
		err := &ExistingResourceError{Status: resp.StatusCode, Body: string(respBody)}
		return nil, err
	default:
		return nil, errors.Errorf("Status code: %d, Body: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}

func (api *API) request(ctx context.Context, method, uri, reqType string, reqBody, query io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, api.baseURL+uri, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request creation failed")
	}
	req.WithContext(ctx)

	req.Header = api.headers

	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	methods := []string{"POST", "PUT"}

	if req.Header.Get("Content-Type") == "" && Contains(methods, method) && reqType != "userdetails" {
		req.Header.Set("Content-Type", "application/json")
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
	return resp, nil
}

// Contains wraps methods to check if List contains an element.
func Contains(a interface{}, x interface{}) bool {
	switch x.(type) {
	case int:
		group := a.([]int)
		return intInList(group, x.(int))
	case string:
		group := a.([]string)
		return strInList(group, x.(string))
	default:
		return false
	}
}

func strInList(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func intInList(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
