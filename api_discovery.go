package wallarm

import (
	"encoding/json"
	"fmt"
)

type (
	// APIDiscovery contains operations available on the API Discovery
	// configuration resource. The config is a singleton per client_id —
	// no Create/Delete; only Read and Update.
	APIDiscovery interface {
		APIDiscoveryConfigRead(clientID int) (*APIDiscoveryConfig, error)
	}

	// APIDiscoveryConfig mirrors the API Discovery config body returned by
	// GET /v1/clients/{client_id}/apid/config and accepted by POST on the
	// same path.
	APIDiscoveryConfig struct {
		ClientID                    int                             `json:"clientid"`
		Enabled                     bool                            `json:"enabled"`
		Protocols                   APIDiscoveryProtocols           `json:"protocols"`
		ApplyExtendedFilter         bool                            `json:"apply_extended_filter"`
		TypeDetectionThreshold      float64                         `json:"type_detection_threshold"`
		PIIDetectionThreshold       float64                         `json:"pii_detection_threshold"`
		CallPointsStorageLimit      int                             `json:"call_points_storage_limit"`
		SensitiveSamples            APIDiscoverySensitiveSamples    `json:"sensitive_samples"`
		DisabledApps                []int                           `json:"disabled_apps"`
		EndpointStability           APIDiscoveryEndpointStability   `json:"endpoint_stability"`
		GroupSOAP                   bool                            `json:"group_soap"`
		ServerVariability           APIDiscoveryServerVariability   `json:"server_variability"`
		AllowedContentTypesPatterns []string                        `json:"allowed_content_types_patterns"`
		ExtensionsWhitelist         APIDiscoveryExtensionsWhitelist `json:"extensions_whitelist"`
	}

	APIDiscoveryProtocols struct {
		REST    bool `json:"rest"`
		GraphQL bool `json:"graphql"`
		SOAP    bool `json:"soap"`
		GRPC    bool `json:"grpc"`
		MCP     bool `json:"mcp"`
	}

	APIDiscoverySensitiveSamples struct {
		Enabled     bool `json:"enabled"`
		MinMasked   int  `json:"min_masked"`
		MaxMasked   int  `json:"max_masked"`
		MaskSymbols bool `json:"mask_symbols"`
	}

	APIDiscoveryEndpointStability struct {
		MinCount int `json:"min_count"`
		MinTime  int `json:"min_time"`
	}

	APIDiscoveryServerVariability struct {
		Enabled                 bool                                       `json:"enabled"`
		ByDateEnabled           bool                                       `json:"by_date_enabled"`
		ByLocalCodeEnabled      bool                                       `json:"by_local_code_enabled"`
		ByEmailEnabled          bool                                       `json:"by_email_enabled"`
		ByAlphanumericIDEnabled bool                                       `json:"by_alphanumeric_id_enabled"`
		ByCustomPaths           APIDiscoveryServerVariabilityByCustomPaths `json:"by_custom_paths"`
	}

	APIDiscoveryServerVariabilityByCustomPaths struct {
		Enabled bool     `json:"enabled"`
		Paths   []string `json:"paths"`
	}

	APIDiscoveryExtensionsWhitelist struct {
		Enabled    bool     `json:"enabled"`
		Extensions []string `json:"extensions"`
	}

	APIDiscoveryConfigResp struct {
		Status int                 `json:"status"`
		Body   *APIDiscoveryConfig `json:"body"`
	}
)

// APIDiscoveryConfigRead fetches the API Discovery configuration for the
// given client.
func (api *api) APIDiscoveryConfigRead(clientID int) (*APIDiscoveryConfig, error) {
	uri := fmt.Sprintf("/v1/clients/%d/apid/config", clientID)
	respBody, err := api.makeRequest("GET", uri, "api_discovery", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("APIDiscoveryConfigRead: %w", err)
	}
	var resp APIDiscoveryConfigResp
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("APIDiscoveryConfigRead: failed to parse response - %w", err)
	}
	return resp.Body, nil
}
