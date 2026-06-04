package main

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_user_failed_logins "msuite-toolkit/pkg/endpoints/get-user-failed-logins"
	get_user_mfa "msuite-toolkit/pkg/endpoints/get-user-mfa"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"
	"os"
)

func main() {
	outputPath := app.Init("users_logins.csv")
	as := &app.AppState

	users := fetchUsersFromOU(as)
	userMFA := fetchUsersMFA(as, users)
	userFailedLogins := fetchUsersFailedLogins(as, users)
	writeUsersLoginsCSV(outputPath, users, userMFA, userFailedLogins)
}

// fetchUsersFromOU retrieves all users for the configured OU.
func fetchUsersFromOU(as *types.AppState) []types.UserInfo {
	return get_users.GetUsersWithProgress(
		as,
		types.NewQueryRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
	)
}

// fetchUsersMFA retrieves MFA info for all users with progress.
func fetchUsersMFA(as *types.AppState, users []types.UserInfo) map[types.UserID]get_user_mfa.UserMFAInfo {
	return get_user_mfa.GetUsersMFAWithProgress(as, users)
}

// fetchUsersFailedLogins retrieves failed login info for all users with progress.
func fetchUsersFailedLogins(as *types.AppState, users []types.UserInfo) map[string][]get_user_failed_logins.FailedLogin {
	return get_user_failed_logins.GetUsersFailedLoginsWithProgress(as, users)
}

// marshalMFA marshals MFA data to JSON, returning "{}" on empty or error.
func marshalMFA(mfa any) string {
	if mfa == nil {
		return "{}"
	}
	b, err := json.Marshal(mfa)
	if err != nil {
		slog.Error("marshalling mfa failed", "err", err)
		return "{}"
	}
	return string(b)
}

// marshalFailedLogins marshals failed logins to JSON, returning "[]" on empty or error.
func marshalFailedLogins(fls []get_user_failed_logins.FailedLogin) string {
	if len(fls) == 0 {
		return "[]"
	}
	b, err := json.Marshal(fls)
	if err != nil {
		slog.Error("marshalling failed logins failed", "err", err)
		return "[]"
	}
	return string(b)
}

// writeUsersLoginsCSV writes user login data to a pipe-delimited CSV.
func writeUsersLoginsCSV(outputPath *string, users []types.UserInfo, userMFA map[types.UserID]get_user_mfa.UserMFAInfo, userFailedLogins map[string][]get_user_failed_logins.FailedLogin) {
	csvFile, err := os.Create(*outputPath)
	if err != nil {
		slog.Error("creating csv file failed", "err", err)
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

	if err := w.Write([]string{"User ID", "Email", "MFA", "Failed logins"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, user := range users {
		mfaB := marshalMFA(userMFA[user.UserID])
		failedLoginsB := marshalFailedLogins(userFailedLogins[user.UserID])

		if err := w.Write([]string{user.UserID, user.UserEmail, mfaB, failedLoginsB}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}
