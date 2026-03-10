package wallarm

type Denylist interface {
	DenylistRead(clientID int) ([]IPRule, error)
	DenylistCreate(clientID int, params AccessRuleCreateRequest) error
	DenylistDelete(clientID int, rules []AccessRuleDeleteEntry) error
}

// DenylistRead requests the current denylist.
func (api *api) DenylistRead(clientID int) ([]IPRule, error) {
	return api.IPListRead(DenylistType, clientID)
}

// DenylistCreate creates a denylist entry in the Wallarm Cloud.
func (api *api) DenylistCreate(clientID int, params AccessRuleCreateRequest) error {
	params.List = DenylistType
	return api.IPListCreate(clientID, params)
}

// DenylistDelete deletes denylist entries for the client.
func (api *api) DenylistDelete(clientID int, rules []AccessRuleDeleteEntry) error {
	return api.IPListDelete(clientID, rules)
}
