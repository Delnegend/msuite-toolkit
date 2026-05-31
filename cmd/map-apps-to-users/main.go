package main

import (
	"encoding/csv"
	"encoding/json"
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
			NewQueryRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
	)

	appsMap := get_user_apps.GetUserAppsWithProgress(as, users)
	filteredAppsMap := filterAppsByDestination(appsMap, as.FilterBy.DestinationHost, as.FilterBy.DestinationPort)

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

	if err := w.Write([]string{"App Name (ID)", "Destination Host", "Destination Port", "User Emails"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for app, users := range filteredAppsMap {

		if err := w.Write(
			[]string{
				fmt.Sprintf("%s (%d)", app.App.Name, app.App.AppID),
				app.App.DestinationSetting.IPDef,
				app.App.DestinationSetting.PortDef,
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

	// ONE USER to ONE APP

	csvFile2, err := os.Create(fmt.Sprintf("ONE_USER-to-ONE_APP_%s", *outputPath))
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

	if err := w2.Write([]string{"User Email", "App Name (ID)", "Destination Host", "Destination Port"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for app, users := range filteredAppsMap {

		for _, userEmail := range users {
			if err := w2.Write(
				[]string{
					userEmail,
					fmt.Sprintf("%s (%d)",
						app.App.Name, app.App.AppID),
					app.App.DestinationSetting.IPDef,
					app.App.DestinationSetting.PortDef,
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

	if err := w3.Write([]string{"User", "Apps", "Apps JSON"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	userToAppsMap := make(map[types.UserEmail][]*get_user_apps.AuthorizedApp)
	for app, users := range filteredAppsMap {
		for _, user := range users {
			userToAppsMap[user] = append(userToAppsMap[user], app)
		}
	}

	for user, apps := range userToAppsMap {
		if err := w3.Write(
			[]string{
				user,
				func() string {
					var appStrings []string
					for _, app := range apps {
						appStrings = append(appStrings, fmt.Sprintf("%s (%d)", app.App.Name, app.App.AppID))
					}
					return strings.Join(appStrings, ",")
				}(),
				func() string {
					// jsonData, err := json.Marshal(apps)
					var simplifiedApps []*get_user_apps.SimplifiedAppInfo
					for _, app := range apps {
						simplifiedAppInfo := app.ToSimplifiedAppInfo()
						simplifiedApps = append(simplifiedApps, simplifiedAppInfo)
					}
					jsonData, err := json.Marshal(simplifiedApps)
					if err != nil {
						slog.Error("marshaling apps to JSON failed", "err", err)
						return "[]"
					}
					return string(jsonData)
				}(),
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
