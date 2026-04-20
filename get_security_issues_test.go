package wallarm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSecurityIssuesReadUsesPOSTBodyTokenAndClientID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/security_issues", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req GetSecurityIssuesRead
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, 130, req.ClientID)
		assert.Equal(t, "test-token", req.Token)
		require.NotNil(t, req.Filter)
		assert.Equal(t, []string{"open"}, req.Filter.State)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":101,"client_id":130,"severity":"high","state":"open","name":"SQL injection","host":"shop.example.com"}]`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.GetSecurityIssuesRead(&GetSecurityIssuesRead{
		ClientID: 130,
		Limit:    20,
		Filter: &GetSecurityIssuesFilter{
			State: []string{"open"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)
	assert.Equal(t, 101, resp[0].ID)
	assert.Equal(t, "high", resp[0].Severity)
}

func TestGetSecurityIssuesCountUsesPOSTBodyTokenAndClientID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/security_issues/count", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req GetSecurityIssuesCount
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, 130, req.ClientID)
		assert.Equal(t, "test-token", req.Token)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":7}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.GetSecurityIssuesCount(&GetSecurityIssuesCount{ClientID: 130})
	require.NoError(t, err)
	assert.Equal(t, 7, resp.Count)
}

func TestGetSecurityIssueGroupsUsesPOSTBodyTokenAndClientID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/security_issues/groups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var req GetSecurityIssueGroups
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, 130, req.ClientID)
		assert.Equal(t, "test-token", req.Token)
		require.NotNil(t, req.OrderBy)
		assert.Equal(t, "severity", req.OrderBy.Name)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"group_id":"sqli","title":"SQL injection","severity":"high","security_issues_count":3,"states":{"open":2,"closed":1,"marked_false":0}}]`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.GetSecurityIssueGroups(&GetSecurityIssueGroups{
		ClientID: 130,
		OrderBy: &GetSecurityIssuesOrderBy{
			Name:      "severity",
			Direction: "desc",
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)
	assert.Equal(t, "sqli", resp[0].GroupID)
	assert.Equal(t, 3, resp[0].SecurityIssuesCount)
}

func TestGetSecurityIssueUsesGETQueryTokenAndClientID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/security_issues/101", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "130", r.URL.Query().Get("client_id"))
		assert.Equal(t, "test-token", r.URL.Query().Get("token"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":101,"client_id":130,"severity":"high","state":"open","name":"SQL injection","description":"demo issue"}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.GetSecurityIssue(&GetSecurityIssue{
		ID:       101,
		ClientID: 130,
	})
	require.NoError(t, err)
	assert.Equal(t, 101, resp.ID)
	assert.Equal(t, "demo issue", resp.Description)
}

func TestSecurityIssueMethodsRejectNilRequest(t *testing.T) {
	client, err := New()
	require.NoError(t, err)

	respRead, err := client.GetSecurityIssuesRead(nil)
	require.EqualError(t, err, "security issues read request is required")
	assert.Nil(t, respRead)

	respCount, err := client.GetSecurityIssuesCount(nil)
	require.EqualError(t, err, "security issues count request is required")
	assert.Nil(t, respCount)

	respGroups, err := client.GetSecurityIssueGroups(nil)
	require.EqualError(t, err, "security issue groups request is required")
	assert.Nil(t, respGroups)

	respGroupsCount, err := client.GetSecurityIssueGroupsCount(nil)
	require.EqualError(t, err, "security issue groups count request is required")
	assert.Nil(t, respGroupsCount)

	respIssue, err := client.GetSecurityIssue(nil)
	require.EqualError(t, err, "security issue request is required")
	assert.Nil(t, respIssue)
}
