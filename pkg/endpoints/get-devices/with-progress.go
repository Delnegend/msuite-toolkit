package get_devices

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"os"
	"sync"
)

func GetDevicesWithProgress(appState *types.AppState, basePayload types.QueryRequestPayload) []types.DeviceInfo {
	fmt.Println("Fetching devices...")
	var wg sync.WaitGroup
	progressPercentChan := make(chan int)
	wg.Go(func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
	})
	devices, err := GetAllDevices(appState, basePayload, progressPercentChan)
	close(progressPercentChan)
	wg.Wait()
	if err != nil {
		slog.Error("failed to get all devices", "err", err)
		os.Exit(1)
	}
	slog.Info("fetched devices", "count", len(devices))
	return devices
}
