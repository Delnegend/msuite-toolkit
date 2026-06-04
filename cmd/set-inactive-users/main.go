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
	as := &app.AppState

	users := fetchInactiveCandidates(as)
	inactiveUsers := selectInactiveUsers(users, as.LastLoginThresholdInMonths, as.IncludeUsersWithUnknownLastLogin, time.Now())

	writeResultsCSV(*outputPath, buildDryRunRows(inactiveUsers))

	if as.DryRun {
		slog.Info("dry run completed", "inactive", len(inactiveUsers), "total", len(users))
		return
	}

	confirmProceed(as)
	results := inactivateUsers(inactiveUsers)
	writeResultsCSV(*outputPath, results)
	reportInactivationResults(results, len(inactiveUsers), len(users))
}

// fetchInactiveCandidates retrieves all users for the configured OU.
func fetchInactiveCandidates(as *types.AppState) []types.UserInfo {
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

// confirmProceed runs interactive confirmations before inactivating users.
func confirmProceed(as *types.AppState) {
	if as.LastLoginThresholdInMonths < 3 {
		if err := confirmDangerousThreshold(as.LastLoginThresholdInMonths); err != nil {
			slog.Error("confirmation failed", "err", err)
			os.Exit(1)
		}
	}

	if err := confirmWouldBeInactiveUsersRead(); err != nil {
		slog.Error("confirmation failed", "err", err)
		os.Exit(1)
	}
}

// inactivateUsers performs the inactivation for each user and returns result rows.
func inactivateUsers(users []types.UserInfo) [][]string {
	rows := make([][]string, 0, len(users))
	for _, user := range users {
		result := "OK"
		if err := inactive_user.InactiveUser(&app.AppState, types.UserID(user.UserID)); err != nil {
			slog.Error("failed to inactive user", "user_id", user.UserID, "err", err)
			result = err.Error()
		} else {
			slog.Info("inactivated user", "user_id", user.UserID, "display_name", user.DisplayName, "email", user.UserEmail)
		}
		rows = append(rows, []string{user.UserID, user.UserEmail, lastLoginString(user), result})
	}
	return rows
}

// reportInactivationResults logs a summary of the inactivation run.
func reportInactivationResults(results [][]string, selectedUsers, totalUsers int) {
	var failed int
	for _, r := range results {
		if r[3] != "OK" {
			failed++
		}
	}
	if failed > 0 {
		slog.Error("completed with failures", "inactive", selectedUsers, "failed", failed, "total", totalUsers)
		os.Exit(1)
	}
	slog.Info("completed", "inactive", selectedUsers, "total", totalUsers)
}

// selectInactiveUsers returns users whose last login is older than thresholdMonths
// (or users with unknown last login if includeUnknownLastLogin is true).
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

// lastLoginString returns the user's last login time as an RFC3339 string, or empty.
func lastLoginString(user types.UserInfo) string {
	if lastLogin := user.LastLoginTime(); lastLogin != nil {
		return lastLogin.UTC().Format(time.RFC3339)
	}
	return ""
}

// buildDryRunRows constructs CSV rows annotated with "DRY_RUN" for preview purposes.
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
