package wallarm

import (
	"encoding/json"
	"net/http"
)

type (
	// SecurityIssues contains operations available on SecurityIssues resource
	SecurityIssues interface {
		GetSecurityIssuesRead(getSecurityIssuesBody *GetSecurityIssuesRead) ([]*GetSecurityIssuesResp, error)
	}

	// GetSecurityIssuesRead is a root object for requesting security issues.
	// Limit is a number between 0 - 1000
	// Offset is a number
	GetSecurityIssuesRead struct {
		ClientID  int                       `json:"client_id"`
		Offset    int                       `json:"offset"`
		Limit     int                       `json:"limit"`
		Unlimited bool                      `json:"unlimited"`
		Filter    *GetSecurityIssuesFilter  `json:"filter"`
		OrderBy   *GetSecurityIssuesOrderBy `json:"order_by"`
	}

	GetSecurityIssuesOrderBy struct {
		Name      string `json:"name"`
		Direction string `json:"direction"`
	}

	GetSecurityIssuesFilter struct {
		ClientId           int      `json:"client_id,omitempty"`
		NotClientId        int      `json:"!client_id,omitempty"`
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
		Id                 int      `json:"id,omitempty"`
		NotId              int      `json:"!id,omitempty"`
		DomainId           int      `json:"domain_id,omitempty"`
		NotDomainId        int      `json:"!domain_id,omitempty"`
		SubdomainId        int      `json:"subdomain_id,omitempty"`
		NotSubdomainId     int      `json:"!subdomain_id,omitempty"`
		IssueType          string   `json:"issue_type,omitempty"`
		NotIssueType       string   `json:"!issue_type,omitempty"`
		Owasp              string   `json:"owasp,omitempty"`
		NotOwasp           string   `json:"!owasp,omitempty"`
		SourceTemplate     string   `json:"source_template,omitempty"`
		NotSourceTemplate  string   `json:"!source_template,omitempty"`
		GroupId            string   `json:"group_id,omitempty"`
		NotGroupId         string   `json:"!group_id,omitempty"`
		SearchQuery        string   `json:"search_query,omitempty"`
		TestRunPublicUuids string   `json:"test_run_public_uuids,omitempty"`
		Verified           bool     `json:"verified,omitempty"`
	}

	GetSecurityIssuesResp struct {
		Id                      int    `json:"id"`
		ClientId                int    `json:"client_id"`
		Severity                string `json:"severity"`
		State                   string `json:"state"`
		Volume                  int    `json:"volume"`
		Name                    string `json:"name"`
		CreatedAt               int    `json:"created_at"`
		DiscoveredAt            int    `json:"discovered_at"`
		DiscoveredBy            string `json:"discovered_by"`
		DiscoveredByDisplayName string `json:"discovered_by_display_name"`
		Url                     string `json:"url"`
		Host                    string `json:"host"`
		Path                    string `json:"path"`
		ParameterDisplayName    string `json:"parameter_display_name"`
		ParameterPosition       string `json:"parameter_position"`
		ParameterName           string `json:"parameter_name"`
		HttpMethod              string `json:"http_method"`
		AasmTemplate            string `json:"aasm_template"`
		Mitigations             struct {
			Vpatch struct {
				RuleId int `json:"rule_id"`
			} `json:"vpatch"`
		} `json:"mitigations"`
		IssueType struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"issue_type"`
		Owasp []struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			Link     string `json:"link"`
		} `json:"owasp"`
		Tags []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"tags"`
		Verified bool `json:"verified"`
	}
)

func (api *api) GetSecurityIssuesRead(getSecurityIssuesBody *GetSecurityIssuesRead) ([]*GetSecurityIssuesResp, error) {
	uri := "/v1/security_issues"
	respBody, err := api.makeRequest(http.MethodGet, uri, "security_issues", getSecurityIssuesBody,
		map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var v []*GetSecurityIssuesResp
	if err = json.Unmarshal(respBody, &v); err != nil {
		return nil, err
	}
	return v, nil
}
