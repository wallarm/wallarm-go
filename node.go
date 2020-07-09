package wallarm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// NodeCreate represents options to set for the Node for creating.
type NodeCreate struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
	Clientid int    `json:"clientid"`
}

// NodeCreateResp is the API response on the Create action containing
// information about one concrete node.
// Used to get specific parameters of the created Node such as
// time of last syncronisation along with relevant LOM and Proton files.
type NodeCreateResp struct {
	Status int         `json:"status"`
	Body   GetNodeBody `json:"body"`
}

// GetNodeBodyPOST is used as an additional response on the GET
// request to fetch the statuses for all the Nodes.
type GetNodeBodyPOST struct {
	Type              string      `json:"type"`
	ID                string      `json:"id"`
	IP                string      `json:"ip"`
	Hostname          string      `json:"hostname"`
	LastActivity      int         `json:"last_activity"`
	Enabled           bool        `json:"enabled"`
	Clientid          int         `json:"clientid"`
	LastAnalytic      int         `json:"last_analytic"`
	CreateTime        int         `json:"create_time"`
	CreateFrom        string      `json:"create_from"`
	ProtondbVersion   int         `json:"protondb_version"`
	LomVersion        int         `json:"lom_version"`
	ProtondbUpdatedAt interface{} `json:"protondb_updated_at"`
	LomUpdatedAt      interface{} `json:"lom_updated_at"`
	NodeEnvParams     struct {
		Packages struct {
		} `json:"packages"`
	} `json:"node_env_params"`
	Active bool `json:"active"`
}

// GetNodePOST represents a root object of the fetching of Nodes.
type GetNodePOST struct {
	Status int                `json:"status"`
	Body   *[]GetNodeBodyPOST `json:"body"`
}

// GetNodeBody is used to find out configurations and parameters of one specific Node.
type GetNodeBody struct {
	Type              string      `json:"type"`
	ID                int         `json:"id"`
	UUID              string      `json:"uuid"`
	IP                interface{} `json:"ip"`
	Hostname          string      `json:"hostname"`
	LastActivity      interface{} `json:"last_activity"`
	Enabled           bool        `json:"enabled"`
	Clientid          int         `json:"clientid"`
	LastAnalytic      interface{} `json:"last_analytic"`
	CreateTime        int         `json:"create_time"`
	CreateFrom        string      `json:"create_from"`
	ProtondbVersion   interface{} `json:"protondb_version"`
	LomVersion        interface{} `json:"lom_version"`
	ProtondbUpdatedAt interface{} `json:"protondb_updated_at"`
	LomUpdatedAt      interface{} `json:"lom_updated_at"`
	NodeEnvParams     struct {
		Packages struct {
		} `json:"packages"`
	} `json:"node_env_params"`
	Active              bool   `json:"active"`
	InstanceCount       int    `json:"instance_count"`
	ActiveInstanceCount int    `json:"active_instance_count"`
	Token               string `json:"token"`
	RequestsAmount      int    `json:"requests_amount"`
	Secret              string `json:"secret"`
}

// GetNode represents a root object of the fetching action for Nodes.
// It allows to iterate over several Nodes.
type GetNode struct {
	Status int            `json:"status"`
	Body   *[]GetNodeBody `json:"body"`
}

// NodeFilter is a filter object to convey for the Node request
type NodeFilter struct {
	UUID     string `json:"uuid,omitempty"`
	IP       string `json:"ip,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

// GetNodeReadByFilter is used to fetch Nodes by POST method using filter by UUID/IP/Hostname
type GetNodeReadByFilter struct {
	Filter    *NodeFilter `json:"filter"`
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
	OrderBy   string      `json:"order_by,omitempty"`
	OrderDesc bool        `json:"order_desc,omitempty"`
}

// NodeCreate returns the info about just created Node
// For example, UUID/ClientID/Token/InstanceCount
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) NodeCreate(nodeBody *NodeCreate) (*NodeCreateResp, error) {

	uri := "/v2/node"
	respBody, err := api.makeRequest("POST", uri, "node", nodeBody)
	if err != nil {
		return nil, err
	}
	var n NodeCreateResp
	if err = json.Unmarshal(respBody, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

// NodeDelete returns nothing if request was successful, otherwise Error.
// It accepts a node ID which is used to delete the specified Node
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) NodeDelete(nodeID int) error {

	uri := fmt.Sprintf("/v2/node/%d", nodeID)
	_, err := api.makeRequest("DELETE", uri, "", nil)
	if err != nil {
		return err
	}
	return nil
}

// GetNodeRead returns statistics about 1000 created Nodes with specified type, for instance, `cloud_node`
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) GetNodeRead(clientID int, typeNode string) (*GetNode, error) {

	uri := "/v2/node"
	q := url.Values{}
	q.Add("order_by", "hostname")
	q.Add("filter[clientid][]", strconv.Itoa(clientID))
	q.Add("limit", "1000")
	q.Add("offset", "0")
	if typeNode == "all" {
		q.Add("filter[!type]", "fast_node")
	} else {
		q.Add("filter[type]", typeNode)
	}
	query := q.Encode()

	respBody, err := api.makeRequest("GET", uri, "", query)
	if err != nil {
		return nil, err
	}

	var n GetNode
	if err = json.Unmarshal(respBody, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

// GetNodeReadByFilter returns settings of the Node specified by body with a filter.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) GetNodeReadByFilter(getNodeBody *GetNodeReadByFilter) (*GetNodePOST, error) {

	uri := "/v1/objects/node"
	respBody, err := api.makeRequest("POST", uri, "", getNodeBody)
	if err != nil {
		return nil, err
	}

	var n GetNodePOST
	if err = json.Unmarshal(respBody, &n); err != nil {
		return nil, err
	}
	return &n, nil
}
