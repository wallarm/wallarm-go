package wallarm

type Allowlist interface {
	AllowlistRead(clientID int) ([]IPRule, error)
	AllowlistCreate(clientID int, params AccessRuleCreateRequest) error
	AllowlistDelete(clientID int, rules []AccessRuleDeleteEntry) error
}

// AllowlistRead requests the current allowlist.
func (api *api) AllowlistRead(clientID int) ([]IPRule, error) {
	return api.IPListRead(AllowlistType, clientID)
}

// AllowlistCreate creates an allowlist entry in the Wallarm Cloud.
func (api *api) AllowlistCreate(clientID int, params AccessRuleCreateRequest) error {
	params.List = AllowlistType
	return api.IPListCreate(clientID, params)
}

// AllowlistDelete deletes allowlist entries for the client.
func (api *api) AllowlistDelete(clientID int, rules []AccessRuleDeleteEntry) error {
	return api.IPListDelete(clientID, rules)
}
