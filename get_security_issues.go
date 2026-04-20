package wallarm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var securityIssuesJSONHeaders = map[string]string{"Content-Type": "application/json"}

type (
	// SecurityIssues contains read operations available on Security Issues resources.
	SecurityIssues interface {
		GetSecurityIssuesRead(body *GetSecurityIssuesRead) ([]*GetSecurityIssuesResp, error)
		GetSecurityIssuesCount(body *GetSecurityIssuesCount) (*GetSecurityIssuesCountResp, error)
		GetSecurityIssueGroups(body *GetSecurityIssueGroups) ([]*GetSecurityIssueGroupResp, error)
		GetSecurityIssueGroupsCount(body *GetSecurityIssueGroupsCount) (*GetSecurityIssueGroupsCountResp, error)
		GetSecurityIssue(body *GetSecurityIssue) (*GetSecurityIssueResp, error)
	}

	GetSecurityIssuesRead struct {
		ClientID  int                       `json:"client_id"`
		Token     string                    `json:"token,omitempty"`
		Offset    int                       `json:"offset,omitempty"`
		Limit     int                       `json:"limit,omitempty"`
		Unlimited bool                      `json:"unlimited,omitempty"`
		Filter    *GetSecurityIssuesFilter  `json:"filter,omitempty"`
		OrderBy   *GetSecurityIssuesOrderBy `json:"order_by,omitempty"`
	}

	GetSecurityIssuesCount struct {
		ClientID int                      `json:"client_id"`
		Token    string                   `json:"token,omitempty"`
		Filter   *GetSecurityIssuesFilter `json:"filter,omitempty"`
	}

	GetSecurityIssueGroups struct {
		ClientID  int                       `json:"client_id"`
		Token     string                    `json:"token,omitempty"`
		Offset    int                       `json:"offset,omitempty"`
		Limit     int                       `json:"limit,omitempty"`
		Unlimited bool                      `json:"unlimited,omitempty"`
		Filter    *GetSecurityIssuesFilter  `json:"filter,omitempty"`
		OrderBy   *GetSecurityIssuesOrderBy `json:"order_by,omitempty"`
	}

	GetSecurityIssueGroupsCount struct {
		ClientID int                      `json:"client_id"`
		Token    string                   `json:"token,omitempty"`
		Filter   *GetSecurityIssuesFilter `json:"filter,omitempty"`
	}

	GetSecurityIssue struct {
		ID       int
		ClientID int
		Token    string
	}

	GetSecurityIssuesOrderBy struct {
		Name      string `json:"name"`
		Direction string `json:"direction"`
	}

	GetSecurityIssuesFilter struct {
		ClientID           int      `json:"client_id,omitempty"`
		NotClientID        int      `json:"!client_id,omitempty"`
		Severity           []string `json:"severity,omitempty"`
		NotSeverity        []string `json:"!severity,omitempty"`
		Host               string   `json:"host,omitempty"`
		NotHost            string   `json:"!host,omitempty"`
		State              []string `json:"state,omitempty"`
		NotState           []string `json:"!state,omitempty"`
		CreatedSince       int      `json:"created_since,omitempty"`
		DiscoveredSince    int      `json:"discovered_since,omitempty"`
		DiscoveredBy       []string `json:"discovered_by,omitempty"`
		NotDiscoveredBy    []string `json:"!discovered_by,omitempty"`
		ID                 int      `json:"id,omitempty"`
		NotID              int      `json:"!id,omitempty"`
		DomainID           int      `json:"domain_id,omitempty"`
		NotDomainID        int      `json:"!domain_id,omitempty"`
		SubdomainID        int      `json:"subdomain_id,omitempty"`
		NotSubdomainID     int      `json:"!subdomain_id,omitempty"`
		IssueType          string   `json:"issue_type,omitempty"`
		NotIssueType       string   `json:"!issue_type,omitempty"`
		Owasp              string   `json:"owasp,omitempty"`
		NotOwasp           string   `json:"!owasp,omitempty"`
		SourceTemplate     string   `json:"source_template,omitempty"`
		NotSourceTemplate  string   `json:"!source_template,omitempty"`
		GroupID            string   `json:"group_id,omitempty"`
		NotGroupID         string   `json:"!group_id,omitempty"`
		SearchQuery        string   `json:"search_query,omitempty"`
		TestRunPublicUuids string   `json:"test_run_public_uuids,omitempty"`
		Verified           bool     `json:"verified,omitempty"`
		Incident           *bool    `json:"incident,omitempty"`
	}

	SecurityIssueType struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	SecurityIssueOWASP struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Link     string `json:"link"`
	}

	SecurityIssueTag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	SecurityIssueVPatch struct {
		RuleID *int `json:"rule_id"`
	}

	SecurityIssueMitigations struct {
		Vpatch *SecurityIssueVPatch `json:"vpatch,omitempty"`
	}

	GetSecurityIssuesResp struct {
		ID                      int                      `json:"id"`
		ClientID                int                      `json:"client_id"`
		Severity                string                   `json:"severity"`
		State                   string                   `json:"state"`
		Volume                  int                      `json:"volume"`
		Name                    string                   `json:"name"`
		CreatedAt               int                      `json:"created_at"`
		DiscoveredAt            int                      `json:"discovered_at"`
		DiscoveredBy            string                   `json:"discovered_by"`
		DiscoveredByDisplayName string                   `json:"discovered_by_display_name"`
		URL                     string                   `json:"url"`
		Host                    string                   `json:"host"`
		Path                    string                   `json:"path"`
		ParameterDisplayName    string                   `json:"parameter_display_name"`
		ParameterPosition       string                   `json:"parameter_position"`
		ParameterName           string                   `json:"parameter_name"`
		HTTPMethod              string                   `json:"http_method"`
		AASMTemplate            string                   `json:"aasm_template"`
		Mitigations             SecurityIssueMitigations `json:"mitigations"`
		IssueType               SecurityIssueType        `json:"issue_type"`
		Owasp                   []SecurityIssueOWASP     `json:"owasp"`
		Tags                    []SecurityIssueTag       `json:"tags"`
		Incident                bool                     `json:"incident"`
		FalsePositiveRuleID     *int                     `json:"false_positive_rule_id"`
		Verified                bool                     `json:"verified"`
	}

	GetSecurityIssuesCountResp struct {
		Count int `json:"count"`
	}

	SecurityIssueStates struct {
		Open        int `json:"open"`
		Closed      int `json:"closed"`
		MarkedFalse int `json:"marked_false"`
		Hidden      int `json:"hidden,omitempty"`
	}

	GetSecurityIssueGroupResp struct {
		GroupID                 string               `json:"group_id"`
		Title                   string               `json:"title"`
		Severity                string               `json:"severity"`
		DiscoveredBy            string               `json:"discovered_by"`
		DiscoveredByDisplayName string               `json:"discovered_by_display_name"`
		IssueType               SecurityIssueType    `json:"issue_type"`
		SecurityIssuesCount     int                  `json:"security_issues_count"`
		States                  SecurityIssueStates  `json:"states"`
		Owasp                   []SecurityIssueOWASP `json:"owasp"`
		HostsCount              int                  `json:"hosts_count"`
		FirstHost               string               `json:"first_host"`
		ClientID                int                  `json:"client_id"`
		Tags                    []SecurityIssueTag   `json:"tags"`
	}

	GetSecurityIssueGroupsCountResp struct {
		Count int `json:"count"`
	}

	GetSecurityIssueResp struct {
		GetSecurityIssuesResp
		Credentials           []string                 `json:"credentials"`
		AdditionalInfo        string                   `json:"additional_info"`
		ExploitationExamples  []interface{}            `json:"exploitation_examples"`
		OOBInteraction        []interface{}            `json:"oob_interaction"`
		PassiveDetectIncident bool                     `json:"passive_detect_incident"`
		Description           string                   `json:"description"`
		References            []string                 `json:"references"`
		Mitigation            string                   `json:"mitigation"`
		CWE                   []interface{}            `json:"cwe"`
		RiskInfo              map[string]interface{}   `json:"risk_info"`
		Source                string                   `json:"source"`
		IssueSubtype          map[string]interface{}   `json:"issue_subtype"`
		PassiveDetectURLs     []string                 `json:"passive_detect_urls"`
		Manual                bool                     `json:"manual"`
		LeaksInfo             []interface{}            `json:"leaks_info"`
		URLActive             bool                     `json:"url_active"`
		StatusHistory         []map[string]interface{} `json:"status_history"`
		Recheckable           bool                     `json:"recheckable"`
		LastRecheck           map[string]interface{}   `json:"last_recheck"`
	}
)

