package wallarm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleSettingsRead(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expected method 'GET', got %s", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"clientid": 123,
				"min_lom_format": 53,
        "max_lom_format": 54,
				"max_lom_size": 10000,
				"lom_disabled": true,
				"lom_compilation_delay": 30,
        "rules_snapshot_enabled": true,
        "rules_snapshot_max_count": 45,
        "rules_manipulation_locked": false,
				"heavy_lom": true,
        "data_in_s3": false,
        "parameters_count_weight": 1,
        "path_variativity_weight": 2,
        "pii_weight": 3,
        "request_content_weight": 4,
        "open_vulns_weight": 5,
        "serialized_data_weight": 6,
        "risk_score_algo": "maximum",
        "pii_fallback": false
			}
		}`)
	}

	mux.HandleFunc("/v2/client/123/rules/settings", handler)
	actual, err := client.RulesSettingsRead(123)

	minLomFormat := 53
	maxLomFormat := 54
	maxLomSize := 10000
	lomDisabled := true
	lomCompilationDelay := 30
	rulesSnapshotEnabled := true
	rulesSnapshotMaxCount := 45
	rulesManipulationLocked := false
	heavyLom := true
	dataInS3 := false
	parametersCountWeight := 1
	pathVariativityWeight := 2
	piiWeight := 3
	requestContentWeight := 4
	openVulnsWeight := 5
	serializedDataWeight := 6
	riskScoreAlgo := "maximum"
	piiFallback := false

	expected := RulesSettingsResponse{
		Status: 200,
		Body: &RulesSettingsResponseBody{
			ClientId: 123,
			RuleSettingsParams: &RuleSettingsParams{
				MinLomFormat:            &minLomFormat,
				MaxLomFormat:            &maxLomFormat,
				MaxLomSize:              &maxLomSize,
				LomDisabled:             &lomDisabled,
				LomCompilationDelay:     &lomCompilationDelay,
				RulesSnapshotEnabled:    &rulesSnapshotEnabled,
				RulesSnapshotMaxCount:   &rulesSnapshotMaxCount,
				RulesManipulationLocked: &rulesManipulationLocked,
				HeavyLom:                &heavyLom,
				DataInS3:                &dataInS3,
				ParametersCountWeight:   &parametersCountWeight,
				PathVariativityWeight:   &pathVariativityWeight,
				PiiWeight:               &piiWeight,
				RequestContentWeight:    &requestContentWeight,
				OpenVulnsWeight:         &openVulnsWeight,
				SerializedDataWeight:    &serializedDataWeight,
				RiskScoreAlgo:           &riskScoreAlgo,
				PiiFallback:             &piiFallback,
			},
		},
	}

	if assert.NoError(t, err) {
		assert.Equal(t, expected, *actual)
	}
}

func TestRuleSettingsUpdate(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT", "Expected method 'PUT', got %s", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"clientid": 123,
				"min_lom_format": 53,
        "max_lom_format": 74,
				"max_lom_size": 10000,
				"lom_disabled": true,
				"lom_compilation_delay": 30,
        "rules_snapshot_enabled": true,
        "rules_snapshot_max_count": 45,
        "rules_manipulation_locked": false,
				"heavy_lom": true,
        "data_in_s3": false,
        "parameters_count_weight": 1,
        "path_variativity_weight": 2,
        "pii_weight": 3,
        "request_content_weight": 4,
        "open_vulns_weight": 5,
        "serialized_data_weight": 6,
        "risk_score_algo": "maximum",
        "pii_fallback": false
			}
		}`)
	}

	mux.HandleFunc("/v2/client/123/rules/settings", handler)

	maxLomFormat := 74
	params := &RuleSettingsParams{
		MaxLomFormat: &maxLomFormat,
	}

	actual, err := client.RulesSettingsUpdate(params, 123)

	minLomFormat := 53
	maxLomSize := 10000
	lomDisabled := true
	lomCompilationDelay := 30
	rulesSnapshotEnabled := true
	rulesSnapshotMaxCount := 45
	rulesManipulationLocked := false
	heavyLom := true
	dataInS3 := false
	parametersCountWeight := 1
	pathVariativityWeight := 2
	piiWeight := 3
	requestContentWeight := 4
	openVulnsWeight := 5
	serializedDataWeight := 6
	riskScoreAlgo := "maximum"
	piiFallback := false

	expected := RulesSettingsResponse{
		Status: 200,
		Body: &RulesSettingsResponseBody{
			ClientId: 123,
			RuleSettingsParams: &RuleSettingsParams{
				MinLomFormat:            &minLomFormat,
				MaxLomFormat:            &maxLomFormat,
				MaxLomSize:              &maxLomSize,
				LomDisabled:             &lomDisabled,
				LomCompilationDelay:     &lomCompilationDelay,
				RulesSnapshotEnabled:    &rulesSnapshotEnabled,
				RulesSnapshotMaxCount:   &rulesSnapshotMaxCount,
				RulesManipulationLocked: &rulesManipulationLocked,
				HeavyLom:                &heavyLom,
				DataInS3:                &dataInS3,
				ParametersCountWeight:   &parametersCountWeight,
				PathVariativityWeight:   &pathVariativityWeight,
				PiiWeight:               &piiWeight,
				RequestContentWeight:    &requestContentWeight,
				OpenVulnsWeight:         &openVulnsWeight,
				SerializedDataWeight:    &serializedDataWeight,
				RiskScoreAlgo:           &riskScoreAlgo,
				PiiFallback:             &piiFallback,
			},
		},
	}

	if assert.NoError(t, err) {
		assert.Equal(t, expected, *actual)
	}
}
