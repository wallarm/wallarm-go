package wallarm

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type (
	// Blacklist contains operations available on Blacklist resource
	Blacklist interface {
		BlacklistRead(clientID int) ([]IPRule, error)
		BlacklistCreate(clientID int, params IPRuleCreationParams) error
		BlacklistDelete(clientID int, ids []int) error
	}

	IPRuleCreationParams struct {
		ExpiredAt int    `json:"expired_at"`
		List      string `json:"list"`
		Pools     []int  `json:"pools"`
		Reason    string `json:"reason"`
		RuleType  string `json:"rule_type"`
		Subnet    string `json:"subnet"`
	}

	IPRule struct {
		ID int `json:"id"`
		ClientID int `json:"clientid"`
		RuleType string `json:"rule_type"`
		List string `json:"list"`
		Author string `json:"author"`
		CreatedAt int `json:"created_at"`
		ExpiredAt int `json:"expired_at"`
		Pools []int `json:"pools"`
		Reason string `json:"reason"`
		AuthorTriggerID int `json:"author_trigger_id"`
		AuthorUserID int `json:"author_user_id"`
		Subnet string `json:"subnet"`
		Country string `json:"country"`
		ProxyType string `json:"proxy_type"`
		Datacenter string `json:"datacenter"`
		SourceValues []string `json:"source_values"`
	}
)

// BlacklistRead requests the current blacklist for the future purposes.
// It is going to respond with the list of IP addresses.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) BlacklistRead(clientID int) ([]IPRule, error) {
	uri := "/v4/ip_rules"

	q := url.Values{}
	q.Set("filter[clientid]", strconv.Itoa(clientID))
	q.Set("filter[list]", "black")
	q.Set("limit", "1000")
	q.Set("offset", "0")

	var bulkIPRules struct {
		Body struct {
			Objects []IPRule `json:"objects"`
		} `json:"body"`
	}

	var result []IPRule

	for offset := 0; len(bulkIPRules.Body.Objects) > 0; offset += 1000 {
		q.Set("offset", strconv.Itoa(offset))

		respBody, err := api.makeRequest("GET", uri, "", q.Encode())
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(respBody, &bulkIPRules); err != nil {
			return nil, err
		}


		result = append(result, bulkIPRules.Body.Objects...)
	}

	return result, nil
}

// BlacklistCreate creates a blacklist in the Wallarm Cloud.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) BlacklistCreate(clientID int, params IPRuleCreationParams) error {
	uri := "/v4/ip_rules"
	reqBody := struct {
		ClientID int `json:"clientid"`
		Force bool `json:"force"`
		IPRule IPRuleCreationParams `json:"ip_rule"`
	}{ClientID: clientID, Force: false, IPRule: params}

	_, err := api.makeRequest("POST", uri, "", &reqBody)

	return err
}

// BlacklistDelete deletes a blacklist for the client.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) BlacklistDelete(clientID int, ids []int) error {
	uri := "/v4/ip_rules"
	reqBody := struct {
		Filter struct {
			ID       []int `json:"id"`
			ClientID int   `json:"clientid"`
		} `json:"filter"`
	}{}
	reqBody.Filter.ID = ids
	reqBody.Filter.ClientID = clientID

	_, err := api.makeRequest("DELETE", uri, "ip_rules", &reqBody)

	return err
}
