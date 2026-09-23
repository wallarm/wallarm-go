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

const attackVectorsByRequestFixture = `{
	"data": [
		{
			"vector_id": "8f6bcf31a1dac9939fb1e0d641658cbc667f6c69d512e05a777be62145e95080",
			"request_id": "req-1",
			"client_id": 8649,
			"application_id": -1,
			"type": "xss",
			"point": "[\"header\",\"XSS\"]",
			"value": "<script>alert(1)</script>",
			"stamps": [454, 7],
			"stamps_hash": 4182624376,
			"host": "127.0.0.1:8080",
			"path": "/a/b/c",
			"normalized_host": "127.0.0.1:8080",
			"normalized_path": "/a/b/c",
			"scheme": "http",
			"method": "GET",
			"protocol": "HTTP/1.1",
			"remote_addr4": "152.32.174.171",
			"remote_addr6": "",
			"response_status_code": 502,
			"request_time": 1790142071217,
			"block_status": "monitored",
			"known_attack": "generic_rce,CVE-2017-9841",
			"node_uuid": ["19cbe661-7cb0-4348-8b79-88ec6f4731b9"],
			"state": "NOT FILLED",
			"action_conditions": [
				{"point": ["instance"], "type": "equal", "value": "13"},
				{"point": ["header", "HOST"], "type": "iequal", "value": "127.0.0.1:8080"},
				{"point": ["path", 0], "type": "equal", "value": "a"},
				{"point": ["action_ext"], "type": "absent"}
			]
		}
	],
	"has_more": true,
	"cursor": "MTc5MDE0MjA3MTIxN3x2ZWN0b3J8cmVxLTE="
}`

func TestAttackVectorsByRequest(t *testing.T) {
	setup()
	defer teardown()

	var sent map[string]interface{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		raw, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		sent = nil
		assert.NoError(t, json.Unmarshal(raw, &sent))
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, attackVectorsByRequestFixture)
	}
	mux.HandleFunc("/v1/client/8649/attack-vectors/by-request", handler)

	params := &AttackVectorsByRequestParams{RequestID: "req-1"}
	resp, err := client.AttackVectorsByRequest(8649, params)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"request_id": "req-1", "limit": float64(100)}, sent)
	assert.Equal(t, 0, params.Limit)
	if !assert.NotNil(t, resp) {
		return
	}

	assert.True(t, resp.HasMore)
	assert.Equal(t, "MTc5MDE0MjA3MTIxN3x2ZWN0b3J8cmVxLTE=", resp.Cursor)
	if !assert.Len(t, resp.Data, 1) {
		return
	}
	v := resp.Data[0]
	assert.Equal(t, "8f6bcf31a1dac9939fb1e0d641658cbc667f6c69d512e05a777be62145e95080", v.VectorID)
	assert.Equal(t, "req-1", v.RequestID)
	assert.Equal(t, 8649, v.ClientID)
	assert.Equal(t, -1, v.ApplicationID)
	assert.Equal(t, "xss", v.Type)
	assert.Equal(t, `["header","XSS"]`, v.Point)
	assert.Equal(t, "<script>alert(1)</script>", v.Value)
	assert.Equal(t, []int{454, 7}, v.Stamps)
	assert.Equal(t, 4182624376, v.StampsHash)
	assert.Equal(t, "127.0.0.1:8080", v.Host)
	assert.Equal(t, "/a/b/c", v.Path)
	assert.Equal(t, "http", v.Scheme)
	assert.Equal(t, "GET", v.Method)
	assert.Equal(t, "HTTP/1.1", v.Protocol)
	assert.Equal(t, "152.32.174.171", v.RemoteAddr4)
	assert.Equal(t, "", v.RemoteAddr6)
	assert.Equal(t, 502, v.ResponseStatusCode)
	assert.Equal(t, 1790142071217, v.RequestTime)
	assert.Equal(t, "monitored", v.BlockStatus)
	assert.Equal(t, "generic_rce,CVE-2017-9841", v.KnownAttack)
	assert.Equal(t, []string{"19cbe661-7cb0-4348-8b79-88ec6f4731b9"}, v.NodeUUID)
	assert.Equal(t, []AttackVectorCondition{
		{Type: "equal", Point: []interface{}{"instance"}, Value: "13"},
		{Type: "iequal", Point: []interface{}{"header", "HOST"}, Value: "127.0.0.1:8080"},
		{Type: "equal", Point: []interface{}{"path", float64(0)}, Value: "a"},
		{Type: "absent", Point: []interface{}{"action_ext"}},
	}, v.ActionConditions)
}

func TestAttackVectorsByRequestExplicitLimit(t *testing.T) {
	setup()
	defer teardown()

	var sent map[string]interface{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.NoError(t, json.Unmarshal(raw, &sent))
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"data":[],"has_more":false}`)
	}
	mux.HandleFunc("/v1/client/8649/attack-vectors/by-request", handler)

	resp, err := client.AttackVectorsByRequest(8649, &AttackVectorsByRequestParams{RequestID: "req-1", Limit: 2})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"request_id": "req-1", "limit": float64(2)}, sent)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Empty(t, resp.Data)
	assert.False(t, resp.HasMore)
	assert.Empty(t, resp.Cursor)
}

func TestAttackVectorsByRequestAPIError(t *testing.T) {
	setup()
	defer teardown()

	calls := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"status":400,"body":{"limit":{"error":"limit must be between 1 and 1000"}}}`)
	}
	mux.HandleFunc("/v1/client/8649/attack-vectors/by-request", handler)

	resp, err := client.AttackVectorsByRequest(8649, &AttackVectorsByRequestParams{RequestID: "req-1", Limit: 1001})
	assert.Nil(t, resp)
	var apiErr *APIError
	if assert.True(t, errors.As(err, &apiErr)) {
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	}
	assert.Equal(t, 1, calls)
}

func TestAttackVectorsByRequestNilParams(t *testing.T) {
	setup()
	defer teardown()

	calls := 0
	mux.HandleFunc("/v1/client/8649/attack-vectors/by-request", func(w http.ResponseWriter, r *http.Request) {
		calls++
	})

	resp, err := client.AttackVectorsByRequest(8649, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, 0, calls)
}
