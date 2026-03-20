package get_user_devices

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUserDevices(t *testing.T) {
	const TEST_USER_ID = "6895fe2e5a62f6f14af9d954"

	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	devices, err := GetUserDevices(&appState, TEST_USER_ID)
	if err != nil {
		t.Fatalf("GetUserDevices failed: %v", err)
	}
	t.Logf("User devices count: %d", len(devices))
	for i, device := range devices {
		if i < 5 {
			t.Logf("Device %d: ID=%s, Name=%s, OS=%s", i, device.DeviceID, device.DeviceName, device.OS)
		}
	}
}

func TestGetUserDevicesWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	users := []types.UserInfo{
		{UserID: "6895fe2e5a62f6f14af9d954"},
	}

	devicesMap := GetUserDevicesWithProgress(&appState, users)
	t.Logf("Fetched devices for %d users. Map size: %d", len(users), len(devicesMap))
}
