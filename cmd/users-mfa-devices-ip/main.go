package main

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"msuite-toolkit/pkg/app"
	"msuite-toolkit/pkg/blocks"
	"msuite-toolkit/pkg/endpoints"
	"msuite-toolkit/pkg/types"
	"os"
)

func main() {
	outputPath := app.Init("users_mfa_devices_ip.csv")

	users := blocks.GetUsersWithProgress(&app.AppState)

	userMFA := blocks.GetUsersMFAWithProgress(&app.AppState, users)

	userDevices := blocks.GetUserDevicesWithProgress(&app.AppState, users)

	usersDevicesLastIP := blocks.GetUsersDevicesLastIPWithProgress(&app.AppState, userDevices)

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
	if err := w.Write([]string{"UserID", "UserEmail", "MFA", "Device"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	type deviceWithIP struct {
		endpoints.DeviceInfo
		IP types.IPAddress `json:"ip"`
	}

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

		devices := userDevices[user.UserID]
		ips := usersDevicesLastIP[user.UserID]
		devs := make([]deviceWithIP, 0, len(devices))
		for _, d := range devices {
			ip := types.IPAddress("")
			if ips != nil {
				if v, ok := ips[d.DeviceID]; ok {
					ip = v
				}
			}
			devs = append(devs, deviceWithIP{DeviceInfo: d, IP: ip})
		}

		devB, err := json.Marshal(devs)
		if err != nil {
			slog.Error("marshalling devices failed", "err", err)
			os.Exit(1)
		}

		if err := w.Write([]string{
			user.UserID,
			user.UserEmail,
			string(mfaB),
			string(devB),
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
