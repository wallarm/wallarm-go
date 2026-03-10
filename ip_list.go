package wallarm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type (
	IPListType string

	IPList interface {
		IPListRead(listType IPListType, clientID int) ([]IPRule, error)
		IPListCreate(clientID int, params AccessRuleCreateRequest) error
		IPListDelete(clientID int, rules []AccessRuleDeleteEntry) error
	}

	AccessRuleCreateRequest struct {
		List           IPListType        `json:"list"`
		Force          bool              `json:"force"`
		Reason         string            `json:"reason"`
		ApplicationIDs []int             `json:"application_ids"`
		ExpiredAt      int               `json:"expired_at"`
		Rules          []AccessRuleEntry `json:"rules"`
	}

	AccessRuleEntry struct {
		RulesType string   `json:"rules_type"`
		Values    []string `json:"values"`
	}

	AccessRuleDeleteRequest struct {
		Rules []AccessRuleDeleteEntry `json:"rules"`
	}

	AccessRuleDeleteEntry struct {
		RuleType string `json:"rule_type"`
		IDs      []int  `json:"ids"`
	}

	IPRule struct {
		ID             int      `json:"id"`
		ClientID       int      `json:"client_id"`
		RuleType       string   `json:"rule_type"`
		List           string   `json:"list"`
		CreatedAt      int      `json:"created_at"`
		ExpiredAt      int      `json:"expired_at"`
		ApplicationIDs []int    `json:"application_ids"`
		Reason         string   `json:"reason"`
		AuthorUserID   int      `json:"author_user_id"`
		Values         []string `json:"values"`
		Status         string   `json:"status"`
	}
)

const (
	DenylistType  IPListType = "block"
	AllowlistType IPListType = "allow"
	GraylistType  IPListType = "gray"
)

func (api *api) IPListRead(listType IPListType, clientID int) ([]IPRule, error) {
	uri := fmt.Sprintf("/v1/blocklist/clients/%d/groups", clientID)

	q := url.Values{}
	q.Add("filter[rule_type][]", "subnet")
	q.Add("filter[rule_type][]", "proxy_type")
	q.Add("filter[rule_type][]", "datacenter")
	q.Add("filter[rule_type][]", "location")
	q.Set("filter[list]", string(listType))
	q.Set("limit", "50")

	var response struct {
		Body struct {
			Objects []IPRule `json:"objects"`
		} `json:"body"`
	}

	var result []IPRule
	offset := 0

	for {
		q.Set("offset", strconv.Itoa(offset))

		respBody, err := api.makeRequest("GET", uri, "", q.Encode(), nil)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(respBody, &response); err != nil {
			return nil, err
		}

		result = append(result, response.Body.Objects...)

		if len(response.Body.Objects) < 50 {
			break
		}

		offset += 50
	}

	return result, nil
}

func (api *api) IPListCreate(clientID int, params AccessRuleCreateRequest) error {
	uri := fmt.Sprintf("/v1/blocklist/clients/%d/access_rules", clientID)

	_, err := api.makeRequest("POST", uri, "", &params, nil)

	return err
}

func (api *api) IPListDelete(clientID int, rules []AccessRuleDeleteEntry) error {
	uri := fmt.Sprintf("/v1/blocklist/clients/%d/groups", clientID)

	reqBody := AccessRuleDeleteRequest{
		Rules: rules,
	}

	_, err := api.makeRequest("DELETE", uri, "ip_rules", &reqBody, nil)

	return err
}
