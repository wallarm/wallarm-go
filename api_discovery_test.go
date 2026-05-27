package wallarm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIDiscoveryConfigRead(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{
			"status": 200,
			"body": {
				"enabled": true,
				"protocols": {"rest": true, "graphql": true, "soap": true, "grpc": true, "mcp": true},
				"apply_extended_filter": true,
				"type_detection_threshold": 0.5,
				"pii_detection_threshold": 0.1,
				"call_points_storage_limit": 50000,
				"sensitive_samples": {"enabled": false, "min_masked": 20, "max_masked": 80, "mask_symbols": false},
				"disabled_apps": [],
				"endpoint_stability": {"min_count": 2, "min_time": 300},
				"group_soap": false,
				"server_variability": {
					"enabled": false,
					"by_date_enabled": false,
					"by_local_code_enabled": false,
					"by_email_enabled": false,
					"by_alphanumeric_id_enabled": false,
					"by_custom_paths": {"enabled": false, "paths": []}
				},
				"allowed_content_types_patterns": ["text/xml", "application/%json"],
				"extensions_whitelist": {"enabled": true, "extensions": ["", "do", "action"]}
			}
		}`)
	}

	mux.HandleFunc("/v1/clients/22510/apid/config", handler)
	cfg, err := client.APIDiscoveryConfigRead(22510)
	assert.NoError(t, err)
	if !assert.NotNil(t, cfg) {
		return
	}

	assert.True(t, cfg.Enabled)
	assert.True(t, cfg.Protocols.REST)
	assert.True(t, cfg.Protocols.GraphQL)
	assert.True(t, cfg.Protocols.SOAP)
	assert.True(t, cfg.Protocols.GRPC)
	assert.True(t, cfg.Protocols.MCP)
	assert.True(t, cfg.ApplyExtendedFilter)
	assert.InEpsilon(t, 0.5, cfg.TypeDetectionThreshold, 0.001)
	assert.InEpsilon(t, 0.1, cfg.PIIDetectionThreshold, 0.001)
	assert.Equal(t, 50000, cfg.CallPointsStorageLimit)
	assert.False(t, cfg.SensitiveSamples.Enabled)
	assert.Equal(t, 20, cfg.SensitiveSamples.MinMasked)
	assert.Equal(t, 80, cfg.SensitiveSamples.MaxMasked)
	assert.False(t, cfg.SensitiveSamples.MaskSymbols)
	assert.Empty(t, cfg.DisabledApps)
	assert.Equal(t, 2, cfg.EndpointStability.MinCount)
	assert.Equal(t, 300, cfg.EndpointStability.MinTime)
	assert.False(t, cfg.GroupSOAP)
	assert.False(t, cfg.ServerVariability.Enabled)
	assert.False(t, cfg.ServerVariability.ByDateEnabled)
	assert.False(t, cfg.ServerVariability.ByCustomPaths.Enabled)
	assert.Empty(t, cfg.ServerVariability.ByCustomPaths.Paths)
	assert.Equal(t, []string{"text/xml", "application/%json"}, cfg.AllowedContentTypesPatterns)
	assert.True(t, cfg.ExtensionsWhitelist.Enabled)
	assert.Equal(t, []string{"", "do", "action"}, cfg.ExtensionsWhitelist.Extensions)
}
