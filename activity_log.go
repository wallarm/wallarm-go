package wallarm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	// ActivityLog contains read operations available on Activity Log resources.
	ActivityLog interface {
		ActivityLogEventsGetFilters() (*ActivityLogEventsGetFiltersResp, error)
		ActivityLogEventsRead(req *ActivityLogEventsRead) (*ActivityLogEventsResp, error)
		ActivityLogEventRead(req *ActivityLogEventRead) (*ActivityLogEventsResp, error)
	}

	ActivityLogEventsGetFiltersResp struct {
		Body ActivityLogEventsGetFiltersRespBody `json:"body"`
	}

	ActivityLogEventsGetFiltersRespBody struct {
		ObjectTypes []ActivityLogValueLabel `json:"object_types"`
		ActionTypes []ActivityLogValueLabel `json:"action_types"`
		Outcomes    []string                `json:"outcomes"`
		Sources     []string                `json:"sources"`
	}

	ActivityLogValueLabel struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	ActivityLogEventsRead struct {
		ClientID  int                      `json:"-"`
		Filter    *ActivityLogEventsFilter `json:"filter,omitempty"`
		Offset    int                      `json:"offset,omitempty"`
		Limit     int                      `json:"limit,omitempty"`
		OrderBy   string                   `json:"order_by,omitempty"`
		OrderDesc bool                     `json:"order_desc,omitempty"`
	}

	ActivityLogEventsFilter struct {
		ObjectTypes []string `json:"object_types,omitempty"`
		ActionTypes []string `json:"action_types,omitempty"`
		Outcomes    []string `json:"outcomes,omitempty"`
		Sources     []string `json:"sources,omitempty"`
		ActorIDs    []int64  `json:"actor_ids,omitempty"`
		TimeStart   int64    `json:"time_start,omitempty"`
		TimeEnd     int64    `json:"time_end,omitempty"`
	}

	ActivityLogEventRead struct {
		ClientID int   `json:"-"`
		EventID  int64 `json:"-"`
	}

	ActivityLogEventsResp struct {
		Body ActivityLogEventsRespBody `json:"body"`
	}

	ActivityLogEventsRespBody struct {
		Objects []ActivityLogEvent `json:"objects"`
	}

	ActivityLogEvent struct {
		ID                uint64                 `json:"id"`
		Time              int64                  `json:"time"`
		ActionType        string                 `json:"action_type"`
		ObjectType        string                 `json:"object_type"`
		ObjectTypeInfo    *ActivityLogValueLabel `json:"object_type_info,omitempty"`
		Outcome           string                 `json:"outcome"`
		Source            string                 `json:"source"`
		ClientID          int64                  `json:"client_id"`
		Initiator         *ActivityLogActor      `json:"initiator,omitempty"`
		ObjectID          *string                `json:"object_id,omitempty"`
		Description       *string                `json:"description,omitempty"`
		ChangedFields     []string               `json:"changed_fields,omitempty"`
		Diff              json.RawMessage        `json:"diff,omitempty"`
		StateBeforeAction json.RawMessage        `json:"state_before_action,omitempty"`
		StateAfterAction  json.RawMessage        `json:"state_after_action,omitempty"`
		Object            ActivityLogObject      `json:"object"`
	}

	ActivityLogActor struct {
		ID       *int64  `json:"id,omitempty"`
		ClientID *int64  `json:"client_id,omitempty"`
		Name     *string `json:"name,omitempty"`
	}

	ActivityLogObject struct {
		Name string `json:"name,omitempty"`
		Info string `json:"info,omitempty"`
	}
)

func (api *api) ActivityLogEventsGetFilters() (*ActivityLogEventsGetFiltersResp, error) {
	respBody, err := api.makeRequest(http.MethodGet, "/v1/activity_log/events_get_filters", "activity_log", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp ActivityLogEventsGetFiltersResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) ActivityLogEventsRead(req *ActivityLogEventsRead) (*ActivityLogEventsResp, error) {
	if req == nil {
		return nil, fmt.Errorf("activity log request is required")
	}

	uri := fmt.Sprintf("/v1/activity_log/%d/events", req.ClientID)
	respBody, err := api.makeRequest(http.MethodGet, uri, "activity_log", req.toQuery(), nil)
	if err != nil {
		return nil, err
	}

	var resp ActivityLogEventsResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (api *api) ActivityLogEventRead(req *ActivityLogEventRead) (*ActivityLogEventsResp, error) {
	if req == nil {
		return nil, fmt.Errorf("activity log event request is required")
	}

	uri := fmt.Sprintf("/v1/activity_log/%d/events/%d", req.ClientID, req.EventID)
	respBody, err := api.makeRequest(http.MethodGet, uri, "activity_log", "", nil)
	if err != nil {
		return nil, err
	}

	var resp ActivityLogEventsResp
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (req *ActivityLogEventsRead) toQuery() string {
	v := url.Values{}
	if req == nil {
		return v.Encode()
	}

	if req.Filter != nil {
		for _, value := range req.Filter.ObjectTypes {
			if value != "" {
				v.Add("filter[object_types][]", value)
			}
		}
		for _, value := range req.Filter.ActionTypes {
			if value != "" {
				v.Add("filter[action_types][]", value)
			}
		}
		for _, value := range req.Filter.Outcomes {
			if value != "" {
				v.Add("filter[outcomes][]", value)
			}
		}
		for _, value := range req.Filter.Sources {
			if value != "" {
				v.Add("filter[sources][]", value)
			}
		}
		for _, actorID := range req.Filter.ActorIDs {
			v.Add("filter[actor_ids][]", strconv.FormatInt(actorID, 10))
		}
		if req.Filter.TimeStart > 0 {
			v.Set("filter[time_start]", strconv.FormatInt(req.Filter.TimeStart, 10))
		}
		if req.Filter.TimeEnd > 0 {
			v.Set("filter[time_end]", strconv.FormatInt(req.Filter.TimeEnd, 10))
		}
	}

	if req.Offset > 0 {
		v.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		v.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.OrderBy != "" {
		v.Set("order_by", req.OrderBy)
	}
	v.Set("order_desc", strconv.FormatBool(req.OrderDesc))

	return v.Encode()
}
