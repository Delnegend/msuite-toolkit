package types

type AppState struct {
	AdminPortalAddress               string `toml:"admin_portal_address" json:"admin_portal_address"`
	BearerToken                      string `toml:"bearer_token" json:"bearer_token"`
	AdminUserID                      string `toml:"admin_user_id" json:"admin_user_id"`
	WorkerCount                      int    `toml:"worker_count" json:"worker_count"`
	DryRun                           bool   `toml:"dry_run" json:"dry_run"`
	IncludeUsersWithUnknownLastLogin bool   `toml:"include_users_with_unknown_last_login" json:"include_users_with_unknown_last_login"`

	FilterBy struct {
		DestinationHost string `toml:"destination_host" json:"destination_host"`
		DestinationPort string `toml:"destination_port" json:"destination_port"`
	} `toml:"filter_by" json:"filter_by"`

	// last_login_threshold_in_month: discard users not logged in within this many months (0 = disabled)
	LastLoginThresholdInMonths int `toml:"last_login_threshold_in_month" json:"last_login_threshold_in_month"`
	// organizational_unit_id: optional OU id to filter users (empty = ignore)
	OrganizationalUnitID string `toml:"organizational_unit_id" json:"organizational_unit_id"`
	// exclude_emails: list of user emails to exclude from deletion operations
	ExcludeEmails []string `toml:"exclude_emails" json:"exclude_emails"`
	// include_emails: list of user emails to include in reports (empty = include all)
	IncludeEmails []string `toml:"include_emails" json:"include_emails"`
}
