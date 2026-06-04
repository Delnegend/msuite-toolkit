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

	users := fetchAllUsers(as)
	active := filterActiveUsers(users, as.LastLoginThresholdInMonths)
	writeActiveUsersCSV(outputPath, active)
	fmt.Printf("Wrote %d active users\n", len(active))
}

// fetchAllUsers retrieves all users for the configured OU.
func fetchAllUsers(as *types.AppState) []types.UserInfo {
	users, err := get_users.GetAllUsers(
		as,
		types.NewQueryRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
		nil,
	)
	if err != nil {
		slog.Error("fetching users failed", "err", err)
		os.Exit(1)
	}
	return users
}

// filterActiveUsers returns users whose last login is within the threshold.
func filterActiveUsers(users []types.UserInfo, thresholdMonths int) []types.UserInfo {
	var cutoff time.Time
	if thresholdMonths > 0 {
		cutoff = time.Now().AddDate(0, -thresholdMonths, 0)
	}

	var active []types.UserInfo
	for _, u := range users {
		if thresholdMonths > 0 {
			if lt := u.LastLoginTime(); lt == nil || lt.Before(cutoff) {
				continue
			}
		}
		active = append(active, u)
	}
	return active
}

// formatLastLogin returns the user's last login as an RFC3339 string, or empty.
func formatLastLogin(u types.UserInfo) string {
	if lt := u.LastLoginTime(); lt != nil {
		return lt.UTC().Format(time.RFC3339)
	}
	return ""
}

// writeActiveUsersCSV writes the active users list to a CSV file.
func writeActiveUsersCSV(outputPath *string, active []types.UserInfo) {
	outName := fmt.Sprintf("active_users_%s", *outputPath)
	csvFile, err := os.Create(outName)
	if err != nil {
		slog.Error("creating output file failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w := csv.NewWriter(csvFile)
	w.Comma = '|'
	defer w.Flush()

	if err := w.Write([]string{"Display name", "Email", "Last login time"}); err != nil {
		slog.Error("writing CSV header failed", "err", err)
		os.Exit(1)
	}

	for _, u := range active {
		if err := w.Write([]string{u.DisplayName, u.UserEmail, formatLastLogin(u)}); err != nil {
			slog.Error("writing CSV row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}
