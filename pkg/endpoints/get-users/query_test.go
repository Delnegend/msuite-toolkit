package get_users

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUsers(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	count, users, err := GetUsers(&appState, 0, 10)
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}

	t.Logf("Users count in batch: %d (total: %d)", len(users), count)
	for i, user := range users {
		t.Logf("User %d: ID=%s, Email=%s", i, user.UserID, user.UserEmail)
	}
}

func TestGetUsersWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	users := GetUsersWithProgress(&appState)
	t.Logf("Total users fetched with progress: %d", len(users))
}
