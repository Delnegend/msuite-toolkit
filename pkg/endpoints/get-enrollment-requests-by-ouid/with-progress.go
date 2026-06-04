package get_enrollment_requests_by_ouid

import (
	"fmt"
	"log/slog"
	"sync"

	get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
)

// GetEnrollmentRequestsByOuidWithProgress fetches all enrollment requests for
// users in the given OU, displaying a progress bar while doing so.
// The signature mirrors GetEnrollmentRequestsWithProgress but accepts an
// additional ouid parameter. When ouid is empty, it delegates directly to
// GetEnrollmentRequestsWithProgress (no OU filtering).
func GetEnrollmentRequestsByOuidWithProgress(
	appState *types.AppState,
	basePayload types.QueryRequestPayload,
	ouid string,
) ([]get_enrollment_requests.EnrollmentRequest, error) {
	if ouid == "" {
		return get_enrollment_requests.GetEnrollmentRequestsWithProgress(appState, basePayload)
	}

	fmt.Println("Fetching enrollment requests by OU...")
	var wg sync.WaitGroup
	progressPercentChan := make(chan int)
	wg.Go(func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
	})
	requests, err := GetAllEnrollmentRequestsByOuid(appState, basePayload, ouid, progressPercentChan)
	close(progressPercentChan)
	wg.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment requests by OU: %w", err)
	}
	slog.Info("fetched enrollment requests by OU", "count", len(requests))
	return requests, nil
}
