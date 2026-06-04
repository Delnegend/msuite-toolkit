package get_enrollment_requests_by_ouid

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"

	get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"

	"github.com/alitto/pond/v2"
)

// GetAllEnrollmentRequestsByOuid fetches all enrollment requests for all users
// belonging to the given organizational unit (OU). It first retrieves all user IDs
// for the OU, then makes concurrent requests to the enrollment requests endpoint
// (one per user ID, since the server only accepts a single user_id filter per call).
func GetAllEnrollmentRequestsByOuid(
	as *types.AppState,
	basePayload types.QueryRequestPayload,
	ouid string,
	progressPercentChan chan<- int,
) ([]get_enrollment_requests.EnrollmentRequest, error) {
	// 1. Get all users for the OU
	usersPayload := types.NewQueryRequestBuilder().
		WithFilterByOrgUnitID(ouid).
		Build()

	users, err := get_users.GetAllUsers(as, usersPayload, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for OU %s: %w", ouid, err)
	}

	if len(users) == 0 {
		slog.Info("no users found for OU", "ouid", ouid)
		return nil, nil
	}

	slog.Info("found users for OU", "count", len(users), "ouid", ouid)

	// 2. Fetch enrollment requests per user in parallel
	pool := pond.NewPool(as.WorkerCount)
	var mu sync.Mutex
	var allRequests []get_enrollment_requests.EnrollmentRequest
	var errsMu sync.Mutex
	var errs []error
	var completed int32
	totalUsers := len(users)

	if progressPercentChan != nil {
		select {
		case progressPercentChan <- 0:
		default:
		}
	}

	tasks := make([]pond.Task, 0, totalUsers)

	for _, user := range users {
		userID := user.UserID
		task := pool.SubmitErr(func() error {
			// Build payload with a single user_id filter.
			// The server only accepts one user_id per call, so we create one request per user.
			userFilter := map[string]any{
				"key":      "user_id",
				"operator": "equal_to",
				"value":    userID,
			}
			// Copy filters to avoid mutating the shared basePayload.Filters slice.
			pp := basePayload
			pp.Filters = append(append([]any{}, basePayload.Filters...), userFilter)

			requests, err := get_enrollment_requests.GetAllEnrollmentRequests(as, pp, nil)
			if err != nil {
				slog.Error("fetching enrollment requests for user failed", "err", err, "user_id", userID)
				atomic.AddInt32(&completed, 1)
				if progressPercentChan != nil {
					percent := int(atomic.LoadInt32(&completed)) * 100 / totalUsers
					select {
					case progressPercentChan <- percent:
					default:
					}
				}
				return err
			}

			if len(requests) > 0 {
				mu.Lock()
				allRequests = append(allRequests, requests...)
				mu.Unlock()
			}

			atomic.AddInt32(&completed, 1)
			if progressPercentChan != nil {
				percent := int(atomic.LoadInt32(&completed)) * 100 / totalUsers
				select {
				case progressPercentChan <- percent:
				default:
				}
			}
			return nil
		})
		tasks = append(tasks, task)
	}

	pool.StopAndWait()

	for _, t := range tasks {
		if tErr := t.Wait(); tErr != nil {
			errsMu.Lock()
			errs = append(errs, tErr)
			errsMu.Unlock()
		}
	}

	if progressPercentChan != nil {
		completedVal := int(atomic.LoadInt32(&completed))
		percent := completedVal * 100 / totalUsers
		if completedVal >= totalUsers {
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
		return allRequests, fmt.Errorf("encountered %d errors: %s", len(errs), b.String())
	}

	return allRequests, nil
}
