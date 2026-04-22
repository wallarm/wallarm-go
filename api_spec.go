package wallarm

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type (
	// APISpec contains operations available on APISpec resource
	APISpec interface {
		APISpecCreate(body *APISpecCreate) (APISpecCreateResp, error)
		APISpecReadByID(clientID, specID int) (APISpecBody, error)
		APISpecUpdate(clientID, specID int, body *APISpecUpdate) (APISpecCreateResp, error)
		APISpecList(clientID int, page, perPage int) (APISpecListResp, error)
		APISpecDelete(clientID, specID int) error
		APISpecPolicyPut(clientID, specID int, body *APISpecPolicy) (APISpecPolicyResp, error)
	}

	APISpecCreate struct {
		Title             string              `json:"title"`
		Description       string              `json:"description"`
		FileRemoteURL     string              `json:"file_remote_url"`
		RegularFileUpdate bool                `json:"regular_file_update"`
		APIDetection      bool                `json:"api_detection"`
		ClientID          int                 `json:"-"`
		Instances         []interface{}       `json:"instances"`
		Domains           []interface{}       `json:"domains"`
		AuthHeaders       []APISpecAuthHeader `json:"auth_headers,omitempty"`
	}

	APISpecUpdate struct {
		Title             *string             `json:"title,omitempty"`
		Description       *string             `json:"description,omitempty"`
		FileRemoteURL     *string             `json:"file_remote_url,omitempty"`
		RegularFileUpdate *bool               `json:"regular_file_update,omitempty"`
		APIDetection      *bool               `json:"api_detection,omitempty"`
		Instances         []interface{}       `json:"instances,omitempty"`
		Domains           []interface{}       `json:"domains,omitempty"`
		AuthHeaders       []APISpecAuthHeader `json:"auth_headers,omitempty"`
	}

	APISpecCreateResp struct {
		Status int          `json:"status"`
		Body   *APISpecBody `json:"body"`
	}

	APISpecPolicyResp struct {
		Status int            `json:"status"`
		Body   *APISpecPolicy `json:"body"`
	}

	APISpecBody struct {
		ID                   int                 `json:"id"`
		ClientID             int                 `json:"client_id"`
		Title                string              `json:"title"`
		Description          string              `json:"description"`
		Status               string              `json:"status"`
		Instances            []interface{}       `json:"instances"`
		Domains              []interface{}       `json:"domains"`
		RegularFileUpdate    bool                `json:"regular_file_update"`
		APIDetection         bool                `json:"api_detection"`
		SpecVersion          string              `json:"spec_version"`
		Version              int                 `json:"version"`
		EndpointsCount       int                 `json:"endpoints_count"`
		ShadowEndpointsCount int                 `json:"shadow_endpoints_count"`
		OrphanEndpointsCount int                 `json:"orphan_endpoints_count"`
		ZombieEndpointsCount int                 `json:"zombie_endpoints_count"`
		OpenAPIVersion       string              `json:"openapi_version"`
		LastSyncedAt         string              `json:"last_synced_at"`
		LastComparedAt       string              `json:"last_compared_at"`
		UpdatedAt            string              `json:"updated_at"`
		CreatedAt            string              `json:"created_at"`
		NodeSyncVersion      int                 `json:"node_sync_version"`
		FileRemoteURL        string              `json:"file_remote_url"`
		File                 *APISpecFile        `json:"file,omitempty"`
		Policy               *APISpecPolicy      `json:"policy,omitempty"`
		AuthHeaders          []APISpecAuthHeader `json:"auth_headers,omitempty"`
		FileChangedAt        string              `json:"file_changed_at,omitempty"`
		Format               int                 `json:"format,omitempty"`
	}

	APISpecPolicyCondition struct {
		Type  string        `json:"type"`
		Value interface{}   `json:"value"`
		Point []interface{} `json:"point"`
	}

	APISpecPolicy struct {
		Enabled                   bool                     `json:"enabled"`
		Conditions                []APISpecPolicyCondition `json:"conditions,omitempty"`
		UndefinedEndpointMode     string                   `json:"undefined_endpoint_mode,omitempty"`
		UndefinedParameterMode    string                   `json:"undefined_parameter_mode,omitempty"`
		MissingParameterMode      string                   `json:"missing_parameter_mode,omitempty"`
		InvalidParameterValueMode string                   `json:"invalid_parameter_value_mode,omitempty"`
		MissingAuthMode           string                   `json:"missing_auth_mode,omitempty"`
		InvalidRequestMode        string                   `json:"invalid_request_mode,omitempty"`
		TimeoutMode               string                   `json:"timeout_mode,omitempty"`
		MaxRequestSizeMode        string                   `json:"max_request_size_mode,omitempty"`
		Timeout                   int                      `json:"timeout,omitempty"`
		MaxRequestSize            int                      `json:"max_request_size,omitempty"`
	}

	APISpecAuthHeader struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	APISpecFile struct {
		Name      string `json:"name"`
		SignedURL string `json:"signed_url"`
		Checksum  string `json:"checksum"`
		MimeType  string `json:"mime_type"`
		Version   int    `json:"version"`
	}

	APISpecListResp struct {
		Items       []APISpecBody `json:"items"`
		CurrentPage int           `json:"current_page"`
		PerPage     int           `json:"per_page"`
		TotalPages  int           `json:"total_pages"`
		TotalCount  int           `json:"total_count"`
	}
)

