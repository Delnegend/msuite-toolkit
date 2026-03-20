package main

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_user_device_last_ip "msuite-toolkit/pkg/endpoints/get-user-device-last-ip"
	get_user_devices "msuite-toolkit/pkg/endpoints/get-user-devices"
	get_user_failed_logins "msuite-toolkit/pkg/endpoints/get-user-failed-logins"
	get_user_mfa "msuite-toolkit/pkg/endpoints/get-user-mfa"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"
	"os"
)

func main() {
	outputPath := app.Init("users_history.csv")

	users := get_users.GetUsersWithProgress(&app.AppState, types.NewGetUsersRequestBuilder().Build())

	userMFA := get_user_mfa.GetUsersMFAWithProgress(&app.AppState, users)

	userDevices := get_user_devices.GetUserDevicesWithProgress(&app.AppState, users)

	usersDevicesLastIP := get_user_device_last_ip.GetUsersDevicesLastIPWithProgress(&app.AppState, userDevices)

	userFailedLogins := get_user_failed_logins.GetUsersFailedLoginsWithProgress(&app.AppState, users)

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
	if err := w.Write([]string{"UserID", "UserEmail", "MFA", "Device", "FailedLogins"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	type deviceWithIP struct {
		types.DeviceInfo
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

		failedLoginsB := []byte("[]")
		if failedLogins, ok := userFailedLogins[user.UserID]; ok && failedLogins != nil {
			var err error
			failedLoginsB, err = json.Marshal(failedLogins)
			if err != nil {
				slog.Error("marshalling failed logins failed", "err", err, "userID", user.UserID)
				os.Exit(1)
			}
		}

		if err := w.Write([]string{
			user.UserID,
			user.UserEmail,
			string(mfaB),
			string(devB),
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
