package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"
	"os"
	"time"
)

func main() {
	outputPath := app.Init("active_users.csv")

	as := &app.AppState

	users, err := get_users.GetAllUsers(
		as,
		types.
			NewGetUsersRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
		nil,
	)
	if err != nil {
		slog.Error("fetching users failed", "err", err)
		os.Exit(1)
	}

	var active []types.UserInfo

	// compute cutoff if configured
	var cutoff time.Time
	if as.LastLoginThresholdInMonths > 0 {
		cutoff = time.Now().AddDate(0, -as.LastLoginThresholdInMonths, 0)
	}

	for _, u := range users {
		// filter by last login threshold if configured
		if as.LastLoginThresholdInMonths > 0 {
			if lt := u.LastLoginTime(); lt == nil || lt.Before(cutoff) {
				continue
			}
		}

		active = append(active, u)
	}

	// write result to output CSV file
	outName := fmt.Sprintf("active_users_%s", *outputPath)
	f, err := os.Create(outName)
	if err != nil {
		slog.Error("creating output file failed", "err", err)
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// header
	if err := w.Write([]string{"Display name", "Email", "Last login time"}); err != nil {
		slog.Error("writing CSV header failed", "err", err)
		os.Exit(1)
	}

	for _, u := range active {
		var lastLoginStr string
		if lt := u.LastLoginTime(); lt != nil {
			lastLoginStr = lt.UTC().Format(time.RFC3339)
		} else {
			lastLoginStr = ""
		}

		if err := w.Write([]string{u.DisplayName, u.UserEmail, lastLoginStr}); err != nil {
			slog.Error("writing CSV row failed", "err", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Wrote %d active users to %s\n", len(active), outName)
}

// Note: organizational unit filtering is applied server-side via the request payload.
