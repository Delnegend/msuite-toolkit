package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_user_apps "msuite-toolkit/pkg/endpoints/get-user-apps"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"
	"os"
	"strings"
)

func main() {
	outputPath := app.Init("apps_to_users.csv")

	as := &app.AppState

	users := get_users.GetUsersWithProgress(
		as,
		types.
			NewGetUsersRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
	)

	appsMap := get_user_apps.GetUserAppsWithProgress(as, users)

	// ONE APP to MANY USERS

	csvFile, err := os.Create(fmt.Sprintf("ONE_APP-to-MANY_USERS_%s", *outputPath))
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

	if err := w.Write([]string{"App", "Users"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for app, users := range appsMap {
		if err := w.Write(
			[]string{
				fmt.Sprintf("%s (%d)", app.App.Name, app.App.AppID),
				strings.Join(users, ","),
			}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}

	// ONE APP to ONE USER

	csvFile2, err := os.Create(fmt.Sprintf("ONE_APP-to-ONE_USER_%s", *outputPath))
	if err != nil {
		slog.Error("creating csv file failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := csvFile2.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w2 := csv.NewWriter(csvFile2)
	w2.Comma = '|'
	defer w2.Flush()

	if err := w2.Write([]string{"App", "User"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for app, users := range appsMap {
		for _, user := range users {
			if err := w2.Write(
				[]string{
					fmt.Sprintf("%s (%d)",
						app.App.Name, app.App.AppID),
					user,
				}); err != nil {
				slog.Error("writing csv row failed", "err", err)
				os.Exit(1)
			}
		}
	}

	if err := w2.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}

	// ONE USER to MANY APPS

	csvFile3, err := os.Create(fmt.Sprintf("ONE_USER-to-MANY_APPS_%s", *outputPath))
	if err != nil {
		slog.Error("creating csv file failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := csvFile3.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w3 := csv.NewWriter(csvFile3)
	w3.Comma = '|'
	defer w3.Flush()

	if err := w3.Write([]string{"User", "Apps"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	userToAppsMap := make(map[types.UserEmail][]string)
	for app, users := range appsMap {
		for _, user := range users {
			userToAppsMap[user] = append(userToAppsMap[user], fmt.Sprintf("%s (%d)", app.App.Name, app.App.AppID))
		}
	}

	for user, apps := range userToAppsMap {
		if err := w3.Write(
			[]string{
				user,
				strings.Join(apps, ","),
			}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w3.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}
