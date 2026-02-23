package types

type AppState struct {
	AdminPortalAddress string `toml:"admin_portal_address" json:"admin_portal_address"`
	BearerToken        string `toml:"bearer_token" json:"bearer_token"`
	AdminUserID        string `toml:"admin_user_id" json:"admin_user_id"`
}
