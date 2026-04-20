package wallarm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPListRead_GroupedEntries(t *testing.T) {
	setup()
	defer teardown()

	requests := 0
	mux.HandleFunc("/v1/blocklist/clients/17/groups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "block", r.URL.Query().Get("filter[list]"))
		assert.ElementsMatch(t, []string{"subnet", "proxy_type", "datacenter", "location"}, r.URL.Query()["filter[rule_type][]"])
		assert.Equal(t, "1", r.URL.Query().Get("limit"))
		assert.Equal(t, fmt.Sprintf("%d", requests), r.URL.Query().Get("offset"))

		w.Header().Set("content-type", "application/json")
		if requests == 0 {
			fmt.Fprint(w, `{"body":{"objects":[{"id":7,"client_id":17,"rule_type":"subnet","list":"block","created_at":1770000000,"expired_at":0,"application_ids":[0],"reason":"manual block","author_user_id":11,"values":["1.2.3.4"],"status":"enabled"}]}}`)
		} else if requests == 1 {
			fmt.Fprint(w, `{"body":{"objects":[{"id":8,"client_id":17,"rule_type":"location","list":"block","created_at":1770000100,"expired_at":1771000000,"application_ids":[12],"reason":"geo block","author_user_id":12,"values":["RU"],"status":"enabled"}]}}`)
		} else {
			fmt.Fprint(w, `{"body":{"objects":[]}}`)
		}
		requests++
	})

	actual, err := client.IPListRead(DenylistType, 17, 1)
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Equal(t, []IPRule{
		{
			ID:             7,
			ClientID:       17,
			RuleType:       "subnet",
			List:           "block",
			CreatedAt:      1770000000,
			ExpiredAt:      0,
			ApplicationIDs: []int{0},
			Reason:         "manual block",
			AuthorUserID:   11,
			Values:         []string{"1.2.3.4"},
			Status:         "enabled",
		},
		{
			ID:             8,
			ClientID:       17,
			RuleType:       "location",
			List:           "block",
			CreatedAt:      1770000100,
			ExpiredAt:      1771000000,
			ApplicationIDs: []int{12},
			Reason:         "geo block",
			AuthorUserID:   12,
			Values:         []string{"RU"},
			Status:         "enabled",
		},
	}, actual)
	assert.Equal(t, 3, requests)
}

func TestIPListSearch_UsesBracketedFilterQuery(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/blocklist/clients/42/groups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, []string{"location"}, r.URL.Query()["filter[rule_type][]"])
		assert.Equal(t, "allow", r.URL.Query().Get("filter[list]"))
		assert.Equal(t, "RU", r.URL.Query().Get("filter[query]"))
		assert.Empty(t, r.URL.Query().Get("filter"))
		assert.Equal(t, "1", r.URL.Query().Get("limit"))
		assert.Equal(t, "0", r.URL.Query().Get("offset"))

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"body":{"objects":[{"id":9,"client_id":42,"rule_type":"location","list":"allow","created_at":1770000200,"expired_at":0,"application_ids":[0],"reason":"allow geo","author_user_id":13,"values":["RU"],"status":"enabled"}]}}`)
	})

	actual, err := client.IPListSearch(AllowlistType, 42, "location", "RU")
	require.NoError(t, err)
	require.Len(t, actual, 1)
	assert.Equal(t, IPRule{
		ID:             9,
		ClientID:       42,
		RuleType:       "location",
		List:           "allow",
		CreatedAt:      1770000200,
		ExpiredAt:      0,
		ApplicationIDs: []int{0},
		Reason:         "allow geo",
		AuthorUserID:   13,
		Values:         []string{"RU"},
		Status:         "enabled",
	}, actual[0])
}

func TestIPListCreate(t *testing.T) {
	setup()
	defer teardown()

	expected := AccessRuleCreateRequest{
		List:           GraylistType,
		Force:          true,
		Reason:         "temporary review",
		ApplicationIDs: []int{0, 12},
		ExpiredAt:      1772000000,
		Rules: []AccessRuleEntry{{
			RulesType: "subnet",
			Values:    []string{"10.0.0.0/8"},
		}},
	}

	mux.HandleFunc("/v1/blocklist/clients/23/access_rules", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var actual AccessRuleCreateRequest
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, expected, actual)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status":200}`)
	})

	err := client.IPListCreate(23, expected)
	require.NoError(t, err)
}

func TestIPListDelete(t *testing.T) {
	setup()
	defer teardown()

	expected := AccessRuleDeleteRequest{
		Rules: []AccessRuleDeleteEntry{
			{RuleType: "subnet", IDs: []int{7, 8}},
			{RuleType: "location", IDs: []int{9}},
		},
	}

	mux.HandleFunc("/v1/blocklist/clients/23/groups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method, "Expected method 'DELETE', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var actual AccessRuleDeleteRequest
		require.NoError(t, json.Unmarshal(body, &actual))
		assert.Equal(t, expected, actual)

		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status":200}`)
	})

	err := client.IPListDelete(23, expected.Rules)
	require.NoError(t, err)
}
