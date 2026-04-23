package wallarm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPISpecCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"id": 111680,
				"client_id": 8649,
				"title": "TEST_SPEC",
				"description": "test",
				"status": "ready",
				"instances": [],
				"domains": [],
				"regular_file_update": true,
				"api_detection": false,
				"spec_version": "1.0.0",
				"version": 0,
				"endpoints_count": 3,
				"openapi_version": "3.0.0"
			}
		}`)
	})

	res, err := client.APISpecCreate(&APISpecCreate{
		Title:             "TEST_SPEC",
		Description:       "test",
		FileRemoteURL:     "https://example.com/spec.yaml",
		RegularFileUpdate: true,
		ClientID:          8649,
	})
	assert.NoError(t, err)
	assert.Equal(t, 111680, res.Body.ID)
	assert.Equal(t, "TEST_SPEC", res.Body.Title)
}

func TestAPISpecReadByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111680", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"id": 111680,
				"client_id": 8649,
				"title": "TEST_SPEC",
				"status": "ready"
			}
		}`)
	})

	res, err := client.APISpecReadByID(8649, 111680)
	assert.NoError(t, err)
	assert.Equal(t, 111680, res.ID)
}

func TestAPISpecDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111680", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status": 200}`)
	})

	err := client.APISpecDelete(8649, 111680)
	assert.NoError(t, err)
}

func TestAPISpecReadByID_NotFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111680", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status": 200, "body": null}`)
	})

	_, err := client.APISpecReadByID(8649, 111680)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound, got %v", err)
}

func TestAPISpecReadByID_HTTP404(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111961", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"status":404,"body":"Api spec with ID:111961 not found"}`)
	})

	_, err := client.APISpecReadByID(8649, 111961)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound, got %v", err)
}

func TestAPISpecUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111680", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v4/clients/8649/rules/api-specs/111680", r.RequestURI)
		body, _ := io.ReadAll(r.Body)
		var req APISpecUpdate
		assert.NoError(t, json.Unmarshal(body, &req))
		if assert.NotNil(t, req.Title) {
			assert.Equal(t, "UPDATED_SPEC", *req.Title)
		}
		if assert.NotNil(t, req.Description) {
			assert.Equal(t, "updated description", *req.Description)
		}
		assert.Nil(t, req.FileRemoteURL)
		assert.Nil(t, req.RegularFileUpdate)
		assert.Nil(t, req.APIDetection)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"id": 111680,
				"client_id": 8649,
				"title": "UPDATED_SPEC",
				"description": "updated description",
				"status": "ready"
			}
		}`)
	})

	title := "UPDATED_SPEC"
	description := "updated description"
	res, err := client.APISpecUpdate(8649, 111680, &APISpecUpdate{
		Title:       &title,
		Description: &description,
	})
	assert.NoError(t, err)
	assert.Equal(t, 200, res.Status)
	assert.NotNil(t, res.Body)
	assert.Equal(t, 111680, res.Body.ID)
	assert.Equal(t, "UPDATED_SPEC", res.Body.Title)
	assert.Equal(t, "updated description", res.Body.Description)
}

func TestAPISpecList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v4/clients/8649/rules/api-specs?page=2&per_page=50", r.RequestURI)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"items": [
				{"id": 111680, "client_id": 8649, "title": "SPEC_A", "status": "ready"},
				{"id": 111681, "client_id": 8649, "title": "SPEC_B", "status": "ready"}
			],
			"current_page": 2,
			"per_page": 50,
			"total_pages": 3,
			"total_count": 120
		}`)
	})

	resp, err := client.APISpecList(8649, 2, 50)
	assert.NoError(t, err)
	assert.Equal(t, 120, resp.TotalCount)
	assert.Equal(t, 2, resp.CurrentPage)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, 111680, resp.Items[0].ID)
}

func TestAPISpecPolicyPut(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111680/policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v4/clients/8649/rules/api-specs/111680/policy", r.RequestURI)
		body, _ := io.ReadAll(r.Body)
		var req APISpecPolicy
		assert.NoError(t, json.Unmarshal(body, &req))
		assert.True(t, req.Enabled)
		assert.Equal(t, "block", req.UndefinedEndpointMode)
		assert.Equal(t, "monitor", req.UndefinedParameterMode)
		assert.Equal(t, "block", req.MissingParameterMode)
		assert.Equal(t, "monitor", req.InvalidParameterValueMode)
		assert.Equal(t, "block", req.MissingAuthMode)
		assert.Equal(t, "monitor", req.InvalidRequestMode)
		assert.Equal(t, "block", req.TimeoutMode)
		assert.Equal(t, "block", req.MaxRequestSizeMode)
		assert.Equal(t, 5000, req.Timeout)
		assert.Equal(t, 1024, req.MaxRequestSize)
		assert.Len(t, req.Conditions, 1)
		assert.Equal(t, "equal", req.Conditions[0].Type)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"enabled": true,
				"undefined_endpoint_mode": "block",
				"undefined_parameter_mode": "monitor",
				"missing_parameter_mode": "block",
				"invalid_parameter_value_mode": "monitor",
				"missing_auth_mode": "block",
				"invalid_request_mode": "monitor",
				"timeout_mode": "block",
				"max_request_size_mode": "block",
				"timeout": 5000,
				"max_request_size": 1024,
				"conditions": [
					{"type": "equal", "value": "example.com", "point": ["header", "HOST"]}
				]
			}
		}`)
	})

	policy := &APISpecPolicy{
		Enabled:                   true,
		UndefinedEndpointMode:     "block",
		UndefinedParameterMode:    "monitor",
		MissingParameterMode:      "block",
		InvalidParameterValueMode: "monitor",
		MissingAuthMode:           "block",
		InvalidRequestMode:        "monitor",
		TimeoutMode:               "block",
		MaxRequestSizeMode:        "block",
		Timeout:                   5000,
		MaxRequestSize:            1024,
		Conditions: []APISpecPolicyCondition{
			{Type: "equal", Value: "example.com", Point: []interface{}{"header", "HOST"}},
		},
	}
	res, err := client.APISpecPolicyPut(8649, 111680, policy)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.Status)
	assert.NotNil(t, res.Body)
	assert.True(t, res.Body.Enabled)
	assert.Equal(t, "block", res.Body.UndefinedEndpointMode)
	assert.Equal(t, "monitor", res.Body.UndefinedParameterMode)
	assert.Equal(t, "block", res.Body.MissingParameterMode)
	assert.Equal(t, "monitor", res.Body.InvalidParameterValueMode)
	assert.Equal(t, "block", res.Body.MissingAuthMode)
	assert.Equal(t, "monitor", res.Body.InvalidRequestMode)
	assert.Equal(t, "block", res.Body.TimeoutMode)
	assert.Equal(t, "block", res.Body.MaxRequestSizeMode)
	assert.Equal(t, 5000, res.Body.Timeout)
	assert.Equal(t, 1024, res.Body.MaxRequestSize)
	assert.Len(t, res.Body.Conditions, 1)
}

func TestAPISpecPolicyPut_Disable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs/111680/policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		body, _ := io.ReadAll(r.Body)
		var req APISpecPolicy
		assert.NoError(t, json.Unmarshal(body, &req))
		assert.False(t, req.Enabled)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status": 200, "body": {"enabled": false}}`)
	})

	res, err := client.APISpecPolicyPut(8649, 111680, &APISpecPolicy{Enabled: false})
	assert.NoError(t, err)
	assert.NotNil(t, res.Body)
	assert.False(t, res.Body.Enabled)
}
