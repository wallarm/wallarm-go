package wallarm

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type (
	// ApiSpec contains operations available on ApiSpec resource
	APISpec interface {
		APISpecCreate(apiSpecBody *APISpecCreate) (APISpecCreateResp, error)
		APISpecDelete(clientID int, apiSpecID int) error
		APISpecRead(clientID int, id int) (APISpecBody, error)
	}

	APISpecCreate struct {
		Title             string        `json:"title"`
		Description       string        `json:"description"`
		FileRemoteURL     string        `json:"file_remote_url"`
		RegularFileUpdate bool          `json:"regular_file_update"`
		APIDetection      bool          `json:"api_detection"`
		ClientID          int           `json:"-"`
		Instances         []interface{} `json:"instances"`
		Domains           []interface{} `json:"domains"`
	}

	APISpecCreateResp struct {
		Status int          `json:"status"`
		Body   *APISpecBody `json:"body"`
	}

	APISpecBody struct {
		ID                   int           `json:"id"`
		ClientID             int           `json:"client_id"`
		Title                string        `json:"title"`
		Description          string        `json:"description"`
		Status               string        `json:"status"`
		Instances            []interface{} `json:"instances"`
		Domains              []interface{} `json:"domains"`
		RegularFileUpdate    bool          `json:"regular_file_update"`
		APIDetection         bool          `json:"api_detection"`
		SpecVersion          string        `json:"spec_version"`
		Version              int           `json:"version"`
		EndpointsCount       int           `json:"endpoints_count"`
		ShadowEndpointsCount int           `json:"shadow_endpoints_count"`
		OrphanEndpointsCount int           `json:"orphan_endpoints_count"`
		ZombieEndpointsCount int           `json:"zombie_endpoints_count"`
		OpenAPIVersion       string        `json:"openapi_version"`
		LastSyncedAt         string        `json:"last_synced_at"`
		LastComparedAt       string        `json:"last_compared_at"`
		UpdatedAt            string        `json:"updated_at"`
		CreatedAt            string        `json:"created_at"`
		NodeSyncVersion      int           `json:"node_sync_version"`
		FileRemoteURL        string        `json:"file_remote_url"`
		File                 struct {
			Name      string `json:"name"`
			SignedURL string `json:"signed_url"`
			Checksum  string `json:"checksum"`
			MimeType  string `json:"mime_type"`
			Version   int    `json:"version"`
		} `json:"file"`
	}

	APISpecRead struct {
		Items       []APISpecBody `json:"items"`
		CurrentPage int           `json:"current_page"`
		PerPage     int           `json:"per_page"`
		TotalPages  int           `json:"total_pages"`
		TotalCount  int           `json:"total_count"`
	}
)

var ErrNotFound = errors.New("ApiSpec not found")

func (api *api) APISpecRead(clientID int, id int) (APISpecBody, error) {

	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs", clientID)
	var apiSpecBody APISpecBody
	respBody, err := api.makeRequest("GET", uri, "api_spec", nil, nil)
	if err != nil {
		return apiSpecBody, fmt.Errorf("ApiSpecRead: failed to make request - %w", err)
	}
	var readResult APISpecRead
	if err = json.Unmarshal(respBody, &readResult); err != nil {
		return apiSpecBody, fmt.Errorf("ApiSpecRead: failed to parse response - %w", err)
	}
	for _, obj := range readResult.Items {
		if obj.ID == id {
			return obj, nil
		}
	}

	return apiSpecBody, fmt.Errorf("ApiSpecRead: %w - body: %s", ErrNotFound, string(respBody))
}

func (api *api) APISpecCreate(apiSpecBody *APISpecCreate) (APISpecCreateResp, error) {

	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs", apiSpecBody.ClientID)
	respBody, err := api.makeRequest("POST", uri, "api_spec", apiSpecBody, nil)
	var a APISpecCreateResp
	if err != nil {
		return a, fmt.Errorf("ApiSpecCreate: failed to make request - %w", err)
	}

	if err = json.Unmarshal(respBody, &a); err != nil {
		return a, fmt.Errorf("ApiSpecCreate: failed to parse response - %w", err)
	}
	return a, nil
}

func (api *api) APISpecDelete(clientID int, apiSpecID int) error {
	uri := fmt.Sprintf("/v4/clients/%d/rules/api-specs/%d", clientID, apiSpecID)

	_, err := api.makeRequest("DELETE", uri, "api_spec", nil, nil)
	if err != nil {
		return fmt.Errorf("ApiSpecDelete: failed to make request - %w", err)
	}
	return nil
}