var ErrNotFound = errors.New("APISpec not found")

func (api *api) APISpecReadByID(clientID, specID int) (APISpecBody, error) {
	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs/%d", clientID, specID)
	respBody, err := api.makeRequest("GET", uri, "api_spec", nil, nil)
	if err != nil {
		return APISpecBody{}, fmt.Errorf("APISpecReadByID: failed to make request - %w", err)
	}
	var resp APISpecCreateResp
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return APISpecBody{}, fmt.Errorf("APISpecReadByID: failed to parse response - %w", err)
	}
	if resp.Body == nil {
		return APISpecBody{}, fmt.Errorf("APISpecReadByID: %w", ErrNotFound)
	}
	return *resp.Body, nil
}

func (api *api) APISpecCreate(apiSpecBody *APISpecCreate) (APISpecCreateResp, error) {

	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs", apiSpecBody.ClientID)
	respBody, err := api.makeRequest("POST", uri, "api_spec", apiSpecBody, nil)
	var a APISpecCreateResp
	if err != nil {
		return a, fmt.Errorf("APISpecCreate: failed to make request - %w", err)
	}

	if err = json.Unmarshal(respBody, &a); err != nil {
		return a, fmt.Errorf("APISpecCreate: failed to parse response - %w", err)
	}
	return a, nil
}

func (api *api) APISpecUpdate(clientID, specID int, body *APISpecUpdate) (APISpecCreateResp, error) {
	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs/%d", clientID, specID)
	respBody, err := api.makeRequest("PUT", uri, "api_spec", body, nil)
	var resp APISpecCreateResp
	if err != nil {
		return resp, fmt.Errorf("APISpecUpdate: failed to make request - %w", err)
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return resp, fmt.Errorf("APISpecUpdate: failed to parse response - %w", err)
	}
	return resp, nil
}

func (api *api) APISpecList(clientID, page, perPage int) (APISpecListResp, error) {
	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs?page=%d&per_page=%d", clientID, page, perPage)
	respBody, err := api.makeRequest("GET", uri, "api_spec", nil, nil)
	if err != nil {
		return APISpecListResp{}, fmt.Errorf("APISpecList: failed to make request - %w", err)
	}
	var resp APISpecListResp
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return APISpecListResp{}, fmt.Errorf("APISpecList: failed to parse response - %w", err)
	}
	return resp, nil
}

func (api *api) APISpecDelete(clientID, apiSpecID int) error {
	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs/%d", clientID, apiSpecID)

	_, err := api.makeRequest("DELETE", uri, "api_spec", nil, nil)
	if err != nil {
		return fmt.Errorf("APISpecDelete: failed to make request - %w", err)
	}
	return nil
}

func (api *api) APISpecPolicyPut(clientID, specID int, body *APISpecPolicy) (APISpecPolicyResp, error) {
	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs/%d/policy", clientID, specID)
	respBody, err := api.makeRequest("PUT", uri, "api_spec_policy", body, nil)
	var resp APISpecPolicyResp
	if err != nil {
		return resp, fmt.Errorf("APISpecPolicyPut: failed to make request - %w", err)
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return resp, fmt.Errorf("APISpecPolicyPut: failed to parse response - %w", err)
	}
	return resp, nil
}
