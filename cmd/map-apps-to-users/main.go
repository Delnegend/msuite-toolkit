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

	users := fetchUsers(as)
	users = filterUsersByIncludeEmails(users, as.IncludeEmails)
	appsMap := fetchUserApps(as, users)
	filteredAppsMap := filterAppsByDestination(appsMap, as.FilterBy.DestinationHost, as.FilterBy.DestinationPort)

	writeOneAppToManyUsersCSV(outputPath, filteredAppsMap)
	writeOneUserToOneAppCSV(outputPath, filteredAppsMap)
	writeOneUserToManyAppsCSV(outputPath, filteredAppsMap)
}

// fetchUsers retrieves all users for the configured OU.
func fetchUsers(as *types.AppState) []types.UserInfo {
	return get_users.GetUsersWithProgress(
		as,
		types.
			NewQueryRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
	)
}

// fetchUserApps retrieves the authorized apps for each user.
func fetchUserApps(as *types.AppState, users []types.UserInfo) map[*get_user_apps.AuthorizedApp][]string {
	return get_user_apps.GetUserAppsWithProgress(as, users)
}

// filterUsersByIncludeEmails narrows the user list to only those whose email
// appears in includeEmails. When includeEmails is empty, all users pass through.
func filterUsersByIncludeEmails(users []types.UserInfo, includeEmails []string) []types.UserInfo {
	if len(includeEmails) == 0 {
		return users
	}
	includeSet := make(map[string]struct{}, len(includeEmails))
	for _, e := range includeEmails {
		includeSet[e] = struct{}{}
	}
	filtered := make([]types.UserInfo, 0, len(users))
	for _, u := range users {
		if _, ok := includeSet[u.UserEmail]; ok {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

// writeOneAppToManyUsersCSV writes a CSV mapping each app to the list of users assigned to it.
func writeOneAppToManyUsersCSV(outputPath *string, appsMap map[*get_user_apps.AuthorizedApp][]string) {
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

	for app, users := range appsMap {
		if err := w.Write([]string{
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
}

// writeOneUserToOneAppCSV writes a CSV with one row per user-app pair.
func writeOneUserToOneAppCSV(outputPath *string, appsMap map[*get_user_apps.AuthorizedApp][]string) {
	csvFile, err := os.Create(fmt.Sprintf("ONE_USER-to-ONE_APP_%s", *outputPath))
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

	if err := w.Write([]string{"User Email", "App Name (ID)", "Destination Host", "Destination Port"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for app, users := range appsMap {
		for _, userEmail := range users {
			if err := w.Write([]string{
				userEmail,
				fmt.Sprintf("%s (%d)", app.App.Name, app.App.AppID),
				app.App.DestinationSetting.IPDef,
				app.App.DestinationSetting.PortDef,
			}); err != nil {
				slog.Error("writing csv row failed", "err", err)
				os.Exit(1)
			}
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}

// buildUserToAppsMap inverts the app→users map into a user→apps map.
func buildUserToAppsMap(appsMap map[*get_user_apps.AuthorizedApp][]string) map[types.UserEmail][]*get_user_apps.AuthorizedApp {
	userToAppsMap := make(map[types.UserEmail][]*get_user_apps.AuthorizedApp)
	for app, users := range appsMap {
		for _, user := range users {
			userToAppsMap[user] = append(userToAppsMap[user], app)
		}
	}
	return userToAppsMap
}

// writeOneUserToManyAppsCSV writes a CSV mapping each user to all their apps.
func writeOneUserToManyAppsCSV(outputPath *string, appsMap map[*get_user_apps.AuthorizedApp][]string) {
	csvFile, err := os.Create(fmt.Sprintf("ONE_USER-to-MANY_APPS_%s", *outputPath))
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

	if err := w.Write([]string{"User", "Apps", "Apps JSON"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	userToAppsMap := buildUserToAppsMap(appsMap)
	for user, apps := range userToAppsMap {
		appNames := formatAppNames(apps)
		appsJSON := marshalAppsJSON(apps)

		if err := w.Write([]string{user, appNames, appsJSON}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}

// formatAppNames returns a comma-separated string of "Name (ID)" for each app.
func formatAppNames(apps []*get_user_apps.AuthorizedApp) string {
	var appStrings []string
	for _, app := range apps {
		appStrings = append(appStrings, fmt.Sprintf("%s (%d)", app.App.Name, app.App.AppID))
	}
	return strings.Join(appStrings, ",")
}

// marshalAppsJSON serializes the apps to a JSON array of simplified app info.
func marshalAppsJSON(apps []*get_user_apps.AuthorizedApp) string {
	var simplifiedApps []*get_user_apps.SimplifiedAppInfo
	for _, app := range apps {
		simplifiedApps = append(simplifiedApps, app.ToSimplifiedAppInfo())
	}
	jsonData, err := json.Marshal(simplifiedApps)
	if err != nil {
		slog.Error("marshaling apps to JSON failed", "err", err)
		return "[]"
	}
	return string(jsonData)
}
