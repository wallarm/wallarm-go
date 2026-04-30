package wallarm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (

	// Action contains operations available on Action resource
	Action interface {
		HintRead(hintBody *HintRead) (*HintReadResp, error)
		ActionList(params *ActionListParams) (*ActionListResponse, error)
		ActionReadByID(actionID int) (*ActionEntry, error)
		ActionReadByHitID(hitID []string) (*ActionByHitResponse, error)
		HintCreate(ruleBody *ActionCreate) (*ActionCreateResp, error)
		HintUpdateV3(ruleID int, hintBody *HintUpdateV3Params) (*ActionCreateResp, error)
		HintDelete(hintbody *HintDelete) (*HintDeleteResp, error)
	}

	// HintUpdateV3Params is used for updating a hint via v3 API.
	// HintUpdateV3Params is the body sent to PUT /v3/hint/{rule_id}. Pointer
	// types + omitempty let callers send only the fields they want to update;
	// nil fields are dropped from the wire payload.
	HintUpdateV3Params struct {
		// Common to every rule type.
		VariativityDisabled *bool   `json:"variativity_disabled,omitempty"`
		Comment             *string `json:"comment,omitempty"`
		Title               *string `json:"title,omitempty"`
		Active              *bool   `json:"active,omitempty"`
		Set                 *string `json:"set,omitempty"`

		// Mode is mutable for wallarm_mode, api_abuse_mode, graphql_detection,
		// file_upload_size_limit, overlimit_res_settings, set_response_header.
		Mode *string `json:"mode,omitempty"`

		// vpatch / disable_attack_type.
		AttackType *string `json:"attack_type,omitempty"`

		// disable_stamp.
		Stamp *int `json:"stamp,omitempty"`

		// regex.
		Regex *string `json:"regex,omitempty"`

		// credentials_regex.
		LoginRegex    *string `json:"login_regex,omitempty"`
		CaseSensitive *bool   `json:"case_sensitive,omitempty"`

		// credentials_regex / credentials_point.
		CredStuffType *string `json:"cred_stuff_type,omitempty"`

		// file_upload_size_limit.
		Size *int `json:"size,omitempty"`

		// graphql_detection.
		MaxDepth       *int  `json:"max_depth,omitempty"`
		MaxValueSizeKb *int  `json:"max_value_size_kb,omitempty"`
		MaxDocSizeKb   *int  `json:"max_doc_size_kb,omitempty"`
		MaxDocPerBatch *int  `json:"max_doc_per_batch,omitempty"`
		Introspection  *bool `json:"introspection,omitempty"`
		DebugEnabled   *bool `json:"debug_enabled,omitempty"`

		// overlimit_res_settings.
		OverlimitTime *int `json:"overlimit_time,omitempty"`

		// parser_state.
		Parser *string `json:"parser,omitempty"`
		State  *string `json:"state,omitempty"`

		// rate_limit.
		Delay     *int    `json:"delay,omitempty"`
		Burst     *int    `json:"burst,omitempty"`
		Rate      *int    `json:"rate,omitempty"`
		RspStatus *int    `json:"rsp_status,omitempty"`
		TimeUnit  *string `json:"time_unit,omitempty"`

		// set_response_header.
		Name   *string   `json:"name,omitempty"`
		Values *[]string `json:"values,omitempty"`

		// uploads.
		FileType *string `json:"file_type,omitempty"`

		// Mitigation controls (bola, brute, enum, forced_browsing, rate_limit_enum).
		Threshold            *Threshold            `json:"threshold,omitempty"`
		Reaction             *Reaction             `json:"reaction,omitempty"`
		EnumeratedParameters *EnumeratedParameters `json:"enumerated_parameters,omitempty"`
	}

	// ActionDetails defines the Action of how to parse the request.
	// Point represents a part of the request where the condition should be satisfied.
	// ActionDetails is used to define the particular assets of the Action field.
	ActionDetails struct {
		Type  string        `json:"type,omitempty"`
		Point []interface{} `json:"point,omitempty"`
		Value interface{}   `json:"value,omitempty"`
	}

	// ActionCreate is a creation skeleton for the Rule.
	ActionCreate struct {
		Type                 string                  `json:"type"`
		Action               *[]ActionDetails        `json:"action,omitempty"`
		Clientid             int                     `json:"clientid,omitempty"`
		Validated            bool                    `json:"validated"`
		Point                TwoDimensionalSlice     `json:"point,omitempty"`
		Rules                []string                `json:"rules,omitempty"`
		AttackType           string                  `json:"attack_type,omitempty"`
		Stamp                int                     `json:"stamp,omitempty"`
		Mode                 string                  `json:"mode,omitempty"`
		Counter              string                  `json:"counter,omitempty"`
		Regex                string                  `json:"regex,omitempty"`
		RegexID              int                     `json:"regex_id,omitempty"`
		Enabled              *bool                   `json:"enabled,omitempty"`
		Name                 string                  `json:"name,omitempty"`
		Values               []string                `json:"values,omitempty"`
		Comment              string                  `json:"comment,omitempty"`
		FileType             string                  `json:"file_type,omitempty"`
		Parser               string                  `json:"parser,omitempty"`
		State                string                  `json:"state,omitempty"`
		VarType              string                  `json:"var_type,omitempty"`
		VariativityDisabled  bool                    `json:"variativity_disabled,omitempty"`
		LoginRegex           string                  `json:"login_regex,omitempty"`
		CredStuffType        string                  `json:"cred_stuff_type,omitempty"`
		CredStuffMode        string                  `json:"cred_stuff_mode,omitempty"`
		CaseSensitive        *bool                   `json:"case_sensitive,omitempty"`
		LoginPoint           TwoDimensionalSlice     `json:"login_point,omitempty"`
		// Delay/Burst/Rate/OverlimitTime are *int with omitempty so callers
		// can transmit a literal 0 (a valid value per the API range
		// constraints). With non-pointer int+omitempty, encoding/json drops
		// the zero value silently and the API rejects with "can't be blank".
		// RspStatus stays a plain int because its valid range (400..599) does
		// not include 0; omitempty correctly drops absent.
		Delay                *int                    `json:"delay,omitempty"`
		Burst                *int                    `json:"burst,omitempty"`
		Rate                 *int                    `json:"rate,omitempty"`
		RspStatus            int                     `json:"rsp_status,omitempty"`
		TimeUnit             string                  `json:"time_unit,omitempty"`
		OverlimitTime        *int                    `json:"overlimit_time,omitempty"`
		Suffix               string                  `json:"suffix,omitempty"`
		MaxDepth             int                     `json:"max_depth,omitempty"`
		MaxValueSizeKb       int                     `json:"max_value_size_kb,omitempty"`
		MaxDocSizeKb         int                     `json:"max_doc_size_kb,omitempty"`
		MaxAliasesSizeKb     int                     `json:"max_aliases,omitempty"`
		MaxDocPerBatch       int                     `json:"max_doc_per_batch,omitempty"`
		Introspection        *bool                   `json:"introspection,omitempty"`
		DebugEnabled         *bool                   `json:"debug_enabled,omitempty"`
		Set                  string                  `json:"set,omitempty"`
		Active               bool                    `json:"active"`
		Title                string                  `json:"title,omitempty"`
		Mitigation           string                  `json:"mitigation,omitempty"`
		Reaction             *Reaction               `json:"reaction,omitempty"`
		Threshold            *Threshold              `json:"threshold,omitempty"`
		EnumeratedParameters *EnumeratedParameters   `json:"enumerated_parameters,omitempty"`
		AdvancedConditions   []AdvancedCondition     `json:"advanced_conditions,omitempty"`
		ArbitraryConditions  []ArbitraryConditionReq `json:"arbitrary_conditions,omitempty"`
		Size                 int                     `json:"size,omitempty"`
		SizeUnit             string                  `json:"size_unit,omitempty"`
	}

	EnumeratedParameters struct {
		Mode                 string    `json:"mode"`
		NameRegexps          []string  `json:"name_regexps,omitempty"`
		ValueRegexp          []string  `json:"value_regexps,omitempty"`
		AdditionalParameters *bool     `json:"additional_parameters,omitempty"`
		PlainParameters      *bool     `json:"plain_parameters,omitempty"`
		Points               []*Points `json:"points,omitempty"`
	}

	AdvancedCondition struct {
		Field    string   `json:"field"`
		Operator string   `json:"operator"`
		Value    []string `json:"value"`
	}

	ArbitraryConditionReq struct {
		Point    TwoDimensionalSlice `json:"point"`
		Operator string              `json:"operator"`
		Value    []string            `json:"value"`
	}

	ArbitraryConditionResp struct {
		Point    []interface{} `json:"point"`
		Operator string        `json:"operator"`
		Value    []string      `json:"value"`
	}

	Points struct {
		Point     []interface{} `json:"point"`
		Sensitive bool          `json:"sensitive"`
	}

	Reaction struct {
		BlockBySession *int `json:"block_by_session,omitempty"`
		BlockByIP      *int `json:"block_by_ip,omitempty"`
		GraylistByIP   *int `json:"graylist_by_ip,omitempty"`
	}

	Threshold struct {
		Count  int `json:"count"`
		Period int `json:"period"`
	}

	// ActionListFilter is the filter for listing actions.
	ActionListFilter struct {
		ID       []int    `json:"id,omitempty"`
		NotID    []int    `json:"!id,omitempty"`
		Clientid []int    `json:"clientid,omitempty"`
		HintType []string `json:"hint_type,omitempty"`
		Empty    *bool    `json:"empty,omitempty"`
	}

	// TwoDimensionalSlice is used for Point and HintsCount structures.
	TwoDimensionalSlice [][]interface{}

	// ActionListParams is the request body for listing actions.
	ActionListParams struct {
		Filter *ActionListFilter `json:"filter"`
		Limit  int               `json:"limit"`
		Offset int               `json:"offset"`
	}

	// ActionEntry represents a single action from the API.
	ActionEntry struct {
		ID               int             `json:"id"`
		Clientid         int             `json:"clientid"`
		Name             *string         `json:"name"`
		Conditions       []ActionDetails `json:"conditions"`
		EndpointPath     *string         `json:"endpoint_path"`
		EndpointDomain   *string         `json:"endpoint_domain"`
		EndpointInstance *string         `json:"endpoint_instance"`
		UpdatedAt        int             `json:"updated_at"`
	}

	// ActionListResponse is the response from listing actions.
	ActionListResponse struct {
		Status int           `json:"status"`
		Body   []ActionEntry `json:"body"`
	}

	// ActionByHitResponse is the response from fetching action conditions by hit ID.
	ActionByHitResponse struct {
		Status int `json:"status"`
		Body   struct {
			Conditions []ActionDetails `json:"conditions"`
			Clientid   int             `json:"clientid"`
		} `json:"body"`
	}

	// ActionBody is an inner body for the Action and Hint responses.
	ActionBody struct {
		ID            int             `json:"id"`
		ActionID      int             `json:"actionid"`
		Clientid      int             `json:"clientid"`
		Action        []ActionDetails `json:"action"`
		CreateTime    int             `json:"create_time"`
		CreateUserid  int             `json:"create_userid"`
		Validated     bool            `json:"validated"`
		System        bool            `json:"system"`
		RegexID       interface{}     `json:"regex_id"`
		UpdatedAt     int             `json:"updated_at"`
		Type          string          `json:"type"`
		Enabled       bool            `json:"enabled"`
		Mode          string          `json:"mode"`
		Regex         string          `json:"regex"`
		Point         []interface{}   `json:"point"`
		AttackType    string          `json:"attack_type"`
		Stamp         int             `json:"stamp,omitempty"`
		Rules         []string        `json:"rules"`
		Counter       string          `json:"counter,omitempty"`
		VarType       string          `json:"var_type"`
		LoginRegex    string          `json:"login_regex"`
		CredStuffType string          `json:"cred_stuff_type"`
		CredStuffMode string          `json:"cred_stuff_mode"`
		CaseSensitive *bool           `json:"case_sensitive"`
		LoginPoint    []interface{}   `json:"login_point"`
		// Headers for the Set Response Headers Rule
		// are defined by these two parameters.
		Name                 string                   `json:"name"`
		Values               []interface{}            `json:"values"`
		Delay                int                      `json:"delay,omitempty"`
		Burst                int                      `json:"burst,omitempty"`
		Rate                 int                      `json:"rate,omitempty"`
		RspStatus            int                      `json:"rsp_status,omitempty"`
		TimeUnit             string                   `json:"time_unit,omitempty"`
		OverlimitTime        int                      `json:"overlimit_time,omitempty"`
		Suffix               string                   `json:"suffix,omitempty"`
		FileType             string                   `json:"file_type,omitempty"`
		Parser               string                   `json:"parser,omitempty"`
		State                string                   `json:"state,omitempty"`
		MaxDepth             int                      `json:"max_depth,omitempty"`
		MaxValueSizeKb       int                      `json:"max_value_size_kb,omitempty"`
		MaxDocSizeKb         int                      `json:"max_doc_size_kb,omitempty"`
		MaxAliasesSizeKb     int                      `json:"max_aliases,omitempty"`
		MaxDocPerBatch       int                      `json:"max_doc_per_batch,omitempty"`
		Introspection        *bool                    `json:"introspection,omitempty"`
		DebugEnabled         *bool                    `json:"debug_enabled,omitempty"`
		Set                  string                   `json:"set,omitempty"`
		Active               bool                     `json:"active"`
		Title                string                   `json:"title,omitempty"`
		Comment              string                   `json:"comment,omitempty"`
		Mitigation           string                   `json:"mitigation,omitempty"`
		Reaction             *Reaction                `json:"reaction,omitempty"`
		Threshold            *Threshold               `json:"threshold,omitempty"`
		EnumeratedParameters *EnumeratedParameters    `json:"enumerated_parameters,omitempty"`
		AdvancedConditions   []AdvancedCondition      `json:"advanced_conditions,omitempty"`
		ArbitraryConditions  []ArbitraryConditionResp `json:"arbitrary_conditions,omitempty"`
		Size                 int                      `json:"size,omitempty"`
		SizeUnit             string                   `json:"size_unit,omitempty"`
		VariativityDisabled  bool                     `json:"variativity_disabled,omitempty"`
	}

	// ActionCreateResp is the response of just created Rule.
	ActionCreateResp struct {
		Status int         `json:"status"`
		Body   *ActionBody `json:"body"`
	}

	// HintReadResp is the response of filtered rules by Action ID.
	HintReadResp struct {
		Status int           `json:"status"`
		Body   *[]ActionBody `json:"body"`
	}

	// HintRead is used to define whether action of the rule exists.
	HintRead struct {
		Filter    *HintFilter `json:"filter"`
		OrderBy   string      `json:"order_by"`
		OrderDesc bool        `json:"order_desc"`
		Limit     int         `json:"limit"`
		Offset    int         `json:"offset"`
	}

	// HintFilter is used as a filter by Action ID.
	HintFilter struct {
		Clientid        []int               `json:"clientid,omitempty"`
		ActionID        []int               `json:"actionid,omitempty"`
		ID              []int               `json:"id,omitempty"`
		NotID           []int               `json:"!id,omitempty"`
		NotActionID     []int               `json:"!actionid,omitempty"`
		CreateUserid    []int               `json:"create_userid,omitempty"`
		NotCreateUserid []int               `json:"!create_userid,omitempty"`
		CreateTime      [][]int             `json:"create_time,omitempty"`
		NotCreateTime   [][]int             `json:"!create_time,omitempty"`
		System          *bool               `json:"system,omitempty"`
		Type            []string            `json:"type,omitempty"`
		Point           TwoDimensionalSlice `json:"point,omitempty"`
	}

	// HintDelete is used for removal of Rules by Hint IDs.
	// Supports batch delete — up to 1000 IDs per request.
	HintDelete struct {
		Filter *HintDeleteFilter `json:"filter"`
	}

	// HintDeleteFilter is used as a filter by Hint IDs.
	HintDeleteFilter struct {
		Clientid []int `json:"clientid"`
		ID       []int `json:"id"`
	}

	// HintDeleteResp is the response of a hint delete request. Body is empty
	// on the no-op path; see HintDelete.
	HintDeleteResp struct {
		Status int          `json:"status"`
		Body   []ActionBody `json:"body"`
	}
)

