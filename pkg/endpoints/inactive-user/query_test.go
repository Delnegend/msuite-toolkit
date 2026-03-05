package inactive_user

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestInactiveUser(t *testing.T) {
	const TEST_USER_ID = "6895fe2e5a62f6f14af9d954"

	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	err := InactiveUser(&appState, TEST_USER_ID)
	if err != nil {
		t.Fatalf("InactiveUser failed: %v", err)
	}
	t.Logf("Successfully (in)activated user %s", TEST_USER_ID)
}
