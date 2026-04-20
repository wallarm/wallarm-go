package wallarm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Attack interface {
		AttackRead(body *AttackReadRequest) (*AttackReadResp, error)
		AttackCount(body *AttackCountRequest) (*AttackCountResp, error)
		AttackIP(body *AttackIPRequest) (*AttackIPResp, error)
		HitDetails(body *HitDetailsRequest) (*HitDetailsResp, error)
		HitRaw(body *HitRawRequest) ([]byte, error)
	}

	AttackReadRequest struct {
		Filter    *AttackFilter `json:"filter"`
		Limit     int           `json:"limit"`
		Offset    int           `json:"offset"`
		OrderBy   string        `json:"order_by"`
		OrderDesc bool          `json:"order_desc"`
		Cursor    string        `json:"cursor,omitempty"`
		Paging    bool          `json:"paging,omitempty"`
	}

	AttackFilter struct {
		ClientID        []int           `json:"clientid,omitempty"`
		Time            [][]interface{} `json:"time,omitempty"`
		Type            []string        `json:"type,omitempty"`
		NotType         []string        `json:"!type,omitempty"`
		Domain          []string        `json:"domain,omitempty"`
		Path            []string        `json:"path,omitempty"`
		Parameter       []string        `json:"parameter,omitempty"`
		PoolID          []int           `json:"poolid,omitempty"`
		IP              []string        `json:"ip,omitempty"`
		AttackID        interface{}     `json:"attackid,omitempty"`
		RequestID       []string        `json:"request_id,omitempty"`
		VulnID          []int           `json:"vulnid,omitempty"`
		NotVulnID       interface{}     `json:"!vulnid,omitempty"`
		Experimental    *bool           `json:"experimental,omitempty"`
		NotExperimental *bool           `json:"!experimental,omitempty"`
		Threat          []int           `json:"threat,omitempty"`
		StatusCode      []int           `json:"statuscode,omitempty"`
		Method          []string        `json:"method,omitempty"`
	}

	AttackBody struct {
		ID            []string    `json:"id"`
		AttackID      string      `json:"attackid"`
		ClientID      int         `json:"clientid"`
		Domain        string      `json:"domain"`
		PoolID        int         `json:"poolid"`
		Method        string      `json:"method"`
		Parameter     string      `json:"parameter"`
		Path          string      `json:"path"`
		Type          string      `json:"type"`
		FirstTime     int         `json:"first_time"`
		LastTime      int         `json:"last_time"`
		Hits          int         `json:"hits"`
		HitsCount     int         `json:"hits_count_by_filter"`
		IPCount       int         `json:"ip_count"`
		IPCountFilter int         `json:"ip_count_by_filter"`
		StatusCodes   []int       `json:"statuscodes"`
		Threat        int         `json:"threat"`
		VulnID        interface{} `json:"vulnid"`
		Target        string      `json:"target"`
		VectorsCount  int         `json:"vectors_count"`
		BlockStatus   string      `json:"block_status"`

		IPTop []struct {
			IP      string `json:"ip"`
			Count   int    `json:"count"`
			Country string `json:"country"`
		} `json:"ip_top"`

		CountryTop []struct {
			Country string `json:"country"`
			Count   int    `json:"count"`
		} `json:"country_top"`

		RecheckStatus string      `json:"recheck_status"`
		Experimental  interface{} `json:"experimental"`
		State         interface{} `json:"state"`
	}

	AttackReadResp struct {
		Status int          `json:"status"`
		Body   []AttackBody `json:"body"`
		Cursor string       `json:"cursor,omitempty"`
	}

	AttackCountRequest struct {
		Filter *AttackCountFilter `json:"filter"`
	}

	AttackCountFilter struct {
		ClientID   []int           `json:"clientid,omitempty"`
		Time       [][]interface{} `json:"time,omitempty"`
		Domain     []string        `json:"domain,omitempty"`
		Path       []string        `json:"path,omitempty"`
		Parameter  []string        `json:"parameter,omitempty"`
		Type       []string        `json:"type,omitempty"`
		RequestID  []string        `json:"request_id,omitempty"`
		AttackID   interface{}     `json:"attackid,omitempty"`
		StatusCode []int           `json:"statuscode,omitempty"`
		IP         []string        `json:"ip,omitempty"`
		ID         []string        `json:"id,omitempty"`
	}

	AttackCountResp struct {
		Status int `json:"status"`
		Body   struct {
			Attacks int     `json:"attacks"`
			Hits    float64 `json:"hits"`
			IPs     int     `json:"ips"`
		} `json:"body"`
	}

	AttackIPRequest struct {
		Filter *AttackCountFilter `json:"filter"`
	}

	AttackIPResp struct {
		Status int      `json:"status"`
		Body   []string `json:"body"`
	}

	HitDetailsRequest struct {
		Filter  *HitFilter `json:"filter"`
		Returns string     `json:"returns,omitempty"`
	}

	HitDetailsResp struct {
		Status int              `json:"status"`
		Body   []HitDetailsBody `json:"body"`
	}

	HitDetailsBody struct {
		ID  []string     `json:"id"`
		Raw HitRawDetail `json:"raw"`
	}

	HitRawDetail struct {
		Method  string                 `json:"method"`
		URI     string                 `json:"uri"`
		Proto   string                 `json:"proto"`
		Headers map[string]interface{} `json:"headers"`
	}

	HitRawRequest struct {
		Filter *HitFilter `json:"filter"`
	}
)

func (api *api) AttackRead(body *AttackReadRequest) (*AttackReadResp, error) {
	if body == nil {
		return nil, fmt.Errorf("attack read request is required")
	}
	uri := "/v1/objects/attack"
	respBody, err := api.makeRequest(http.MethodPost, uri, "attack", body,
		map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var resp AttackReadResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) AttackCount(body *AttackCountRequest) (*AttackCountResp, error) {
	if body == nil {
		return nil, fmt.Errorf("attack count request is required")
	}
	uri := "/v1/objects/attack/count"
	respBody, err := api.makeRequest(http.MethodPost, uri, "attack", body,
		map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var resp AttackCountResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) AttackIP(body *AttackIPRequest) (*AttackIPResp, error) {
	if body == nil {
		return nil, fmt.Errorf("attack ip request is required")
	}
	uri := "/v1/objects/attack/ip"
	respBody, err := api.makeRequest(http.MethodPost, uri, "attack", body,
		map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var resp AttackIPResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) HitDetails(body *HitDetailsRequest) (*HitDetailsResp, error) {
	if body == nil {
		return nil, fmt.Errorf("hit details request is required")
	}
	uri := "/v1/objects/hit/details"
	respBody, err := api.makeRequest(http.MethodPost, uri, "attack", body,
		map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var resp HitDetailsResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) HitRaw(body *HitRawRequest) ([]byte, error) {
	if body == nil {
		return nil, fmt.Errorf("hit raw request is required")
	}
	uri := "/v1/objects/hit/raw"
	return api.makeRequest(http.MethodPost, uri, "attack", body,
		map[string]string{"Content-Type": "application/json"})
}
