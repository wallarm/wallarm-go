package wallarm

import (
	"encoding/json"
	"fmt"
)

type (
	OverlimitResSettings interface {
		OverlimitResSettingsUpdate(body *OverlimitResSettingsParams, clientID int) (*OverlimitResSettingsResponse, error)
		OverlimitResSettingsRead(clientID int) (*OverlimitResSettingsResponse, error)
	}

	OverlimitResSettingsParams struct {
		OverlimitTime int    `json:"overlimit_time"`
		Mode          string `json:"mode"`
	}

	OverlimitResSettingsResponse struct {
		Status int                        `json:"status"`
		Body   OverlimitResSettingsParams `json:"body"`
	}
)

// OverlimitResSettingsUpdate changes the global overlimit_res_settings.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) OverlimitResSettingsUpdate(body *OverlimitResSettingsParams, clientID int) (*OverlimitResSettingsResponse, error) {
	uri := fmt.Sprintf("/v2/client/%d/rules/overlimit_res_settings", clientID)
	rawResp, err := api.makeRequest("PUT", uri, "overlimit_res_settings", body, nil)
	if err != nil {
		return nil, err
	}

	var resp OverlimitResSettingsResponse
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// OverlimitResSettingsRead requests the global overlimit_res_settings.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) OverlimitResSettingsRead(clientID int) (*OverlimitResSettingsResponse, error) {
	uri := fmt.Sprintf("/v2/client/%d/rules/overlimit_res_settings", clientID)
	rawResp, err := api.makeRequest("GET", uri, "overlimit_res_settings", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp OverlimitResSettingsResponse
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
