package get_user_mfa

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

func GetUsersMFAWithProgress(appState *types.AppState, users []types.UserInfo) map[types.UserID]UserMFAInfo {
	fmt.Println("Fetching users MFA info...")

	userMFA := make(map[types.UserID]UserMFAInfo)
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
		return userMFA
	}

	var completed int32
	tasks := make([]pond.Task, 0, totalUsers)

	for _, user := range users {
		u := user
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

			mfa, err := GetUserMFA(appState, u.UserID)
			if err != nil {
				return fmt.Errorf("user %s: %w", u.UserID, err)
			}
			mu.Lock()
			userMFA[u.UserID] = mfa
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
			slog.Error("failed to get MFA for a user", "err", e)
		}
	}

	return userMFA
}
