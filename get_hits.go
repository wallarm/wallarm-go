package wallarm

import (
	"encoding/json"
	"net/http"
)

type (
	// Hits contains operations available on Hit resource.
	Hits interface {
		GetHitsRead(body *GetHitsRead) ([]*HitObject, error)
	}

	// GetHitsRead is the request body for fetching hits from /v1/objects/hit.
	GetHitsRead struct {
		Filter    *GetHitsFilter `json:"filter"`
		Limit     int            `json:"limit"`
		Offset    int            `json:"offset"`
		OrderBy   string         `json:"order_by"`
		OrderDesc bool           `json:"order_desc"`
	}

	// GetHitsFilter defines filter parameters for the hits request.
	// Fields with JSON key prefixed by "!" are negation filters.
	// State and SecurityIssueID are interface{} so they marshal as null when nil.
	GetHitsFilter struct {
		State            interface{} `json:"state"`
		NotType          []string    `json:"!type"`
		Time             [][]int64   `json:"time,omitempty"`
		NotState         string      `json:"!state"`
		SecurityIssueID  interface{} `json:"security_issue_id"`
		Clientid         int         `json:"clientid"`
		NotExperimental  bool        `json:"!experimental"`
		NotAasmEvent     bool        `json:"!aasm_event"`
		NotWallarmScanner bool       `json:"!wallarm_scanner"`
		RequestID        string      `json:"request_id,omitempty"`
	}

	// HitObject represents a single hit returned by the Wallarm API.
	HitObject struct {
		ID                []string      `json:"id"`
		Type              string        `json:"type"`
		IP                string        `json:"ip"`
		Size              int           `json:"size"`
		Statuscode        int           `json:"statuscode"`
		Time              int64         `json:"time"`
		Value             string        `json:"value"`
		Impression        interface{}   `json:"impression"`
		Stamps            []int         `json:"stamps"`
		StampsHash        int           `json:"stamps_hash"`
		Regex             []interface{} `json:"regex"`
		ResponseTime      int           `json:"response_time"`
		RemoteCountry     interface{}   `json:"remote_country"`
		Point             []string      `json:"point"`
		RemotePort        int           `json:"remote_port"`
		Poolid            int           `json:"poolid"`
		IPBlocked         bool          `json:"ip_blocked"`
		Experimental      interface{}   `json:"experimental"`
		WallarmScanner    interface{}   `json:"wallarm_scanner"`
		AttackID          []string      `json:"attackid"`
		BlockStatus       string        `json:"block_status"`
		RequestID         string        `json:"request_id"`
		Datacenter        string        `json:"datacenter"`
		ProxyType         interface{}   `json:"proxy_type"`
		Tor               string        `json:"tor"`
		State             interface{}   `json:"state"`
		KnownAttack       []string      `json:"known_attack"`
		KnownFalse        interface{}   `json:"known_false"`
		Protocol          string        `json:"protocol"`
		AuthProtocol      []string      `json:"auth_protocol"`
		EndpointID        interface{}   `json:"endpoint_id"`
		Path              string        `json:"path"`
		Domain            string        `json:"domain"`
		NodeUUID          []string      `json:"node_uuid"`
		CompromisedLogins []interface{} `json:"compromised_logins"`
		Ebpf              interface{}   `json:"ebpf"`
		AasmEvent         bool          `json:"aasm_event"`
		APISpecViolation  interface{}   `json:"api_spec_violation"`
		APISpecID         interface{}   `json:"api_spec_id"`
	}

	// GetHitsResponse is the top-level API response envelope for /v1/objects/hit.
	GetHitsResponse struct {
		Status int          `json:"status"`
		Body   []*HitObject `json:"body"`
	}
)

// GetHitsRead fetches hits from the Wallarm API matching the provided filter.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) GetHitsRead(body *GetHitsRead) ([]*HitObject, error) {
	uri := "/v1/objects/hit"
	respBody, err := api.makeRequest(http.MethodPost, uri, "hit", body, nil)
	if err != nil {
		return nil, err
	}
	var resp GetHitsResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp.Body, nil
}
