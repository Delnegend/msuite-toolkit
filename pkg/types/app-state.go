package types

type AppState struct {
	AdminPortalAddress string `toml:"admin_portal_address" json:"admin_portal_address"`
	BearerToken        string `toml:"bearer_token" json:"bearer_token"`
	AdminUserID        string `toml:"admin_user_id" json:"admin_user_id"`
	WorkerCount        int    `toml:"worker_count" json:"worker_count"`

	// last_login_threshold_in_month: discard users not logged in within this many months (0 = disabled)
	LastLoginThresholdInMonths int `toml:"last_login_threshold_in_month" json:"last_login_threshold_in_month"`
	// organizational_unit_id: optional OU id to filter users (empty = ignore)
	OrganizationalUnitID string `toml:"organizational_unit_id" json:"organizational_unit_id"`
}
