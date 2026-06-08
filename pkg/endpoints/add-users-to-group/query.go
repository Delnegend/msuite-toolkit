package add_users_to_group

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"msuite-toolkit/pkg/httpclient"
	"msuite-toolkit/pkg/types"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

var getHTTPClient = httpclient.GetHTTPClient

// addBatchSize is the maximum number of user IDs sent to the add-users endpoint
// in a single request.
const addBatchSize = 10

type addUsersToGroupPayload struct {
	UserIDs []string `json:"user_ids"`
}

// AddUsersToGroup adds the provided users to the given group through the bulk
// add endpoint, splitting them into batches of addBatchSize that are sent
// concurrently via a worker pool. When progressPercentChan is non-nil,
// completion progress (0-100) is reported as batches finish.
func AddUsersToGroup(as *types.AppState, groupID string, userIDs []string, progressPercentChan chan<- int) error {
	if groupID == "" {
		return fmt.Errorf("no group id provided")
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("no user ids provided")
	}

	totalBatches := (len(userIDs) + addBatchSize - 1) / addBatchSize

	if progressPercentChan != nil {
		select {
		case progressPercentChan <- 0:
		default:
		}
	}

	pool := pond.NewPool(as.WorkerCount)
	tasks := make([]pond.Task, 0, totalBatches)
	var completed int32

	for batchIndex := 0; batchIndex < totalBatches; batchIndex++ {
		start := batchIndex * addBatchSize
		end := start + addBatchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[start:end]
		batchNumber := batchIndex + 1

		task := pool.SubmitErr(func() error {
			err := addUsersToGroupBatch(as, groupID, batch)
			if err != nil {
				slog.Error("adding users to group batch failed", "batch", batchNumber, "total_batches", totalBatches, "err", err)
			}

			done := atomic.AddInt32(&completed, 1)
			if progressPercentChan != nil {
				percent := int(done) * 100 / totalBatches
				select {
				case progressPercentChan <- percent:
				default:
				}
			}

			if err != nil {
				return fmt.Errorf("adding batch %d/%d failed: %w", batchNumber, totalBatches, err)
			}
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

	if len(errs) > 0 {
		var b strings.Builder
		for i, e := range errs {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(e.Error())
		}
		return fmt.Errorf("%d of %d batches failed: %s", len(errs), totalBatches, b.String())
	}

	return nil
}

// addUsersToGroupBatch adds a single batch of users to the group through the
// bulk add endpoint.
func addUsersToGroupBatch(as *types.AppState, groupID string, userIDs []string) error {
	endpoint := fmt.Sprintf("https://%s/identity-api/v1/domains/default/groups/%s/users/@/all/add", as.AdminPortalAddress, groupID)
	payloadBytes, err := json.Marshal(addUsersToGroupPayload{UserIDs: userIDs})
	if err != nil {
		slog.Error("marshalling request payload failed", "err", err)
		return fmt.Errorf("marshalling request payload failed: %w", err)
	}

	client := getHTTPClient()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		slog.Error("creating request failed", "err", err)
		return fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+as.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "err", err)
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body failed", "err", err)
		return fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
