package blocks

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/endpoints"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

func GetUserDevicesWithProgress(appState *types.AppState, users []endpoints.UserInfo) map[types.UserID][]endpoints.DeviceInfo {
	fmt.Println("Fetching user devices...")

	userDeviceMap := make(map[types.UserID][]endpoints.DeviceInfo)
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

	pool := pond.NewPool(100)

	totalUsers := len(users)
	if totalUsers == 0 {
		// nothing to do; ensure progress shows 100% and clean up
		select {
		case progressPercentChan <- 100:
		default:
		}
		close(progressPercentChan)
		<-donePrinter

		return userDeviceMap
	}
	var completed int32
	tasks := make([]pond.Task, 0, totalUsers)

	for _, user := range users {
		task := pool.SubmitErr(func() error {
			// ensure progress is accounted for even on error
			defer func() {
				atomic.AddInt32(&completed, 1)
				if progressPercentChan != nil {
					percent := int(atomic.LoadInt32(&completed)) * 100 / totalUsers
					// do not send the final 100% from workers to avoid duplicate final prints
					if percent < 100 {
						select {
						case progressPercentChan <- percent:
						default:
						}
					}
				}
			}()

			devices, err := endpoints.GetUserDevices(appState, user.UserID)
			if err != nil {
				return fmt.Errorf("user %s: %w", user.UserID, err)
			}
			mu.Lock()
			userDeviceMap[user.UserID] = devices
			mu.Unlock()
			return nil
		})
		tasks = append(tasks, task)
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

	// report final progress (non-blocking) and close the progress channel
	select {
	case progressPercentChan <- 100:
	default:
	}
	close(progressPercentChan)
	<-donePrinter

	// print accumulated errors after progress printing finishes
	if len(errs) > 0 {
		for _, e := range errs {
			slog.Error("failed to get devices for a user", "err", e)
		}
	}

	return userDeviceMap
}