func (api *api) GetSecurityIssuesRead(body *GetSecurityIssuesRead) ([]*GetSecurityIssuesResp, error) {
	if body == nil {
		return nil, fmt.Errorf("security issues read request is required")
	}
	request := *body
	request.Token = api.attackSurfaceToken(request.Token)

	respBody, err := api.makeRequest(http.MethodPost, "/v1/security_issues", "security_issues", &request, securityIssuesJSONHeaders)
	if err != nil {
		return nil, err
	}

	var resp []*GetSecurityIssuesResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (api *api) GetSecurityIssuesCount(body *GetSecurityIssuesCount) (*GetSecurityIssuesCountResp, error) {
	if body == nil {
		return nil, fmt.Errorf("security issues count request is required")
	}
	request := *body
	request.Token = api.attackSurfaceToken(request.Token)

	respBody, err := api.makeRequest(http.MethodPost, "/v1/security_issues/count", "security_issues", &request, securityIssuesJSONHeaders)
	if err != nil {
		return nil, err
	}

	var resp GetSecurityIssuesCountResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) GetSecurityIssueGroups(body *GetSecurityIssueGroups) ([]*GetSecurityIssueGroupResp, error) {
	if body == nil {
		return nil, fmt.Errorf("security issue groups request is required")
	}
	request := *body
	request.Token = api.attackSurfaceToken(request.Token)

	respBody, err := api.makeRequest(http.MethodPost, "/v1/security_issues/groups", "security_issues", &request, securityIssuesJSONHeaders)
	if err != nil {
		return nil, err
	}

	var resp []*GetSecurityIssueGroupResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (api *api) GetSecurityIssueGroupsCount(body *GetSecurityIssueGroupsCount) (*GetSecurityIssueGroupsCountResp, error) {
	if body == nil {
		return nil, fmt.Errorf("security issue groups count request is required")
	}
	request := *body
	request.Token = api.attackSurfaceToken(request.Token)

	respBody, err := api.makeRequest(http.MethodPost, "/v1/security_issues/groups_count", "security_issues", &request, securityIssuesJSONHeaders)
	if err != nil {
		return nil, err
	}

	var resp GetSecurityIssueGroupsCountResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) GetSecurityIssue(body *GetSecurityIssue) (*GetSecurityIssueResp, error) {
	if body == nil {
		return nil, fmt.Errorf("security issue request is required")
	}
	query := url.Values{}
	query.Set("client_id", fmt.Sprintf("%d", body.ClientID))
	query.Set("token", api.attackSurfaceToken(body.Token))

	respBody, err := api.makeRequest(http.MethodGet, fmt.Sprintf("/v1/security_issues/%d", body.ID), "security_issues", query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var resp GetSecurityIssueResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) attackSurfaceToken(explicit string) string {
	if strings.TrimSpace(explicit) != "" {
		return strings.TrimSpace(explicit)
	}
	if token := api.headers.Get("X-WallarmAPI-Token"); strings.TrimSpace(token) != "" {
		return strings.TrimSpace(token)
	}
	if token := api.headers.Get("X-WallarmApi-Token"); strings.TrimSpace(token) != "" {
		return strings.TrimSpace(token)
	}
	for name, values := range api.headers {
		if !(strings.EqualFold(name, "X-WallarmAPI-Token") || strings.EqualFold(name, "X-WallarmApi-Token")) {
			continue
		}
		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}
