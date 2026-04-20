package wallarm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const attackResponseJSON = `{
  "status": 200,
  "cursor": "cursor-2",
  "body": [
    {
      "id": ["attacks_production_2_202603_v_1", "AM5CP50BU-M_02QxawVC"],
      "attackid": "attacks_production_2_202603_v_1:AM5CP50BU-M_02QxawVC",
      "clientid": 2,
      "domain": "node-data.audit.wallarm.com",
      "poolid": -1,
      "method": "GET",
      "parameter": "GET_filename_value",
      "path": "/log/view",
      "type": "ptrav",
      "first_time": 1774882806,
      "last_time": 1774892699,
      "hits": 111,
      "ip_count": 1,
      "ip_top": [{"ip": "161.97.68.155", "count": 111, "country": "FR"}],
      "country_top": [{"country": "FR", "count": 111}],
      "statuscodes": [200, 404],
      "vulnid": null,
      "threat": 90,
      "target": "server",
      "experimental": null,
      "recheck_status": "unrecheckable",
      "vectors_count": 43,
      "block_status": "monitored",
      "state": null
    },
    {
      "id": ["attacks_production_2_202603_v_1", "BN6DQ61CU-N_13RybyWD"],
      "attackid": "attacks_production_2_202603_v_1:BN6DQ61CU-N_13RybyWD",
      "clientid": 2,
      "domain": "api.audit.wallarm.com",
      "poolid": 5,
      "method": "POST",
      "parameter": "POST_body_value",
      "path": "/v1/login",
      "type": "sqli",
      "first_time": 1774880000,
      "last_time": 1774885000,
      "hits": 42,
      "ip_count": 3,
      "ip_top": [{"ip": "10.0.0.1", "count": 30, "country": "US"}],
      "country_top": [{"country": "US", "count": 30}],
      "statuscodes": [403],
      "vulnid": null,
      "threat": 80,
      "target": "server",
      "experimental": null,
      "recheck_status": "in_progress",
      "vectors_count": 5,
      "block_status": "blocked",
      "state": null
    }
  ]
}`

const attackCountResponseJSON = `{
  "status": 200,
  "body": {
    "attacks": 621,
    "hits": 61818.0,
    "ips": 16
  }
}`

const attackIPResponseJSON = `{
  "status": 200,
  "body": ["11.22.33.44", "11.22.33.45"]
}`

const hitResponseJSON = `{
  "status": 200,
  "body": [
    {
      "id": ["hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"],
      "attackid": ["attacks_test_130_202604_v_1", "3z5nB4P16C7u8Q"],
      "type": "sqli",
      "ip": "11.22.33.45",
      "statuscode": 403,
      "time": 1743531300,
      "value": "' UNION SELECT password FROM users --",
      "remote_country": "US",
      "block_status": "blocked",
      "request_id": "core-2252-login",
      "path": "/api/v1/login",
      "domain": "shop.example.com"
    }
  ]
}`

const hitDetailsResponseJSON = `{
  "status": 200,
  "body": [
    {
      "id": ["hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"],
      "raw": {
        "method": "POST",
        "uri": "/api/v1/login",
        "proto": "HTTP/1.1",
        "headers": {
          "host": ["shop.example.com"],
          "content-length": 123
        }
      }
    }
  ]
}`

func TestAttackRead(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/objects/attack", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req AttackReadRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, 2, req.Limit)
		assert.Equal(t, "last_time", req.OrderBy)
		assert.Equal(t, "cursor-1", req.Cursor)
		assert.True(t, req.Paging)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(attackResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.AttackRead(&AttackReadRequest{
		Filter: &AttackFilter{
			ClientID: []int{2},
			Time:     [][]interface{}{{1774800000, nil}},
		},
		Limit:     2,
		Offset:    0,
		OrderBy:   "last_time",
		OrderDesc: true,
		Cursor:    "cursor-1",
		Paging:    true,
	})
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Status)
	assert.Equal(t, "cursor-2", resp.Cursor)
	require.Len(t, resp.Body, 2)

	atk := resp.Body[0]
	assert.Equal(t, "attacks_production_2_202603_v_1:AM5CP50BU-M_02QxawVC", atk.AttackID)
	assert.Equal(t, "node-data.audit.wallarm.com", atk.Domain)
	assert.Equal(t, "/log/view", atk.Path)
	assert.Equal(t, "ptrav", atk.Type)
	assert.Equal(t, "GET", atk.Method)
	assert.Equal(t, 111, atk.Hits)
	assert.Equal(t, 1, atk.IPCount)
	assert.Equal(t, 90, atk.Threat)
	assert.Equal(t, 1774882806, atk.FirstTime)
	assert.Equal(t, 1774892699, atk.LastTime)
	assert.Equal(t, 2, atk.ClientID)
	assert.Equal(t, -1, atk.PoolID)
	assert.Equal(t, "server", atk.Target)
	assert.Equal(t, 43, atk.VectorsCount)
	assert.Equal(t, "monitored", atk.BlockStatus)
	assert.Equal(t, "GET_filename_value", atk.Parameter)

	atk2 := resp.Body[1]
	assert.Equal(t, "sqli", atk2.Type)
	assert.Equal(t, "api.audit.wallarm.com", atk2.Domain)
	assert.Equal(t, 42, atk2.Hits)
	assert.Equal(t, "blocked", atk2.BlockStatus)
}

