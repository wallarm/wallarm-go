package wallarm

import (
	"fmt"
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

func TestAPISpecRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v4/clients/8649/rules/api-specs", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"items": [{
				"id": 111680,
				"client_id": 8649,
				"title": "TEST_SPEC",
				"status": "ready"
			}],
			"current_page": 1,
			"per_page": 20,
			"total_pages": 1,
			"total_count": 1
		}`)
	})

	res, err := client.APISpecRead(8649, 111680)
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
