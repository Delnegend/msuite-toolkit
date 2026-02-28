package get_user_failed_logins

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

// GetUsersFailedLoginsWithProgress fetches all failed logins for multiple users with an overall progress bar.
func GetUsersFailedLoginsWithProgress(appState *types.AppState, users []types.UserInfo) map[types.UserID][]FailedLogin {
	fmt.Println("Fetching users failed logins...")

	userFailedLogins := make(map[types.UserID][]FailedLogin)
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
		select {
		case progressPercentChan <- 100:
		default:
		}
		close(progressPercentChan)
		<-donePrinter
		return userFailedLogins
	}

	var completed int32
	tasks := make([]pond.Task, 0, totalUsers)

	for _, user := range users {
		u := user
		task := pool.SubmitErr(func() error {
			defer func() {
				c := atomic.AddInt32(&completed, 1)
				percent := int(c) * 100 / totalUsers
				if percent < 100 {
					select {
					case progressPercentChan <- percent:
					default:
					}
				}
			}()

			var allFailedLogins []FailedLogin
			limit := 200
			offset := 0
			for {
				total, batch, err := GetUserFailedLogins(appState, u.UserID, offset, limit)
				if err != nil {
					return fmt.Errorf("user %s: %w", u.UserID, err)
				}
				allFailedLogins = append(allFailedLogins, batch...)
				if offset+len(batch) >= total || len(batch) == 0 {
					break
				}
				offset += limit
			}

			mu.Lock()
			userFailedLogins[u.UserID] = allFailedLogins
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
			slog.Error("failed to get failed logins for a user", "err", e)
		}
	}

	return userFailedLogins
}
