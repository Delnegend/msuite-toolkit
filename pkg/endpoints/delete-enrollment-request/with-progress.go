package delete_enrollement_request

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
)

// DeleteEnrollmentRequestsWithProgress deletes the provided enrollment requests
// in batches while rendering a progress bar to stdout.
func DeleteEnrollmentRequestsWithProgress(as *types.AppState, enrollmentRequestIDs []string) error {
	fmt.Println("Deleting enrollment requests...")
	var wg sync.WaitGroup
	progressPercentChan := make(chan int)
	wg.Go(func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
	})
	err := DeleteEnrollmentRequests(as, enrollmentRequestIDs, progressPercentChan)
	close(progressPercentChan)
	wg.Wait()
	if err != nil {
		return fmt.Errorf("failed to delete enrollment requests: %w", err)
	}
	slog.Info("deleted enrollment requests", "count", len(enrollmentRequestIDs))
	return nil
}
