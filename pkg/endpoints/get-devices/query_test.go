package get_devices

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetDevices(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	count, devices, err := GetDevices(&appState, types.NewQueryRequestBuilder().Build())
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	t.Logf("Devices count in batch: %d (total: %d)", len(devices), count)
	for i, device := range devices {
		if i >= 5 {
			break
		}
		t.Logf("Device %d: ID=%s, Name=%s, OS=%s", i, device.DeviceID, device.DeviceName, device.MetaData.OS)
	}
}

func TestGetDevicesWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	devices := GetDevicesWithProgress(&appState, types.NewQueryRequestBuilder().Build())
	t.Logf("Total devices fetched with progress: %d", len(devices))
}
