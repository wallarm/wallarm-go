package wallarm

// ClientFields defines fields which are subject to update.
type ClientFields struct {
	Mode        string `json:"mode,omitempty"`
	ScannerMode string `json:"scanner_mode,omitempty"`
}

// ClientFilter is used for filtration.
// ID is a Client ID entity.
type ClientFilter struct {
	ID int `json:"id"`
}

// ClientUpdateCreate is a root object for updating.
type ClientUpdateCreate struct {
	Filter *ClientFilter `json:"filter"`
	Fields *ClientFields `json:"fields"`
}

// ClientUpdateCreate changes client state.
// It can be used with global WAF mode, Scanner, Attack Rechecker Statuses.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) ClientUpdateCreate(updatebody *ClientUpdateCreate) error {

	uri := "/v1/objects/client/update"
	_, err := api.makeRequest("POST", uri, "client", updatebody)
	if err != nil {
		return err
	}
	return nil
}
