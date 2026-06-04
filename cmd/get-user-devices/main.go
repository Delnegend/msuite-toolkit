package main

import (
	"encoding/csv"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_user_devices "msuite-toolkit/pkg/endpoints/get-user-devices"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"
	"os"
)

func main() {
	outputPath := app.Init("user_devices.csv")
	as := &app.AppState

	users := fetchUsersWithProgress(as)
	userDevices := fetchUserDevices(as, users)
	writeUserDevicesCSV(outputPath, users, userDevices)
}

// fetchUsersWithProgress retrieves all users for the configured OU.
func fetchUsersWithProgress(as *types.AppState) []types.UserInfo {
	return get_users.GetUsersWithProgress(
		as,
		types.NewQueryRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).
			Build(),
	)
}

// fetchUserDevices retrieves devices for all users with progress.
func fetchUserDevices(as *types.AppState, users []types.UserInfo) map[string][]types.DeviceInfo {
	return get_user_devices.GetUserDevicesWithProgress(as, users)
}

// writeUserDevicesCSV writes the user devices to a pipe-delimited CSV.
func writeUserDevicesCSV(outputPath *string, users []types.UserInfo, userDevices map[string][]types.DeviceInfo) {
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

	if err := w.Write([]string{"UserID", "UserEmail", "DeviceID", "DeviceName", "DeviceType", "LastUsed"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, user := range users {
		devices, ok := userDevices[user.UserID]
		if !ok {
			continue
		}
		for _, device := range devices {
			if err := w.Write([]string{
				user.UserID,
				user.UserEmail,
				device.DeviceID,
				device.DeviceName,
				device.MetaData.ProductName,
				device.UpdatedTimeString(),
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
