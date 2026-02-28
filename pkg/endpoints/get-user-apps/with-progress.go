package get_user_apps

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

// GetUserAppsWithProgress fetches authorized apps for each user and returns
// a map keyed by formatted app string "<name> (<id>)" to the slice of user IDs
// who have that app.
func GetUserAppsWithProgress(appState *types.AppState, users []types.UserInfo) map[string][]types.UserEmail {
	fmt.Println("Fetching user apps...")

	appsMap := make(map[string][]types.UserEmail)
	var mu sync.Mutex

	// progress printer
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
		select {
		case progressPercentChan <- 100:
		default:
		}
		close(progressPercentChan)
		<-donePrinter
		return appsMap
	}

	var completed int32
	tasks := make([]pond.Task, 0, totalUsers)

	for _, u := range users {
		task := pool.SubmitErr(func() error {
			defer func() {
				atomic.AddInt32(&completed, 1)
				if progressPercentChan != nil {
					percent := int(atomic.LoadInt32(&completed)) * 100 / totalUsers
					if percent < 100 {
						select {
						case progressPercentChan <- percent:
						default:
						}
					}
				}
			}()

			apps, err := GetUserApps(appState, u.UserID)
			if err != nil {
				return fmt.Errorf("user %s: %w", u.UserEmail, err)
			}

			mu.Lock()
			for _, a := range apps {
				key := fmt.Sprintf("%s (%d)", a.Name, a.AppID)
				appsMap[key] = append(appsMap[key], u.UserEmail)
			}
			mu.Unlock()
			return nil
		})
		tasks = append(tasks, task)
	}

	pool.StopAndWait()

	var errs []error
	for _, t := range tasks {
		if tErr := t.Wait(); tErr != nil {
			errs = append(errs, tErr)
		}
	}

	select {
	case progressPercentChan <- 100:
	default:
	}
	close(progressPercentChan)
	<-donePrinter

	if len(errs) > 0 {
		for _, e := range errs {
			slog.Error("failed to get apps for a user", "err", e)
		}
	}

	return appsMap
}
