package wallarm

type Graylist interface {
	GraylistRead(clientID int) ([]IPRule, error)
	GraylistCreate(clientID int, params AccessRuleCreateRequest) error
	GraylistDelete(clientID int, rules []AccessRuleDeleteEntry) error
}

// GraylistRead requests the current graylist.
func (api *api) GraylistRead(clientID int) ([]IPRule, error) {
	return api.IPListRead(GraylistType, clientID)
}

// GraylistCreate creates a graylist entry in the Wallarm Cloud.
func (api *api) GraylistCreate(clientID int, params AccessRuleCreateRequest) error {
	params.List = GraylistType
	return api.IPListCreate(clientID, params)
}

// GraylistDelete deletes graylist entries for the client.
func (api *api) GraylistDelete(clientID int, rules []AccessRuleDeleteEntry) error {
	return api.IPListDelete(clientID, rules)
}
