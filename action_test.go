package wallarm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHintRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects/hint", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": [{
				"id": 100,
				"actionid": 200,
				"type": "wallarm_mode",
				"action": [{"type": "equal", "point": ["header", "HOST"], "value": "example.com"}]
			}]
		}`)
	})

	res, err := client.HintRead(&HintRead{
		Limit:  1,
		Filter: &HintFilter{Clientid: []int{8649}, ID: []int{100}},
	})
	assert.NoError(t, err)
	assert.Len(t, *res.Body, 1)
	assert.Equal(t, 100, (*res.Body)[0].ID)
	assert.Equal(t, "wallarm_mode", (*res.Body)[0].Type)
}

func TestHintCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects/hint/create", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"id": 500,
				"actionid": 300,
				"type": "vpatch",
				"action": []
			}
		}`)
	})

	action := []ActionDetails{}
	res, err := client.HintCreate(&ActionCreate{
		Type:     "vpatch",
		Clientid: 8649,
		Action:   &action,
	})
	assert.NoError(t, err)
	assert.Equal(t, 500, res.Body.ID)
	assert.Equal(t, 300, res.Body.ActionID)
}

func TestHintDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects/hint/delete", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status": 200, "body": [{"id": 500, "actionid": 300, "type": "wallarm_mode"}]}`)
	})

	resp, err := client.HintDelete(&HintDelete{
		Filter: &HintDeleteFilter{
			Clientid: []int{8649},
			ID:       []int{500},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.Status)
	assert.Len(t, resp.Body, 1)
	assert.Equal(t, 500, resp.Body[0].ID)
}

// Verifies the no-op contract documented on HintDelete.
func TestHintDelete_NoOpReturnsEmptyBody(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects/hint/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status": 200, "body": []}`)
	})

	resp, err := client.HintDelete(&HintDelete{
		Filter: &HintDeleteFilter{
			Clientid: []int{8649},
			ID:       []int{999},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.Status)
	assert.Empty(t, resp.Body)
}

func TestActionReadByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/action/42", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"id": 42,
				"clientid": 8649,
				"conditions": [{"type": "equal", "point": ["header", "HOST"], "value": "test.com"}]
			}
		}`)
	})

	res, err := client.ActionReadByID(42)
	assert.NoError(t, err)
	assert.Equal(t, 42, res.ID)
}

func TestActionReadByHitID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects/action/by_hit", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"conditions": [{"type": "iequal", "point": ["header", "HOST"], "value": "example.com"}],
				"clientid": 8649
			}
		}`)
	})

	res, err := client.ActionReadByHitID([]string{"abc123"})
	assert.NoError(t, err)
	assert.Equal(t, 8649, res.Body.Clientid)
}

func TestHintUpdateV3(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/hint/500", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"id": 500,
				"actionid": 300,
				"type": "vpatch",
				"action": []
			}
		}`)
	})

	comment := "updated"
	res, err := client.HintUpdateV3(500, &HintUpdateV3Params{
		Comment: &comment,
	})
	assert.NoError(t, err)
	assert.Equal(t, 500, res.Body.ID)
}

func TestActionList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects/action", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": [{
				"id": 1,
				"conditions": [{"type": "equal", "point": ["header", "HOST"], "value": "test.com"}]
			}]
		}`)
	})

	res, err := client.ActionList(&ActionListParams{
		Filter: &ActionListFilter{Clientid: []int{8649}},
		Limit:  10,
	})
	assert.NoError(t, err)
	assert.Len(t, res.Body, 1)
	assert.Equal(t, 1, res.Body[0].ID)
}

// TestActionCreate_RateLimitFieldsAcceptZero is a regression test for the
// silent-zero-drop bug fixed in v0.12.1. With non-pointer int + omitempty,
// json.Marshal would drop a literal zero from the wire, which the API then
// rejected as "can't be blank". Switching to *int allows callers to send 0.
func TestActionCreate_RateLimitFieldsAcceptZero(t *testing.T) {
	zero := 0
	body := ActionCreate{
		Type:          "rate_limit",
		Rate:          &zero,
		Burst:         &zero,
		Delay:         &zero,
		OverlimitTime: &zero,
	}
	out, err := json.Marshal(body)
	assert.NoError(t, err)
	got := string(out)
	for _, want := range []string{`"rate":0`, `"burst":0`, `"delay":0`, `"overlimit_time":0`} {
		if !strings.Contains(got, want) {
			t.Errorf("expected JSON to contain %q, got %s", want, got)
		}
	}
}

// TestActionCreate_RateLimitFieldsOmitNil verifies that nil pointers continue
// to be omitted from the wire payload, so callers can opt out cleanly.
func TestActionCreate_RateLimitFieldsOmitNil(t *testing.T) {
	body := ActionCreate{Type: "rate_limit"}
	out, err := json.Marshal(body)
	assert.NoError(t, err)
	got := string(out)
	for _, unwanted := range []string{`"rate"`, `"burst"`, `"delay"`, `"overlimit_time"`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("expected JSON to omit %q for nil pointer, got %s", unwanted, got)
		}
	}
}
