package wallarm

import (
	"encoding/json"
	"net/http"
)

type (
	// Hits contains operations available on Hit resource
	Hits interface {
		HitRead(hitBody *HitReadRequest) ([]*Hit, error)
	}

	// HitReadRequest is a root object for requesting hits.
	// Limit is a number between 0 - 1000
	// Offset is a number
	HitReadRequest struct {
		Filter    *HitFilter `json:"filter"`
		Limit     int        `json:"limit"`
		Offset    int        `json:"offset"`
		OrderBy   string     `json:"order_by"`
		OrderDesc bool       `json:"order_desc"`
	}

	// HitFilter defines the filtering criteria for hit queries.
	// Fields prefixed with Not (json "!field") represent negation filters.
	// Pointer fields with no omitempty are serialized as null when nil.
	HitFilter struct {
		ClientID          int             `json:"clientid"`
		RequestID         string          `json:"request_id,omitempty"`
		ID                []string        `json:"id,omitempty"`
		State             *string         `json:"state"`
		NotType           []string        `json:"!type,omitempty"`
		Time              [][]interface{} `json:"time,omitempty"`
		NotState          string          `json:"!state,omitempty"`
		SecurityIssueID   *int            `json:"security_issue_id"`
		NotExperimental   bool            `json:"!experimental"`
		NotAasmEvent      bool            `json:"!aasm_event"`
		NotWallarmScanner bool            `json:"!wallarm_scanner"`
	}

	// Hit represents a single detection event within a request.
	// A single HTTP request may produce multiple hits if different
	// attack vectors are found at different points in the request.
	Hit struct {
		ID                []string      `json:"id"`
		Type              string        `json:"type"`
		IP                string        `json:"ip"`
		Size              int           `json:"size"`
		StatusCode        int           `json:"statuscode"`
		Time              int           `json:"time"`
		Value             string        `json:"value"`
		Stamps            []int         `json:"stamps"`
		StampsHash        int           `json:"stamps_hash"`
		Regex             []interface{} `json:"regex"`
		ResponseTime      int           `json:"response_time"`
		RemoteCountry     *string       `json:"remote_country"`
		Point             []interface{} `json:"point"`
		RemotePort        int           `json:"remote_port"`
		PoolID            int           `json:"poolid"`
		IPBlocked         bool          `json:"ip_blocked"`
		AttackID          []string      `json:"attackid"`
		BlockStatus       string        `json:"block_status"`
		RequestID         string        `json:"request_id"`
		Datacenter        string        `json:"datacenter"`
		ProxyType         *string       `json:"proxy_type"`
		Tor               string        `json:"tor"`
		State             *string       `json:"state"`
		KnownAttack       []string      `json:"known_attack"`
		KnownFalse        *string       `json:"known_false"`
		Protocol          string        `json:"protocol"`
		AuthProtocol      []string      `json:"auth_protocol"`
		EndpointID        *string       `json:"endpoint_id"`
		Path              string        `json:"path"`
		Domain            string        `json:"domain"`
		NodeUUID          []string      `json:"node_uuid"`
		CompromisedLogins []string      `json:"compromised_logins"`
		Ebpf              *bool         `json:"ebpf"`
		AasmEvent         bool          `json:"aasm_event"`
		ApiSpecViolation  *string       `json:"api_spec_violation"`
		ApiSpecID         *string       `json:"api_spec_id"`
	}

	// hitReadRawResponse wraps the API response envelope for hits.
	// The Wallarm /v1/objects/hit endpoint returns {"status": int, "body": [...]}.
	hitReadRawResponse struct {
		Status int    `json:"status"`
		Body   []*Hit `json:"body"`
	}
)

// HitRead queries the Wallarm API for hits matching the given filter criteria.
// Returns a slice of Hit objects. A single request_id may yield multiple hits
// when different attack vectors are detected at different request points.
func (api *api) HitRead(hitBody *HitReadRequest) ([]*Hit, error) {
	uri := "/v1/objects/hit"
	respBody, err := api.makeRequest(http.MethodPost, uri, "hit", hitBody,
		map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var resp hitReadRawResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp.Body, nil
}
