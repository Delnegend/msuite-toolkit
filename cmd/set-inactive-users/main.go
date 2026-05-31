package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	inactive_user "msuite-toolkit/pkg/endpoints/inactive-user"
	"msuite-toolkit/pkg/types"
	"os"
	"strings"
	"time"
)

func main() {
	outputPath := app.Init("inactive_users.csv")

	users, err := get_users.GetAllUsers(
		&app.AppState,
		types.
			NewQueryRequestBuilder().
			WithFilterByOrgUnitID(app.AppState.OrganizationalUnitID).
			Build(),
		nil,
	)
	if err != nil {
		slog.Error("fetching users failed", "err", err)
		os.Exit(1)
	}

	now := time.Now()
	inactiveUsers := selectInactiveUsers(users, app.AppState.LastLoginThresholdInMonths, app.AppState.IncludeUsersWithUnknownLastLogin, now)
	totalUsers := len(users)
	selectedUsers := len(inactiveUsers)

	if err := writeResultsCSV(*outputPath, buildDryRunRows(inactiveUsers)); err != nil {
		slog.Error("writing csv file failed", "err", err)
		os.Exit(1)
	}

	if app.AppState.DryRun {
		slog.Info("dry run completed", "inactive", selectedUsers, "total", totalUsers)
		return
	}

	if app.AppState.LastLoginThresholdInMonths < 3 {
		if err := confirmDangerousThreshold(app.AppState.LastLoginThresholdInMonths); err != nil {
			slog.Error("confirmation failed", "err", err)
			os.Exit(1)
		}
	}

	if err := confirmWouldBeInactiveUsersRead(); err != nil {
		slog.Error("confirmation failed", "err", err)
		os.Exit(1)
	}

	slog.Info("users selected for inactivation", "inactive", selectedUsers, "total", totalUsers)

	var failed int
	rows := make([][]string, 0, selectedUsers)
	for _, user := range inactiveUsers {
		if err := inactive_user.InactiveUser(&app.AppState, types.UserID(user.UserID)); err != nil {
			msg := err.Error()
			slog.Error("failed to inactive user", "user_id", user.UserID, "err", err)
			failed++
			rows = append(rows, []string{user.UserID, user.UserEmail, lastLoginString(user), msg})
			continue
		}
		slog.Info("inactivated user", "user_id", user.UserID, "display_name", user.DisplayName, "email", user.UserEmail)
		rows = append(rows, []string{user.UserID, user.UserEmail, lastLoginString(user), "OK"})
	}

	if err := writeResultsCSV(*outputPath, rows); err != nil {
		slog.Error("writing csv file failed", "err", err)
		os.Exit(1)
	}

	if failed > 0 {
		slog.Error("completed with failures", "inactive", selectedUsers, "failed", failed, "total", totalUsers)
		os.Exit(1)
	}
	slog.Info("completed", "inactive", selectedUsers, "total", totalUsers)
}

func selectInactiveUsers(users []types.UserInfo, thresholdMonths int, includeUnknownLastLogin bool, now time.Time) []types.UserInfo {
	cutoff := now.AddDate(0, -thresholdMonths, 0)
	var inactive []types.UserInfo
	for _, user := range users {
		lastLogin := user.LastLoginTime()
		if lastLogin == nil {
			if includeUnknownLastLogin {
				inactive = append(inactive, user)
			}
			continue
		}
		if lastLogin.Before(cutoff) {
			inactive = append(inactive, user)
		}
	}
	return inactive
}

func lastLoginString(user types.UserInfo) string {
	if lastLogin := user.LastLoginTime(); lastLogin != nil {
		return lastLogin.UTC().Format(time.RFC3339)
	}
	return ""
}

func buildDryRunRows(users []types.UserInfo) [][]string {
	rows := make([][]string, 0, len(users))
	for _, user := range users {
		rows = append(rows, []string{user.UserID, user.UserEmail, lastLoginString(user), "DRY_RUN"})
	}
	return rows
}

func confirmDangerousThreshold(thresholdMonths int) error {
	confirmation := fmt.Sprintf("I WANT TO INACTIVE USERS WITH LAST LOGIN OLDER THAN %d MONTHS", thresholdMonths)
	fmt.Printf("Threshold is %d months. Type %q to continue: ", thresholdMonths, confirmation)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil && len(strings.TrimSpace(input)) == 0 {
		return fmt.Errorf("reading confirmation failed: %w", err)
	}
	if strings.TrimSpace(input) != confirmation {
		return fmt.Errorf("confirmation text did not match")
	}
	return nil
}

func confirmWouldBeInactiveUsersRead() error {
	confirmation := "I HAVE READ THE WOULD-BE-INACTIVE USERS LIST AND WANT TO PROCEED"
	fmt.Printf("Type %q to confirm you have read the would-be-inactive users list: ", confirmation)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil && len(strings.TrimSpace(input)) == 0 {
		return fmt.Errorf("reading confirmation failed: %w", err)
	}
	if strings.TrimSpace(input) != confirmation {
		return fmt.Errorf("confirmation text did not match")
	}
	return nil
}

func writeResultsCSV(path string, rows [][]string) error {
	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w := csv.NewWriter(csvFile)
	w.Comma = '|'
	defer w.Flush()

	if err := w.Write([]string{"UserID", "Email", "LastLoginTime", "Result"}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := w.Write(r); err != nil {
			return err
		}
	}

	return w.Error()
}
