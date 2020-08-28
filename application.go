package wallarm

import "encoding/json"

// AppCreate is a request body to set ID and Name for the App.
type AppCreate struct {
	*AppFilter
	Name string `json:"name"`
}

// AppFilter is used to filter applications by ID and ClientID.
type AppFilter struct {
	ID       int `json:"id"`
	Clientid int `json:"clientid"`
}

// AppDelete is a root object for deleting filter.
type AppDelete struct {
	Filter *AppFilter `json:"filter"`
}

// AppReadFilter is a filter by Client ID.
type AppReadFilter struct {
	Clientid []int `json:"clientid"`
}

// AppRead is a request body for filtration of the App.
type AppRead struct {
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Filter *AppReadFilter `json:"filter"`
}

// AppReadResp is a response with parameters of the application.
type AppReadResp struct {
	Status int `json:"status"`
	Body   []struct {
		*AppCreate
		Deleted bool `json:"deleted"`
	} `json:"body"`
}

// AppUpdate is a request body to update Fields in the existing Application defined by Filter.
type AppUpdate struct {
	Filter *AppUpdateFilter `json:"filter"`
	Fields *AppUpdateFields `json:"fields"`
}

// AppUpdateFilter is used as a filter with ID of the Application.
type AppUpdateFilter struct {
	*AppReadFilter
	ID int `json:"id"`
}

// AppUpdateFields is used as identificator what should be changed in Application.
// Only Name is supported.
type AppUpdateFields struct {
	Name string `json:"name"`
}

// AppRead reads Applications and returns params of the Application.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) AppRead(appBody *AppRead) (*AppReadResp, error) {
	uri := "/v1/objects/pool"
	respBody, err := api.makeRequest("POST", uri, "app", appBody)
	if err != nil {
		return nil, err
	}
	var a AppReadResp
	if err = json.Unmarshal(respBody, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// AppCreate returns nothing if Application has been created successfully, otherwise error.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) AppCreate(appBody *AppCreate) error {

	uri := "/v1/objects/pool/create"
	_, err := api.makeRequest("POST", uri, "app", appBody)
	if err != nil {
		return err
	}
	return nil
}

// AppDelete returns nothing if Application has been deleted successfully, otherwise error.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) AppDelete(appBody *AppDelete) error {

	uri := "/v1/objects/pool/delete"
	_, err := api.makeRequest("POST", uri, "app", appBody)
	if err != nil {
		return err
	}
	return nil
}

// AppUpdate returns nothing if the Application has been updated successfully, otherwise error.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) AppUpdate(appBody *AppUpdate) error {

	uri := "/v1/objects/pool/update"
	_, err := api.makeRequest("POST", uri, "app", appBody)
	if err != nil {
		return err
	}
	return nil
}
