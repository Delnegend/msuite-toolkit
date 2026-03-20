package get_user_device_last_ip

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUserDeviceLastIP(t *testing.T) {
	const TEST_USER_ID = "6895fe2e5a62f6f14af9d954"
	const TEST_DEVICE_ID = "00000000-0000-0000-0000-000000000000-9adbfe72-ca4e-4b7f-9c3b-262685c26741"

	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}
	t.Logf("Loaded config: admin_portal=%s, admin_user_id=%s, bearer_token_len=%d",
		appState.AdminPortalAddress, appState.AdminUserID, len(appState.BearerToken))

	ip, err := GetUserDeviceLastIP(&appState, TEST_USER_ID, TEST_DEVICE_ID)
	if err != nil {
		t.Fatalf("GetUserDeviceLastIP failed: %v", err)
	}
	t.Logf("Found IP: %s", ip)
}

func TestGetUsersDevicesLastIPWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	usersDevices := map[types.UserID][]types.DeviceInfo{
		"6895fe2e5a62f6f14af9d954": {
			{DeviceID: "00000000-0000-0000-0000-000000000000-9adbfe72-ca4e-4b7f-9c3b-262685c26741"},
		},
	}

	ipsMap := GetUsersDevicesLastIPWithProgress(&appState, usersDevices)
	t.Logf("Fetched IPs for %d. Map size: %d", len(usersDevices), len(ipsMap))
}
