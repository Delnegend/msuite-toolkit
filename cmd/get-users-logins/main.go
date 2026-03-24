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

	users := get_users.GetUsersWithProgress(
		as,
		types.
			NewGetUsersRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
	)

	userMFA := get_user_mfa.GetUsersMFAWithProgress(as, users)

	userFailedLogins := get_user_failed_logins.GetUsersFailedLoginsWithProgress(as, users)

	// create CSV file
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

	// write header
	if err := w.Write([]string{"User ID", "Email", "MFA", "Failed logins"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	// write rows
	for _, user := range users {
		mfaB := []byte("{}")
		if mfa, ok := userMFA[user.UserID]; ok {
			var err error
			mfaB, err = json.Marshal(mfa)
			if err != nil {
				slog.Error("marshalling mfa failed", "err", err, "userID", user.UserID)
				os.Exit(1)
			}
		}

		failedLoginsB := []byte("[]")
		if fls, ok := userFailedLogins[user.UserID]; ok && len(fls) > 0 {
			var err error
			failedLoginsB, err = json.Marshal(fls)
			if err != nil {
				slog.Error("marshalling failed logins failed", "err", err, "userID", user.UserID)
				os.Exit(1)
			}
		}

		if err := w.Write([]string{
			user.UserID,
			user.UserEmail,
			string(mfaB),
			string(failedLoginsB),
		}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}
