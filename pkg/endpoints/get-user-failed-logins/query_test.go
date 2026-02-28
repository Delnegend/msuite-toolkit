package get_user_failed_logins

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUserFailedLogins(t *testing.T) {
	const TEST_USER_ID = "6895fe2e5a62f6f14af9d954"

	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}
	t.Logf("Loaded config: admin_portal=%s, admin_user_id=%s, bearer_token_len=%d",
		appState.AdminPortalAddress, appState.AdminUserID, len(appState.BearerToken))

	_, failedLogins, err := GetUserFailedLogins(&appState, TEST_USER_ID, 0, 100)
	if err != nil {
		t.Fatalf("GetUserFailedLogins failed: %v", err)
	}
	t.Logf("User failed logins count: %d", len(failedLogins))
	for i, login := range failedLogins {
		if i < 5 {
			t.Logf("Failed login %d: %+v", i, login)
		}
	}
}
