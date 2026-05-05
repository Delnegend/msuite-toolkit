package get_enrollment_requests

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
)

func GetEnrollmentRequestsWithProgress(appState *types.AppState, basePayload types.QueryRequestPayload) ([]EnrollmentRequest, error) {
	fmt.Println("Fetching enrollment requests...")
	var wg sync.WaitGroup
	progressPercentChan := make(chan int)
	wg.Go(func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
	})
	requests, err := GetAllEnrollmentRequests(appState, basePayload, progressPercentChan)
	close(progressPercentChan)
	wg.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to get all enrollment requests: %w", err)
	}
	slog.Info("fetched enrollment requests", "count", len(requests))
	return requests, nil
}