func TestAttackCount(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/objects/attack/count", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(attackCountResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.AttackCount(&AttackCountRequest{
		Filter: &AttackCountFilter{
			ClientID: []int{2},
			Time:     [][]interface{}{{1774800000, nil}},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Status)
	assert.Equal(t, 621, resp.Body.Attacks)
	assert.Equal(t, 61818.0, resp.Body.Hits)
	assert.Equal(t, 16, resp.Body.IPs)
}

func TestAttackIP(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/objects/attack/ip", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req AttackIPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, []int{130}, req.Filter.ClientID)
		assert.Equal(t, []string{"11.22.33.44"}, req.Filter.IP)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(attackIPResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.AttackIP(&AttackIPRequest{
		Filter: &AttackCountFilter{
			ClientID: []int{130},
			IP:       []string{"11.22.33.44"},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Status)
	assert.Equal(t, []string{"11.22.33.44", "11.22.33.45"}, resp.Body)
}

func TestHitReadWithAttackIDTupleFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/objects/hit", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, float64(3), req["limit"])
		assert.Equal(t, "time", req["order_by"])
		assert.Equal(t, true, req["order_desc"])

		filter := req["filter"].(map[string]interface{})
		attackIDs := filter["attackid"].([]interface{})
		firstTuple := attackIDs[0].([]interface{})
		assert.Equal(t, "attacks_test_130_202604_v_1", firstTuple[0])
		assert.Equal(t, "3z5nB4P16C7u8Q", firstTuple[1])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(hitResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.HitRead(&HitReadRequest{
		Filter: &HitFilter{
			ClientID: 130,
			AttackID: [][]string{{"attacks_test_130_202604_v_1", "3z5nB4P16C7u8Q"}},
		},
		Limit:     3,
		OrderBy:   "time",
		OrderDesc: true,
	})
	require.NoError(t, err)

	require.Len(t, resp, 1)
	assert.Equal(t, []string{"hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"}, resp[0].ID)
	assert.Equal(t, "11.22.33.45", resp[0].IP)
	assert.Equal(t, 403, resp[0].StatusCode)
}

func TestHitDetails(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/objects/hit/details", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req HitDetailsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, []string{"hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"}, req.Filter.ID)
		assert.Equal(t, "raw", req.Returns)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(hitDetailsResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.HitDetails(&HitDetailsRequest{
		Filter:  &HitFilter{ID: []string{"hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"}},
		Returns: "raw",
	})
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Status)
	require.Len(t, resp.Body, 1)
	assert.Equal(t, "POST", resp.Body[0].Raw.Method)
	assert.Equal(t, "/api/v1/login", resp.Body[0].Raw.URI)
	assert.Equal(t, "HTTP/1.1", resp.Body[0].Raw.Proto)
}

func TestHitRaw(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/objects/hit/raw", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req HitRawRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, []string{"hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"}, req.Filter.ID)

		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte("raw-hit-payload"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.HitRaw(&HitRawRequest{
		Filter: &HitFilter{ID: []string{"hits_test_130_202604_v_1", "Vr0ix9J6f5Wgk3dS"}},
	})
	require.NoError(t, err)

	assert.Equal(t, []byte("raw-hit-payload"), resp)
}

func TestAttackMethodsRejectNilRequest(t *testing.T) {
	client, err := New()
	require.NoError(t, err)

	respRead, err := client.AttackRead(nil)
	require.EqualError(t, err, "attack read request is required")
	assert.Nil(t, respRead)

	respCount, err := client.AttackCount(nil)
	require.EqualError(t, err, "attack count request is required")
	assert.Nil(t, respCount)

	respIP, err := client.AttackIP(nil)
	require.EqualError(t, err, "attack ip request is required")
	assert.Nil(t, respIP)

	respDetails, err := client.HitDetails(nil)
	require.EqualError(t, err, "hit details request is required")
	assert.Nil(t, respDetails)

	respRaw, err := client.HitRaw(nil)
	require.EqualError(t, err, "hit raw request is required")
	assert.Nil(t, respRaw)
}
