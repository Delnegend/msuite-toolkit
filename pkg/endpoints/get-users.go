package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

type UserInfo struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"email"`
}

// GetUsers fetches a batch of users starting from the given offset with the specified limit.
// It returns the total count of users, the list of UserInfo, and any error encountered.
func GetUsers(as *types.AppState, offset int, limit int) (int, []UserInfo, error) {
	endpoint := fmt.Sprintf("https://%s/identity-api/v1/domains/default/users", as.AdminPortalAddress)

	reqPayloadBytes, err := json.Marshal(struct {
		Offset      int            `json:"offset"`
		Limit       int            `json:"limit"`
		Orders      map[string]int `json:"orders"`
		Search      string         `json:"search"`
		Filters     []any          `json:"filters"`
		ExtraParams map[string]any `json:"extra_params"`
	}{
		Offset:      offset,
		Limit:       limit,
		Orders:      map[string]int{"created_time": 1},
		Search:      "",
		Filters:     []any{},
		ExtraParams: map[string]any{},
	})
	if err != nil {
		slog.Error("marshalling request payload failed", "err", err)
		return 0, nil, fmt.Errorf("marshalling request payload failed: %w", err)
	}

	values := url.Values{}
	values.Set("ctx.user_id", as.AdminUserID)
	values.Set("request_payload", string(reqPayloadBytes))
	reqURL := endpoint + "?" + values.Encode()

	client := getHTTPClient()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
	if err != nil {
		slog.Error("creating request failed", "err", err)
		return 0, nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+as.BearerToken)

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "err", err)
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body failed", "err", err)
		return 0, nil, fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
		return 0, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respPayload struct {
		Data  []UserInfo `json:"data"`
		Count int        `json:"count"`
	}
	if err := json.Unmarshal(body, &respPayload); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return 0, nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	return respPayload.Count, respPayload.Data, nil
}

// GetAllUsers fetches all users by making paginated requests.
func GetAllUsers(as *types.AppState, progressPercentChan chan<- int) ([]UserInfo, error) {
	pool := pond.NewPool(16)

	limit := 100

	// send initial 0% (non-blocking)
	if progressPercentChan != nil {
		select {
		case progressPercentChan <- 0:
		default:
		}
	}

	// first request to get total count and first batch
	total, firstBatch, err := GetUsers(as, 0, limit)
	if err != nil {
		slog.Error("initial fetch of users failed", "err", err)
		return nil, err
	}

	users := make([]UserInfo, 0, total)
	mutex := &sync.Mutex{}
	users = append(users, firstBatch...)

	// calculate how many pages total (we already fetched page 0)
	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	// track per-task errors
	var errs []error
	tasks := make([]pond.Task, 0, max(0, pages-1))

	var completed int32 = 1 // first page already done
	// report progress after initial page (non-blocking)
	if progressPercentChan != nil {
		percent := int(atomic.LoadInt32(&completed)) * 100 / pages
		select {
		case progressPercentChan <- percent:
		default:
		}
	}

	// if there's nothing else to do, report 100% and return
	if pages == 1 {
		if progressPercentChan != nil {
			select {
			case progressPercentChan <- 100:
			default:
			}
		}
		return users, nil
	}

	// submit jobs for remaining pages using SubmitErr so each task can return an error
	for page := 1; page < pages; page++ {
		offset := page * limit
		task := pool.SubmitErr(func() error {
			_, batch, err := GetUsers(as, offset, limit)
			if err != nil {
				slog.Error("fetching users failed", "err", err, "offset", offset)
				atomic.AddInt32(&completed, 1)
				if progressPercentChan != nil {
					percent := int(atomic.LoadInt32(&completed)) * 100 / pages
					select {
					case progressPercentChan <- percent:
					default:
					}
				}
				return err
			}
			if len(batch) != 0 {
				mutex.Lock()
				users = append(users, batch...)
				mutex.Unlock()
			}
			atomic.AddInt32(&completed, 1)
			if progressPercentChan != nil {
				percent := int(atomic.LoadInt32(&completed)) * 100 / pages
				select {
				case progressPercentChan <- percent:
				default:
				}
			}
			return nil
		})
		tasks = append(tasks, task)
	}

	// wait for all submitted jobs to finish
	pool.StopAndWait()

	// collect task errors
	for _, t := range tasks {
		if tErr := t.Wait(); tErr != nil {
			errs = append(errs, tErr)
		}
	}

	// report final progress (non-blocking). If all pages completed, report 100.
	if progressPercentChan != nil {
		completedVal := int(atomic.LoadInt32(&completed))
		percent := completedVal * 100 / pages
		if completedVal >= pages {
			percent = 100
		}
		select {
		case progressPercentChan <- percent:
		default:
		}
	}

	if len(errs) > 0 {
		var b strings.Builder
		for i, e := range errs {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(e.Error())
		}
		return users, fmt.Errorf("encountered %d errors: %s", len(errs), b.String())
	}
	return users, nil
}
