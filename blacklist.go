package wallarm

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Bulk is used to define IP address, applications, time and reason
type Bulk struct {
	IP       string `json:"ip"`
	ExpireAt int    `json:"expire_at"`
	Reason   string `json:"reason"`
	Poolid   int    `json:"poolid"`
	Clientid int    `json:"clientid"`
}

// BlacklistCreate is a root object to fill the blacklist
type BlacklistCreate struct {
	Bulks *[]Bulk `json:"bulk"`
}

// BlacklistRead is used to unmarshal blacklist Read function
type BlacklistRead struct {
	Body struct {
		Result            string      `json:"result"`
		Total             int         `json:"total"`
		Continuation      interface{} `json:"continuation"`
		EventContinuation string      `json:"event_continuation"`
		Objects           []struct {
			ID           int           `json:"id"`
			Clientid     int           `json:"clientid"`
			Country      string        `json:"country"`
			Poolid       int           `json:"poolid"`
			StillAttacks bool          `json:"still_attacks"`
			IP           string        `json:"ip"`
			ExpireAt     int           `json:"expire_at"`
			Tags         []interface{} `json:"tags"`
			BlockedAt    int           `json:"blocked_at"`
			Reason       string        `json:"reason"`
			Tor          interface{}   `json:"tor"`
			Datacenter   interface{}   `json:"datacenter"`
			ProxyType    interface{}   `json:"proxy_type"`
		} `json:"objects"`
	} `json:"body"`
}

// BlacklistRead requests the current blacklist for the future purposes.
// It is going to respond with the list of IP addresses.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) BlacklistRead(clientID int) ([]BlacklistRead, error) {
	uri := "/v3/blacklist"

	q := url.Values{}
	q.Add("filter[clientid]", strconv.Itoa(clientID))
	q.Add("filter[attack_delay]", "300")
	q.Add("limit", "1000")
	query := q.Encode()

	respBody, err := api.makeRequest("GET", uri, "", query)
	if err != nil {
		return nil, err
	}

	var data BlacklistRead
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	var resp []BlacklistRead
	resp = append(resp, data)

	for {
		var data BlacklistRead
		if err = json.Unmarshal(respBody, &data); err != nil {
			return nil, err
		}
		if data.Body.Continuation != nil {
			if q.Get("continuation") == "" {
				q.Add("continuation", data.Body.Continuation.(string))
			} else {
				q.Set("continuation", data.Body.Continuation.(string))
			}
			query = q.Encode()

			respBody, err = api.makeRequest("GET", uri, "", query)
			if err != nil {
				return nil, err
			}
			resp = append(resp, data)
		} else {
			return resp, nil
		}
	}
}

// BlacklistCreate creates a blacklist in the Wallarm Cloud.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) BlacklistCreate(blacklistBody *BlacklistCreate) error {

	uri := "/v3/blacklist/bulk"
	_, err := api.makeRequest("POST", uri, "", blacklistBody)
	if err != nil {
		return err
	}
	return nil
}

// BlacklistDelete deletes a blacklist for the client.
// Currently, it will flush the entire blacklist, then it will be changed for granular deletion.
// API reference: https://apiconsole.eu1.wallarm.com
func (api *API) BlacklistDelete(clientID int, ids []int) error {
	var rest int

	uri := "/v3/blacklist/all"

	q := url.Values{}
	q.Add("filter[clientid]", strconv.Itoa(clientID))
	for k, id := range ids {
		lengthID := len([]byte(strconv.Itoa(id)))
		lengthQuery := len([]byte(q.Encode()))
		if lengthQuery+lengthID > 7000 {
			rest = k
			break
		}
		q.Add("filter[id][]", strconv.Itoa(id))
	}
	query := q.Encode()

	_, err := api.makeRequest("DELETE", uri, "", query)
	if err != nil {
		return err
	}

	idRest := ids
	for rest != 0 {
		q := url.Values{}
		q.Add("filter[clientid]", strconv.Itoa(clientID))
		idRest = idRest[rest:]
		for k, id := range idRest {
			lengthID := len([]byte(strconv.Itoa(id)))
			lengthQuery := len([]byte(q.Encode()))
			if lengthQuery+lengthID > 7000 {
				rest = k
				break
			}
			q.Add("filter[id][]", strconv.Itoa(id))

			if k == len(idRest)-1 {
				rest = 0
			}
		}
		query := q.Encode()

		_, err := api.makeRequest("DELETE", uri, "", query)
		if err != nil {
			return err
		}
	}

	return nil
}