// HintRead reads the Rules defined by Action ID.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) HintRead(hintBody *HintRead) (*HintReadResp, error) {

	uri := "/v1/objects/hint"
	respBody, err := api.makeRequest(http.MethodPost, uri, "hint", hintBody, nil)
	if err != nil {
		return nil, err
	}
	var h HintReadResp
	if err = json.Unmarshal(respBody, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// ActionList lists actions matching the given filter.
// Endpoint: POST /v1/objects/action
func (api *api) ActionList(params *ActionListParams) (*ActionListResponse, error) {
	uri := "/v1/objects/action"
	respBody, err := api.makeRequest(http.MethodPost, uri, "action", params, nil)
	if err != nil {
		return nil, err
	}
	var a ActionListResponse
	if err = json.Unmarshal(respBody, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// ActionReadByID fetches a single action by its ID.
// Endpoint: GET /v3/action/{id}
func (api *api) ActionReadByID(actionID int) (*ActionEntry, error) {
	uri := fmt.Sprintf("/v3/action/%d", actionID)
	respBody, err := api.makeRequest(http.MethodGet, uri, "action", nil, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Status int         `json:"status"`
		Body   ActionEntry `json:"body"`
	}
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp.Body, nil
}

// ActionReadByHitID fetches action conditions for a given hit.
// The hit ID is an Elasticsearch tuple [index_name, document_id].
// Returns conditions and client ID but no action ID (the action may not exist yet).
// Endpoint: POST /v1/objects/action/by_hit
func (api *api) ActionReadByHitID(hitID []string) (*ActionByHitResponse, error) {
	uri := "/v1/objects/action/by_hit"
	body := struct {
		ID []string `json:"id"`
	}{ID: hitID}
	respBody, err := api.makeRequest(http.MethodPost, uri, "action", body, nil)
	if err != nil {
		return nil, err
	}
	var resp ActionByHitResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// HintCreate creates Rules in Wallarm Cloud.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) HintCreate(ruleBody *ActionCreate) (*ActionCreateResp, error) {

	uri := "/v1/objects/hint/create"
	respBody, err := api.makeRequest(http.MethodPost, uri, "rule", ruleBody, nil)
	if err != nil {
		return nil, err
	}
	var a ActionCreateResp
	if err = json.Unmarshal(respBody, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// HintDelete deletes the Rule defined by the unique Hint ID.
//
// The Wallarm API returns HTTP 200 in all cases, even when nothing was
// deleted (rule already absent, never existed, or delete blocked server-side
// for counter hints). Callers should inspect the returned Body to detect the
// no-op case (len(Body) == 0).
//
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) HintDelete(hintbody *HintDelete) (*HintDeleteResp, error) {
	uri := "/v1/objects/hint/delete"
	respBody, err := api.makeRequest(http.MethodPost, uri, "hint", hintbody, nil)
	if err != nil {
		return nil, err
	}
	var resp HintDeleteResp
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// HintUpdateV3 updates a hint by rule ID via v3 API.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *api) HintUpdateV3(ruleID int, hintBody *HintUpdateV3Params) (*ActionCreateResp, error) {
	uri := fmt.Sprintf("/v3/hint/%d", ruleID)
	respBody, err := api.makeRequest(http.MethodPut, uri, "hint", hintBody, nil)
	if err != nil {
		return nil, err
	}
	var a ActionCreateResp
	if err = json.Unmarshal(respBody, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
