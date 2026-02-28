package main

import (
	"encoding/csv"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_user_devices "msuite-toolkit/pkg/endpoints/get-user-devices"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"os"
)

func main() {
	outputPath := app.Init("user_devices.csv")

	users := get_users.GetUsersWithProgress(&app.AppState)

	userDevices := get_user_devices.GetUserDevicesWithProgress(&app.AppState, users)

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
	if err := w.Write([]string{"UserID", "UserEmail", "DeviceID", "DeviceName", "DeviceType", "LastUsed"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	// write rows
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
				device.ProductName,
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
