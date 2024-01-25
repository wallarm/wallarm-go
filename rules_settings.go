package wallarm

import (
	"encoding/json"
	"fmt"
)

type (
	RuleSettings interface {
		RulesSettingsRead(clientID int) (*RulesSettingsResponse, error)
		RulesSettingsUpdate(params *RuleSettingsParams, clientID int) (*RulesSettingsResponse, error)
	}

	RulesSettingsResponse struct {
		Status int                        `json:"status"`
		Body   *RulesSettingsResponseBody `json:"body"`
	}

	RulesSettingsResponseBody struct {
		ClientId int `json:"clientid"`
		*RuleSettingsParams
	}

	RuleSettingsParams struct {
		MinLomFormat            *int    `json:"min_lom_format,omitempty"`
		MaxLomFormat            *int    `json:"max_lom_format,omitempty"`
		MaxLomSize              *int    `json:"max_lom_size,omitempty"`
		LomDisabled             *bool   `json:"lom_disabled,omitempty"`
		LomCompilationDelay     *int    `json:"lom_compilation_delay,omitempty"`
		RulesSnapshotEnabled    *bool   `json:"rules_snapshot_enabled,omitempty"`
		RulesSnapshotMaxCount   *int    `json:"rules_snapshot_max_count,omitempty"`
		RulesManipulationLocked *bool   `json:"rules_manipulation_locked,omitempty"`
		HeavyLom                *bool   `json:"heavy_lom,omitempty"`
		DataInS3                *bool   `json:"data_in_s3,omitempty"`
		ParametersCountWeight   *int    `json:"parameters_count_weight,omitempty"`
		PathVariativityWeight   *int    `json:"path_variativity_weight,omitempty"`
		PiiWeight               *int    `json:"pii_weight,omitempty"`
		RequestContentWeight    *int    `json:"request_content_weight,omitempty"`
		OpenVulnsWeight         *int    `json:"open_vulns_weight,omitempty"`
		SerializedDataWeight    *int    `json:"serialized_data_weight,omitempty"`
		RiskScoreAlgo           *string `json:"risk_score_algo,omitempty"`
		PiiFallback             *bool   `json:"pii_fallback,omitempty"`
	}
)

func (api *api) RulesSettingsRead(clientID int) (*RulesSettingsResponse, error) {
	uri := fmt.Sprintf("/v2/client/%d/rules/settings", clientID)

	rawResponse, err := api.makeRequest("GET", uri, "rules_settings", nil)
	if err != nil {
		return nil, err
	}
	var response RulesSettingsResponse
	if err = json.Unmarshal(rawResponse, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (api *api) RulesSettingsUpdate(params *RuleSettingsParams, clientID int) (*RulesSettingsResponse, error) {
	uri := fmt.Sprintf("/v2/client/%d/rules/settings", clientID)

	rawResponse, err := api.makeRequest("PUT", uri, "rules_settings", params)
	if err != nil {
		return nil, err
	}
	var response RulesSettingsResponse
	if err = json.Unmarshal(rawResponse, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
