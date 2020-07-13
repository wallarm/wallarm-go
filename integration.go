package wallarm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// IntegrationEvents represents `Events` object while creating a new integration.
// Event possible values: "hit", "vuln", "system", "scope".
// If `IntegrationObject.Type` is "opsgenie" possible values: "hit", "vuln".
// `Active` identifies whether the current Event should be reported.
type IntegrationEvents struct {
	Event  string `json:"event"`
	Active bool   `json:"active"`
}

// IntegrationObject is an inner object for the Read function containing.
// ID is a unique identifier of the Integration.
type IntegrationObject struct {
	ID        int         `json:"id"`
	Active    bool        `json:"active"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	CreatedAt int         `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	Target    interface{} `json:"target"`
	Events    []struct {
		Event  string `json:"event"`
		Active bool   `json:"active"`
	} `json:"events"`
}

// IntegrationRead is the response on the Read action.
// This is used for correct Unmarshalling of the response as a container.
type IntegrationRead struct {
	Body struct {
		Result string               `json:"result"`
		Object *[]IntegrationObject `json:"object"`
	} `json:"body"`
}

// IntegrationCreate defines how to configure Integration.
// `Type` possible values: "insight_connect", "opsgenie", "slack",
//  "pager_duty", "splunk", "sumo_logic"
type IntegrationCreate struct {
	Name     string               `json:"name"`
	Active   bool                 `json:"active"`
	Target   string               `json:"target"`
	Events   *[]IntegrationEvents `json:"events"`
	Type     string               `json:"type"`
	Clientid int                  `json:"clientid,omitempty"`
}

// IntegrationCreate returns nothing if Integration has been created succesfully, otherwise Error.
// It accepts a body with defined settings namely Event types, Name, Target.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) IntegrationCreate(integrationBody *IntegrationCreate) error {

	uri := "/v2/integration"
	_, err := api.makeRequest("POST", uri, "integration", integrationBody)
	if err != nil {
		return err
	}
	return nil
}

// IntegrationUpdate is used to Update existing resources.
// It utilises the same format of body as the Create function.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) IntegrationUpdate(integrationBody *IntegrationCreate, integrationID int) error {

	uri := fmt.Sprintf("/v2/integration/%d", integrationID)
	_, err := api.makeRequest("PUT", uri, "integration", integrationBody)
	if err != nil {
		return err
	}
	return nil
}

// IntegrationRead is used to read existing integrations.
// It returns the list of Integrations
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) IntegrationRead(clientID int, name, IntType string) (*IntegrationObject, error) {

	uri := "/v2/integration"
	q := url.Values{}
	q.Add("clientid", strconv.Itoa(clientID))
	query := q.Encode()
	respBody, err := api.makeRequest("GET", uri, "integration", query)
	if err != nil {
		return nil, err
	}
	var i IntegrationRead
	if err = json.Unmarshal(respBody, &i); err != nil {
		return nil, err
	}
	for _, obj := range *i.Body.Object {
		if obj.Name == name && obj.Type == IntType {
			return &obj, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Body: %s", string(respBody)))
}

// IntegrationDelete is used to delete an existing integration.
// If successful, returns nothing.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) IntegrationDelete(integrationID int) error {

	uri := fmt.Sprintf("/v2/integration/%d", integrationID)
	_, err := api.makeRequest("DELETE", uri, "integration", nil)
	if err != nil {
		return err
	}
	return nil
}

// IntegrationWithAPITarget is used to create an Integration with the following parameters.
// On purpose to fulfil a custom Webhooks integration.
type IntegrationWithAPITarget struct {
	Token       string                 `json:"token,omitempty"`
	API         string                 `json:"api,omitempty"`
	URL         string                 `json:"url,omitempty"`
	HTTPMethod  string                 `json:"http_method,omitempty"`
	Headers     map[string]interface{} `json:"headers"`
	CaFile      string                 `json:"ca_file"`
	CaVerify    bool                   `json:"ca_verify"`
	Timeout     int                    `json:"timeout,omitempty"`
	OpenTimeout int                    `json:"open_timeout,omitempty"`
}

// IntegrationWithAPICreate is a root object of Create action for Integrations.
// It aids to set `Events` to trigger this integration.
// `Type` possible values: "web_hooks"
// `Target` is a struct for a Webhooks endpoint containing params such as URL, Token, etc.
type IntegrationWithAPICreate struct {
	Name     string                    `json:"name"`
	Active   bool                      `json:"active"`
	Target   *IntegrationWithAPITarget `json:"target"`
	Events   *[]IntegrationEvents      `json:"events"`
	Type     string                    `json:"type"`
	Clientid int                       `json:"clientid,omitempty"`
}

// IntegrationWithAPICreate returns nothing if Integration has been created succesfully, otherwise - error.
// It accepts defined settings namely Event types, Name, Target.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) IntegrationWithAPICreate(integrationBody *IntegrationWithAPICreate) error {

	uri := "/v2/integration"
	_, err := api.makeRequest("POST", uri, "integration", integrationBody)
	if err != nil {
		return err
	}
	return nil
}

// IntegrationWithAPIUpdate is used to Update existing resources.
// It utilises the same format of body as the Create function.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) IntegrationWithAPIUpdate(integrationBody *IntegrationWithAPICreate, integrationID int) error {

	uri := fmt.Sprintf("/v2/integration/%d", integrationID)
	_, err := api.makeRequest("PUT", uri, "integration", integrationBody)
	if err != nil {
		return err
	}
	return nil
}

// EmailIntegrationCreate is a root object of `Create` action for the `email` integration.
// Temporary workaround (`Target` is a slice instead of string) to not check type many times.
// Then it will be changed to interface{} with type checking
type EmailIntegrationCreate struct {
	Name     string               `json:"name"`
	Active   bool                 `json:"active"`
	Target   []string             `json:"target"`
	Events   *[]IntegrationEvents `json:"events"`
	Type     string               `json:"type"`
	Clientid int                  `json:"clientid,omitempty"`
}

// EmailIntegrationCreate returns nothing if the `email` Integration has been created succesfully, otherwise - error.
// It accepts defined settings namely Event types, Name, Target.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) EmailIntegrationCreate(emailBody *EmailIntegrationCreate) error {

	uri := "/v2/integration"
	_, err := api.makeRequest("POST", uri, "email", emailBody)
	if err != nil {
		return err
	}
	return nil
}

// EmailIntegrationUpdate is used to Update existing resources.
// It utilises the same format of body as the Create function.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) EmailIntegrationUpdate(integrationBody *EmailIntegrationCreate, integrationID int) error {

	uri := fmt.Sprintf("/v2/integration/%d", integrationID)
	_, err := api.makeRequest("PUT", uri, "email", integrationBody)
	if err != nil {
		return err
	}
	return nil
}
