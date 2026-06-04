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

type deviceWithIP struct {
	types.DeviceInfo
	IP types.IPAddress `json:"ip"`
}

func main() {
	outputPath := app.Init("users_history.csv")
	as := &app.AppState

	users := fetchHistoryUsers(as)
	userMFA := fetchHistoryMFA(as, users)
	userDevices := fetchHistoryDevices(as, users)
	usersDevicesLastIP := fetchDevicesLastIP(as, userDevices)
	userFailedLogins := fetchHistoryFailedLogins(as, users)

	writeUsersHistoryCSV(outputPath, users, userMFA, userDevices, usersDevicesLastIP, userFailedLogins)
}

// fetchHistoryUsers retrieves all users for the configured OU.
func fetchHistoryUsers(as *types.AppState) []types.UserInfo {
	return get_users.GetUsersWithProgress(
		as,
		types.NewQueryRequestBuilder().
			WithFilterByOrgUnitID(as.OrganizationalUnitID).Build(),
	)
}

// fetchHistoryMFA retrieves MFA info for all users with progress.
func fetchHistoryMFA(as *types.AppState, users []types.UserInfo) map[types.UserID]get_user_mfa.UserMFAInfo {
	return get_user_mfa.GetUsersMFAWithProgress(as, users)
}

// fetchHistoryDevices retrieves devices for all users with progress.
func fetchHistoryDevices(as *types.AppState, users []types.UserInfo) map[string][]types.DeviceInfo {
	return get_user_devices.GetUserDevicesWithProgress(as, users)
}

// fetchDevicesLastIP retrieves the last IP for each device with progress.
func fetchDevicesLastIP(as *types.AppState, userDevices map[string][]types.DeviceInfo) map[string]map[string]types.IPAddress {
	return get_user_device_last_ip.GetUsersDevicesLastIPWithProgress(as, userDevices)
}

// fetchHistoryFailedLogins retrieves failed login info for all users with progress.
func fetchHistoryFailedLogins(as *types.AppState, users []types.UserInfo) map[string][]get_user_failed_logins.FailedLogin {
	return get_user_failed_logins.GetUsersFailedLoginsWithProgress(as, users)
}

// marshalUserMFA marshals MFA data to JSON, returning "{}" on error.
func marshalUserMFA(mfa get_user_mfa.UserMFAInfo) string {
	b, err := json.Marshal(mfa)
	if err != nil {
		slog.Error("marshalling mfa failed", "err", err)
		return "{}"
	}
	return string(b)
}

// buildDevicesWithIP merges device info with last known IP for a user.
func buildDevicesWithIP(devices []types.DeviceInfo, ips map[string]types.IPAddress) []deviceWithIP {
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
	return devs
}

// marshalDevicesWithIP marshals devices with IPs to JSON, returning "[]" on error.
func marshalDevicesWithIP(devs []deviceWithIP) string {
	b, err := json.Marshal(devs)
	if err != nil {
		slog.Error("marshalling devices failed", "err", err)
		return "[]"
	}
	return string(b)
}

// marshalUserFailedLogins marshals failed logins to JSON, returning "[]" on error.
func marshalUserFailedLogins(fls []get_user_failed_logins.FailedLogin) string {
	if fls == nil {
		return "[]"
	}
	b, err := json.Marshal(fls)
	if err != nil {
		slog.Error("marshalling failed logins failed", "err", err)
		return "[]"
	}
	return string(b)
}

// writeUsersHistoryCSV writes the full user history report to a pipe-delimited CSV.
func writeUsersHistoryCSV(
	outputPath *string,
	users []types.UserInfo,
	userMFA map[types.UserID]get_user_mfa.UserMFAInfo,
	userDevices map[string][]types.DeviceInfo,
	usersDevicesLastIP map[string]map[string]types.IPAddress,
	userFailedLogins map[string][]get_user_failed_logins.FailedLogin,
) {
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

	if err := w.Write([]string{"UserID", "UserEmail", "MFA", "Device", "FailedLogins"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, user := range users {
		mfaB := marshalUserMFA(userMFA[user.UserID])
		devs := buildDevicesWithIP(userDevices[user.UserID], usersDevicesLastIP[user.UserID])
		devB := marshalDevicesWithIP(devs)
		failedLoginsB := marshalUserFailedLogins(userFailedLogins[user.UserID])

		if err := w.Write([]string{user.UserID, user.UserEmail, mfaB, devB, failedLoginsB}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}
