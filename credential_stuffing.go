package wallarm

import (
	"encoding/json"
	"fmt"
)

type (
	// CredentialStuffingConfigs contains operations for reading credential stuffing configurations.
	CredentialStuffingConfigs interface {
		CredentialStuffingConfigsRead(clientID int) ([]ActionBody, error)
	}

	// CredentialStuffingConfigsResp is the response from
	// GET /v4/clients/{clientID}/credential_stuffing/configs.
	// Configs are split into "default" and "custom" buckets.
	CredentialStuffingConfigsResp struct {
		Status int `json:"status"`
		Body   struct {
			Default []ActionBody `json:"default"`
			Custom  []ActionBody `json:"custom"`
		} `json:"body"`
	}
)

// CredentialStuffingConfigsRead fetches all credential stuffing configs for a client.
// API: GET /v4/clients/{clientID}/credential_stuffing/configs
func (api *api) CredentialStuffingConfigsRead(clientID int) ([]ActionBody, error) {
	uri := fmt.Sprintf("/v4/clients/%d/credential_stuffing/configs", clientID)
	respBody, err := api.makeRequest("GET", uri, "credential_stuffing", nil, nil)
	if err != nil {
		return nil, err
	}
	var resp CredentialStuffingConfigsResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("CredentialStuffingConfigsRead: failed to parse response: %w", err)
	}
	// Merge both buckets into a single slice.
	configs := make([]ActionBody, 0, len(resp.Body.Default)+len(resp.Body.Custom))
	configs = append(configs, resp.Body.Default...)
	configs = append(configs, resp.Body.Custom...)
	return configs, nil
}
