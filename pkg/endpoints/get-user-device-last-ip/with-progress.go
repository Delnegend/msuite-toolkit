package get_user_device_last_ip

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

func GetUsersDevicesLastIPWithProgress(
	appState *types.AppState,
	usersDevices map[types.UserID][]types.DeviceInfo,
) map[types.UserID]map[types.DeviceID]types.IPAddress {
	fmt.Println("Fetching devices last IPs...")

	result := make(map[types.UserID]map[types.DeviceID]types.IPAddress)
	var mu sync.Mutex

	// start progress printer
	progressPercentChan := make(chan int)
	donePrinter := make(chan struct{})
	go func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
		close(donePrinter)
	}()

	// count total device tasks
	totalTasks := 0
	for _, devices := range usersDevices {
		totalTasks += len(devices)
	}

	if totalTasks == 0 {
		select {
		case progressPercentChan <- 100:
		default:
		}
		close(progressPercentChan)
		<-donePrinter
		return result
	}

	pool := pond.NewPool(appState.WorkerCount)

	var completed int32
	tasks := make([]pond.Task, 0, totalTasks)

	for uid, devices := range usersDevices {
		for _, d := range devices {
			userID := uid
			deviceID := d.DeviceID
			task := pool.SubmitErr(func() error {
				defer func() {
					atomic.AddInt32(&completed, 1)
					if progressPercentChan != nil {
						percent := int(atomic.LoadInt32(&completed)) * 100 / totalTasks
						if percent < 100 {
							select {
							case progressPercentChan <- percent:
							default:
							}
						}
					}
				}()

				ip, err := GetUserDeviceLastIP(appState, userID, deviceID)
				if err != nil {
					return fmt.Errorf("user %s device %s: %w", userID, deviceID, err)
				}
				mu.Lock()
				if _, ok := result[userID]; !ok {
					result[userID] = make(map[types.DeviceID]types.IPAddress)
				}
				result[userID][deviceID] = types.IPAddress(ip)
				mu.Unlock()
				return nil
			})
			tasks = append(tasks, task)
		}
	}

	// wait for all submitted jobs to finish
	pool.StopAndWait()

	// collect task errors
	var errs []error
	for _, t := range tasks {
		if tErr := t.Wait(); tErr != nil {
			errs = append(errs, tErr)
		}
	}

	// report final progress and close the progress channel
	select {
	case progressPercentChan <- 100:
	default:
	}
	close(progressPercentChan)
	<-donePrinter

	if len(errs) > 0 {
		for _, e := range errs {
			slog.Error("failed to get device last IP", "err", e)
		}
	}

	return result
}
