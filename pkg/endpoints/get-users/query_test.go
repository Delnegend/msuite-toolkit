package get_users

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetUsers(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	// example: default call — use builder to get defaults
	count, users, err := GetUsers(&appState, types.NewGetUsersRequestBuilder().Build())
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
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	users := GetUsersWithProgress(
		&appState,
		types.
			NewGetUsersRequestBuilder().
			WithFilterByOrgUnitID(appState.OrganizationalUnitID).
			Build(),
	)
	t.Logf("Total users fetched with progress: %d", len(users))
}

func TestGetUsersWithCustomPayload(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	// only specify fields we want to change (limit and orders) using builder pattern
	count, users, err := GetUsers(&appState, types.NewGetUsersRequestBuilder().WithLimit(5).WithOrders(map[string]int{"created_time": -1}).Build())
	if err != nil {
		t.Fatalf("GetUsers with custom payload failed: %v", err)
	}

	t.Logf("Custom payload fetch: returned %d users (total %d)", len(users), count)
}
