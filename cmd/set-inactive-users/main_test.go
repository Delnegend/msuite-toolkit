package main

import (
	"strconv"
	"testing"
	"time"

	"msuite-toolkit/pkg/types"
)

func TestSelectInactiveUsers(t *testing.T) {
	currentTime := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	oldLogin := strconv.FormatInt(currentTime.AddDate(0, -4, 0).Unix(), 10)
	recentLogin := strconv.FormatInt(currentTime.AddDate(0, -1, 0).Unix(), 10)

	users := []types.UserInfo{
		{UserID: "1", Attributes: map[string]string{"last_login_time": oldLogin}},
		{UserID: "2", Attributes: map[string]string{"last_login_time": recentLogin}},
		{UserID: "3", Attributes: map[string]string{}},
	}

	inactive := selectInactiveUsers(users, 3, false, currentTime)
	if len(inactive) != 1 {
		t.Fatalf("unexpected inactive count: got %d want 1", len(inactive))
	}

	if inactive[0].UserID != "1" {
		t.Fatalf("unexpected inactive users: got %#v", inactive)
	}

	withUnknown := selectInactiveUsers(users, 3, true, currentTime)
	if len(withUnknown) != 2 {
		t.Fatalf("unexpected inactive count with unknown included: got %d want 2", len(withUnknown))
	}

	if withUnknown[0].UserID != "1" || withUnknown[1].UserID != "3" {
		t.Fatalf("unexpected inactive users with unknown included: got %#v", withUnknown)
	}
}
