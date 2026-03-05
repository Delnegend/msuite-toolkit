package get_user_mfa

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUserMFA(t *testing.T) {
	const TEST_USER_ID = "68ff97e403eb7ce1b9967dcf"

	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}
	t.Logf("Loaded config: admin_portal=%s, admin_user_id=%s, bearer_token_len=%d",
		appState.AdminPortalAddress, appState.AdminUserID, len(appState.BearerToken))

	userMFA, err := GetUserMFA(&appState, TEST_USER_ID)
	if err != nil {
		t.Fatalf("GetUserMFA failed: %v", err)
	}
	t.Logf("User MFA info: %+v", userMFA)
}

func TestGetUsersMFAWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	users := []types.UserInfo{
		{UserID: "68ff97e403eb7ce1b9967dcf"},
	}

	mfaMap := GetUsersMFAWithProgress(&appState, users)
	t.Logf("Fetched MFA for %d users. Map size: %d", len(users), len(mfaMap))
}
