package wallarm

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const activityLogFiltersResponseJSON = `{
  "body": {
    "object_types": [
      {"label": "Client", "value": "client"},
      {"label": "API token", "value": "api_token"}
    ],
    "action_types": [
      {"label": "Create", "value": "create"},
      {"label": "Mode change", "value": "mode_change"}
    ],
    "outcomes": ["success", "failure"],
    "sources": ["unknown", "ui_page", "node", "api_token"]
  }
}`

const activityLogEventsResponseJSON = `{
  "body": {
    "objects": [
      {
        "id": 101,
        "time": 1744531200,
        "action_type": "Create",
        "object_type": "API token",
        "object_type_info": {"label": "API token", "value": "api_token"},
        "outcome": "success",
        "source": "api_token",
        "client_id": 17,
        "initiator": {
          "id": 11,
          "client_id": 17,
          "name": "ops-admin"
        },
        "object_id": "token-uuid-1",
        "object": {
          "name": "CI token",
          "info": "Used by deploy pipeline"
        },
        "description": "Created API token",
        "changed_fields": ["name", "scope"],
        "diff": {"scope": {"old": null, "new": "write"}},
        "state_after_action": {"enabled": true}
      }
    ]
  }
}`

func TestActivityLogEventsGetFilters(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/activity_log/events_get_filters", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(activityLogFiltersResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.ActivityLogEventsGetFilters()
	require.NoError(t, err)

	require.Len(t, resp.Body.ObjectTypes, 2)
	assert.Equal(t, "client", resp.Body.ObjectTypes[0].Value)
	assert.Equal(t, "Mode change", resp.Body.ActionTypes[1].Label)
	assert.Equal(t, []string{"success", "failure"}, resp.Body.Outcomes)
	assert.Equal(t, []string{"unknown", "ui_page", "node", "api_token"}, resp.Body.Sources)
}

func TestActivityLogEventsRead(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/activity_log/17/events", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, []string{"api_token", "user"}, r.URL.Query()["filter[object_types][]"])
		assert.Equal(t, []string{"create", "mode_change"}, r.URL.Query()["filter[action_types][]"])
		assert.Equal(t, []string{"success"}, r.URL.Query()["filter[outcomes][]"])
		assert.Equal(t, []string{"ui_page", "api_token"}, r.URL.Query()["filter[sources][]"])
		assert.Equal(t, []string{"11", "12"}, r.URL.Query()["filter[actor_ids][]"])
		assert.Equal(t, "1744527600", r.URL.Query().Get("filter[time_start]"))
		assert.Equal(t, "1744531200", r.URL.Query().Get("filter[time_end]"))
		assert.Equal(t, "5", r.URL.Query().Get("offset"))
		assert.Equal(t, "25", r.URL.Query().Get("limit"))
		assert.Equal(t, "timestamp", r.URL.Query().Get("order_by"))
		assert.Equal(t, "true", r.URL.Query().Get("order_desc"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(activityLogEventsResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.ActivityLogEventsRead(&ActivityLogEventsRead{
		ClientID: 17,
		Filter: &ActivityLogEventsFilter{
			ObjectTypes: []string{"api_token", "user"},
			ActionTypes: []string{"create", "mode_change"},
			Outcomes:    []string{"success"},
			Sources:     []string{"ui_page", "api_token"},
			ActorIDs:    []int64{11, 12},
			TimeStart:   1744527600,
			TimeEnd:     1744531200,
		},
		Offset:    5,
		Limit:     25,
		OrderBy:   "timestamp",
		OrderDesc: true,
	})
	require.NoError(t, err)

	require.Len(t, resp.Body.Objects, 1)
	event := resp.Body.Objects[0]
	assert.Equal(t, uint64(101), event.ID)
	assert.Equal(t, "Create", event.ActionType)
	assert.Equal(t, "API token", event.ObjectType)
	require.NotNil(t, event.ObjectTypeInfo)
	assert.Equal(t, "API token", event.ObjectTypeInfo.Label)
	assert.Equal(t, "api_token", event.ObjectTypeInfo.Value)
	assert.Equal(t, "success", event.Outcome)
	assert.Equal(t, "api_token", event.Source)
	assert.Equal(t, int64(17), event.ClientID)
	require.NotNil(t, event.Initiator)
	require.NotNil(t, event.Initiator.Name)
	assert.Equal(t, "ops-admin", *event.Initiator.Name)
	require.NotNil(t, event.ObjectID)
	assert.Equal(t, "token-uuid-1", *event.ObjectID)
	assert.Equal(t, []string{"name", "scope"}, event.ChangedFields)
	assert.JSONEq(t, `{"scope":{"old":null,"new":"write"}}`, string(event.Diff))
	assert.JSONEq(t, `{"enabled":true}`, string(event.StateAfterAction))
	assert.Equal(t, "CI token", event.Object.Name)
}

func TestActivityLogEventRead(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/activity_log/23/events/777", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if got := r.URL.RawQuery; got != "" {
			t.Fatalf("unexpected query string %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(activityLogEventsResponseJSON))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(
		UsingBaseURL(server.URL),
		Headers(http.Header{"X-WallarmAPI-Token": []string{"test-token"}}),
	)
	require.NoError(t, err)

	resp, err := client.ActivityLogEventRead(&ActivityLogEventRead{
		ClientID: 23,
		EventID:  777,
	})
	require.NoError(t, err)
	require.Len(t, resp.Body.Objects, 1)
	assert.Equal(t, uint64(101), resp.Body.Objects[0].ID)
}

func TestActivityLogEventsReadToQuery(t *testing.T) {
	req := &ActivityLogEventsRead{
		ClientID: 1,
		Filter: &ActivityLogEventsFilter{
			ObjectTypes: []string{"client"},
			ActionTypes: []string{"delete"},
			Outcomes:    []string{"failure"},
			Sources:     []string{"node"},
			ActorIDs:    []int64{44},
			TimeStart:   100,
			TimeEnd:     200,
		},
		OrderDesc: false,
	}

	values, err := url.ParseQuery(req.toQuery())
	require.NoError(t, err)

	assert.Equal(t, []string{"client"}, values["filter[object_types][]"])
	assert.Equal(t, []string{"delete"}, values["filter[action_types][]"])
	assert.Equal(t, []string{"failure"}, values["filter[outcomes][]"])
	assert.Equal(t, []string{"node"}, values["filter[sources][]"])
	assert.Equal(t, []string{"44"}, values["filter[actor_ids][]"])
	assert.Equal(t, "100", values.Get("filter[time_start]"))
	assert.Equal(t, "200", values.Get("filter[time_end]"))
	assert.Equal(t, "false", values.Get("order_desc"))
}
