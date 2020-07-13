package wallarm

import (
	"encoding/json"
	"fmt"
)

// ActionDetails defines the Action of how to parse the request.
// Point represents a part of the request where the condition should be satisfied.
// ActionDetails is used to define the particular assets of the Action field.
type ActionDetails struct {
	Type  string        `json:"type,omitempty"`
	Point []interface{} `json:"point,omitempty"`
	Value string        `json:"value,omitempty"`
}

// ActionCreate is a creation skeleton for the Rule.
type ActionCreate struct {
	Type       string           `json:"type"`
	Action     *[]ActionDetails `json:"action,omitempty"`
	Clientid   int              `json:"clientid,omitempty"`
	Validated  bool             `json:"validated,omitempty"`
	Point      [][]string       `json:"point,omitempty"`
	Rules      []string         `json:"rules,omitempty"`
	AttackType string           `json:"attack_type,omitempty"`
	Mode       string           `json:"mode,omitempty"`
	Regex      string           `json:"regex,omitempty"`
	RegexID    int              `json:"regex_id,omitempty"`
	Enabled    *bool            `json:"enabled,omitempty"`
	Name       string           `json:"name,omitempty"`
	Values     []string         `json:"values,omitempty"`
}

// ActionFilter is the specific filter for getting the rules.
// This is an inner structure.
type ActionFilter struct {
	ID         []int           `json:"id"`
	NotID      []int           `json:"!id"`
	Clientid   []int           `json:"clientid"`
	HintsCount [][]interface{} `json:"hints_count"`
	HintType   []string        `json:"hint_type"`
}

// ActionRead is used as a filter to fetch the rules.
type ActionRead struct {
	Filter *ActionFilter `json:"filter"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

// ActionBody is an inner body for the Action and Hint responses.
type ActionBody struct {
	ID           int             `json:"id"`
	Actionid     int             `json:"actionid"`
	Clientid     int             `json:"clientid"`
	Action       []ActionDetails `json:"action"`
	CreateTime   int             `json:"create_time"`
	CreateUserid int             `json:"create_userid"`
	Validated    bool            `json:"validated"`
	System       bool            `json:"system"`
	RegexID      interface{}     `json:"regex_id"`
	UpdatedAt    int             `json:"updated_at"`
	Type         string          `json:"type"`
	Point        []string        `json:"point"`
	AttackType   string          `json:"attack_type"`
}

// ActionCreateResp is the response of just created Rule.
type ActionCreateResp struct {
	Status int         `json:"status"`
	Body   *ActionBody `json:"body"`
}

// HintReadResp is the response of filtered rules by Action ID.
type HintReadResp struct {
	Status int           `json:"status"`
	Body   *[]ActionBody `json:"body"`
}

// HintRead is used to define whether action of the rule exists.
type HintRead struct {
	Filter    *HintFilter `json:"filter"`
	OrderBy   string      `json:"order_by"`
	OrderDesc bool        `json:"order_desc"`
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
}

// HintFilter is used as a filter by Action ID.
type HintFilter struct {
	Clientid []int `json:"clientid"`
	Actionid []int `json:"actionid"`
}

// HintDelete is used for removal of Rule by Hint ID.
type HintDelete struct {
	Filter *HintDeleteFilter `json:"filter"`
}

// HintDeleteFilter is used as a filter by Hint ID.
type HintDeleteFilter struct {
	Clientid []int `json:"clientid"`
	ID       int   `json:"id"`
}

// HintRead reads the Rules defined by Action ID.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) HintRead(hintBody *HintRead) (*HintReadResp, error) {

	uri := "/v1/objects/hint"
	respBody, err := api.makeRequest("POST", uri, "hint", hintBody)
	if err != nil {
		return nil, err
	}
	var h HintReadResp
	if err = json.Unmarshal(respBody, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// RuleRead reads the Rules defined by a filter.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) RuleRead(ruleBody *ActionRead) error {

	uri := "/v1/objects/action"
	_, err := api.makeRequest("POST", uri, "rule", ruleBody)
	if err != nil {
		return err
	}
	return nil
}

// RuleCreate creates Rules in Wallarm Cloud.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) RuleCreate(ruleBody *ActionCreate) (*ActionCreateResp, error) {

	uri := "/v1/objects/hint/create"
	respBody, err := api.makeRequest("POST", uri, "rule", ruleBody)
	if err != nil {
		return nil, err
	}
	var a ActionCreateResp
	if err = json.Unmarshal(respBody, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// RuleDelete deletes the Rule defined by unique ID.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) RuleDelete(actionID int) error {

	uri := fmt.Sprintf("/v2/action/%d", actionID)
	_, err := api.makeRequest("DELETE", uri, "rule", nil)
	if err != nil {
		return err
	}
	return nil
}

// HintDelete deletes the Rule defined by the unique Hint ID.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) HintDelete(hintbody *HintDelete) error {

	uri := "/v1/objects/hint/delete"
	_, err := api.makeRequest("POST", uri, "hint", hintbody)
	if err != nil {
		return err
	}
	return nil
}
