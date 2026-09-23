package wallarm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AttackVectorsDefaultLimit is the page size sent when AttackVectorsByRequestParams.Limit is zero.
const AttackVectorsDefaultLimit = 100

type (
	// AttackVectors contains operations available on the attack-vectors resource.
	AttackVectors interface {
		AttackVectorsByRequest(clientID int, params *AttackVectorsByRequestParams) (*AttackVectorsByRequestResp, error)
	}

	// AttackVectorsByRequestParams is the body of POST
	// /v1/client/{client_id}/attack-vectors/by-request. Limit accepts 1 to 1000.
	AttackVectorsByRequestParams struct {
		RequestID string `json:"request_id"`
		Limit     int    `json:"limit,omitempty"`
	}

	// AttackVectorsByRequestResp is one page of a request's attack vectors.
	// Cursor is set only when HasMore is true.
	AttackVectorsByRequestResp struct {
		Data    []AttackVector `json:"data"`
		HasMore bool           `json:"has_more"`
		Cursor  string         `json:"cursor,omitempty"`
	}

	// AttackVectorCondition is one action condition of a vector. Point mixes
	// strings and numbers (float64 once decoded); Value is empty for "absent".
	AttackVectorCondition struct {
		Type  string        `json:"type"`
		Point []interface{} `json:"point"`
		Value string        `json:"value,omitempty"`
	}

	// AttackVector is one detected attack vector of a request. Point is the
	// wire's JSON-encoded string; ApplicationID is -1 when the request has
	// no application.
	AttackVector struct {
		VectorID           string                  `json:"vector_id"`
		RequestID          string                  `json:"request_id"`
		ClientID           int                     `json:"client_id"`
		ApplicationID      int                     `json:"application_id"`
		Type               string                  `json:"type"`
		Point              string                  `json:"point"`
		Value              string                  `json:"value"`
		Stamps             []int                   `json:"stamps"`
		StampsHash         int                     `json:"stamps_hash"`
		Host               string                  `json:"host"`
		Path               string                  `json:"path"`
		Scheme             string                  `json:"scheme"`
		Method             string                  `json:"method"`
		Protocol           string                  `json:"protocol"`
		RemoteAddr4        string                  `json:"remote_addr4"`
		RemoteAddr6        string                  `json:"remote_addr6"`
		ResponseStatusCode int                     `json:"response_status_code"`
		RequestTime        int                     `json:"request_time"`
		BlockStatus        string                  `json:"block_status"`
		KnownAttack        string                  `json:"known_attack"`
		NodeUUID           []string                `json:"node_uuid"`
		ActionConditions   []AttackVectorCondition `json:"action_conditions"`
	}
)

// AttackVectorsByRequest returns one page of the attack vectors of a request.
// A zero Limit is sent as AttackVectorsDefaultLimit; params is not modified.
func (api *api) AttackVectorsByRequest(clientID int, params *AttackVectorsByRequestParams) (*AttackVectorsByRequestResp, error) {
	if params == nil {
		return nil, fmt.Errorf("attack vectors by request params are required")
	}
	body := *params
	if body.Limit == 0 {
		body.Limit = AttackVectorsDefaultLimit
	}
	uri := fmt.Sprintf("/v1/client/%d/attack-vectors/by-request", clientID)
	respBody, err := api.makeRequest(http.MethodPost, uri, "attack_vectors", &body, nil)
	if err != nil {
		return nil, err
	}
	var resp AttackVectorsByRequestResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
