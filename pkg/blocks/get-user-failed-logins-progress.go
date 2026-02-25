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

// GetUserFailedLoginsWithProgress fetches all failed logins for a specific user by making paginated requests.
func GetUserFailedLoginsWithProgress(as *types.AppState, userID string) ([]endpoints.FailedLogin, error) {
	fmt.Printf("Fetching failed logins for user %s...\n", userID)

	// start progress printer
	progressPercentChan := make(chan int)
	donePrinter := make(chan struct{})
	go func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
		close(donePrinter)
	}()

	pool := pond.NewPool(16)
	limit := 200

	// first request to get total count and first batch
	total, firstBatch, err := endpoints.GetUserFailedLogins(as, userID, 0, limit)
	if err != nil {
		slog.Error("initial fetch of failed logins failed", "err", err)
		close(progressPercentChan)
		<-donePrinter
		return nil, err
	}

	if total == 0 {
		select {
		case progressPercentChan <- 100:
		default:
		}
		close(progressPercentChan)
		<-donePrinter
		return nil, nil
	}

	failedLogins := make([]endpoints.FailedLogin, 0, total)
	var mu sync.Mutex
	failedLogins = append(failedLogins, firstBatch...)

	// calculate how many pages total (we already fetched page 0)
	pages := (total + limit - 1) / limit

	var completed int32 = 1 // first page already done
	// report initial progress
	select {
	case progressPercentChan <- int(completed * 100 / int32(pages)):
	default:
	}

	// if there's nothing else to do, report 100% and return
	if pages == 1 {
		select {
		case progressPercentChan <- 100:
		default:
		}
		close(progressPercentChan)
		<-donePrinter
		return failedLogins, nil
	}

	tasks := make([]pond.Task, 0, pages-1)
	// submit jobs for remaining pages
	for page := 1; page < pages; page++ {
		offset := page * limit
		task := pool.SubmitErr(func() error {
			defer func() {
				c := atomic.AddInt32(&completed, 1)
				percent := int(c * 100 / int32(pages))
				if percent < 100 {
					select {
					case progressPercentChan <- percent:
					default:
					}
				}
			}()

			_, batch, err := endpoints.GetUserFailedLogins(as, userID, offset, limit)
			if err != nil {
				slog.Error("fetching failed logins failed", "err", err, "offset", offset)
				return err
			}
			if len(batch) != 0 {
				mu.Lock()
				failedLogins = append(failedLogins, batch...)
				mu.Unlock()
			}
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

	// report final progress and close the progress channel
	select {
	case progressPercentChan <- 100:
	default:
	}
	close(progressPercentChan)
	<-donePrinter

	if len(errs) > 0 {
		return failedLogins, fmt.Errorf("encountered %d errors during fetch", len(errs))
	}
	return failedLogins, nil
}
