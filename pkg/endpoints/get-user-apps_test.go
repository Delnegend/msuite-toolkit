package endpoints_test

import (
	"msuite-toolkit/pkg/endpoints"
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUserApps(t *testing.T) {
	const TEST_USER_ID = "68ff97e403eb7ce1b9967dcf"

	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}
	t.Logf("Loaded config: admin_portal=%s, admin_user_id=%s, bearer_token_len=%d",
		appState.AdminPortalAddress, appState.AdminUserID, len(appState.BearerToken))

	apps, err := endpoints.GetUserApps(&appState, TEST_USER_ID)
	if err != nil {
		t.Fatalf("GetUserApps failed: %v", err)
	}
	t.Logf("User Apps: %+v", apps)
}
