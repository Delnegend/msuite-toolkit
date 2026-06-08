package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	delete_enrollment_request "msuite-toolkit/pkg/endpoints/delete-enrollment-request"
	get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
	get_enrollment_requests_by_ouid "msuite-toolkit/pkg/endpoints/get-enrollment-requests-by-ouid"
	"msuite-toolkit/pkg/types"
	"os"
)

func main() {
	app.Init("")

	as := &app.AppState
	requests := fetchPendingEnrollmentRequests(as)

	deletionCandidates := computeIncludedRequests(requests)
	excludedRequests := computeExcludedRequests(requests, deletionCandidates)

	if len(deletionCandidates) == 0 {
		slog.Info("no pending enrollment requests to delete")
		reportExcludedRequests(excludedRequests)
		return
	}

	reportIncludedRequests(as, deletionCandidates)
	reportExcludedRequests(excludedRequests)

	if as.DryRun {
		slog.Info("dry-run mode: pending enrollment requests were NOT deleted", "count", len(deletionCandidates))
		fmt.Printf("Would delete %d pending enrollment request(s) — see report CSVs\n", len(deletionCandidates))
		return
	}

	deletePendingRequests(as, deletionCandidates, len(requests))
	fmt.Printf("Deleted %d pending enrollment request(s) — see report CSVs\n", len(deletionCandidates))
}

// fetchPendingEnrollmentRequests retrieves pending enrollment requests for the configured OU.
func fetchPendingEnrollmentRequests(as *types.AppState) []get_enrollment_requests.EnrollmentRequest {
	requests, err := get_enrollment_requests_by_ouid.GetEnrollmentRequestsByOuidWithProgress(
		as,
		types.NewQueryRequestBuilder().
			WithFilters([]any{map[string]any{
				"key":   "status",
				"value": get_enrollment_requests.Pending,
			}}).
			Build(),
		as.OrganizationalUnitID,
	)
	if err != nil {
		slog.Error("failed to get enrollment requests", "err", err)
		os.Exit(1)
	}
	if len(requests) == 0 {
		slog.Info("no enrollment requests found")
		os.Exit(0)
	}
	return requests
}

// extractRequestIDs pulls the EnrollmentRequestID from each request into a string slice.
func extractRequestIDs(candidates []get_enrollment_requests.EnrollmentRequest) []string {
	requestIDs := make([]string, 0, len(candidates))
	for _, request := range candidates {
		requestIDs = append(requestIDs, request.EnrollmentRequestID)
	}
	return requestIDs
}

// deletePendingRequests performs the bulk deletion and logs progress.
func deletePendingRequests(as *types.AppState, candidates []get_enrollment_requests.EnrollmentRequest, totalFetched int) {
	requestIDs := extractRequestIDs(candidates)
	slog.Info("deleting pending enrollment requests", "total_fetched", totalFetched, "to_delete", len(requestIDs), "excluded", totalFetched-len(requestIDs))

	if err := delete_enrollment_request.DeleteEnrollmentRequestsWithProgress(as, requestIDs); err != nil {
		slog.Error("deleting pending enrollment requests failed", "err", err)
		os.Exit(1)
	}
}

// computeIncludedRequests filters out non-pending enrollment requests and returns
// only those that are candidates for deletion.
func computeIncludedRequests(requests []get_enrollment_requests.EnrollmentRequest) []get_enrollment_requests.EnrollmentRequest {
	candidates := make([]get_enrollment_requests.EnrollmentRequest, 0, len(requests))
	for _, request := range requests {
		if request.Status == get_enrollment_requests.Pending {
			candidates = append(candidates, request)
		}
	}
	return candidates
}

// computeExcludedRequests returns the requests that are NOT in deletionCandidates
// (i.e. the set difference allRequests - deletionCandidates).
func computeExcludedRequests(
	allRequests,
	deletionCandidates []get_enrollment_requests.EnrollmentRequest,
) []get_enrollment_requests.EnrollmentRequest {
	deleteSet := make(map[string]struct{}, len(deletionCandidates))
	for _, r := range deletionCandidates {
		deleteSet[r.EnrollmentRequestID] = struct{}{}
	}
	excluded := make([]get_enrollment_requests.EnrollmentRequest, 0, len(allRequests)-len(deletionCandidates))
	for _, r := range allRequests {
		if _, ok := deleteSet[r.EnrollmentRequestID]; !ok {
			excluded = append(excluded, r)
		}
	}
	return excluded
}

// reportIncludedRequests writes a CSV with the pending enrollment requests that are
// candidates for deletion (or were deleted). The file is named based on dry_run mode.
func reportIncludedRequests(
	as *types.AppState,
	candidates []get_enrollment_requests.EnrollmentRequest,
) {
	csvFileName := "deleted-pending-enrollment-requests.csv"
	if as.DryRun {
		csvFileName = "to-be-deleted-pending-enrollment-requests.csv"
	}

	csvFile, err := os.Create(csvFileName)
	if err != nil {
		slog.Error("creating csv file failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w := csv.NewWriter(csvFile)
	w.Comma = '|'
	defer w.Flush()

	if err := w.Write([]string{"EnrollmentRequestID", "UserID", "Email", "Status", "DeviceName", "DeviceID", "CreatedTime"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, request := range candidates {
		row := []string{
			request.EnrollmentRequestID,
			request.UserID,
			request.UserInfo.UserEmail,
			string(request.Status),
			request.DeviceName,
			request.DeviceID,
			fmt.Sprintf("%d", request.CreatedTime),
		}
		if err := w.Write(row); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		slog.Error("csv flush failed", "err", err)
		os.Exit(1)
	}

	slog.Info("wrote included pending enrollment requests report", "csv", csvFileName, "count", len(candidates))
}

// reportExcludedRequests writes a CSV with the enrollment requests that were
// excluded from deletion (not in Pending status).
func reportExcludedRequests(excluded []get_enrollment_requests.EnrollmentRequest) {
	if len(excluded) == 0 {
		return
	}

	csvFileName := "excluded-pending-enrollment-requests.csv"

	csvFile, err := os.Create(csvFileName)
	if err != nil {
		slog.Error("creating csv file failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w := csv.NewWriter(csvFile)
	w.Comma = '|'
	defer w.Flush()

	if err := w.Write([]string{"EnrollmentRequestID", "UserID", "Email", "Status", "DeviceName", "DeviceID", "CreatedTime"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, request := range excluded {
		row := []string{
			request.EnrollmentRequestID,
			request.UserID,
			request.UserInfo.UserEmail,
			string(request.Status),
			request.DeviceName,
			request.DeviceID,
			fmt.Sprintf("%d", request.CreatedTime),
		}
		if err := w.Write(row); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		slog.Error("csv flush failed", "err", err)
		os.Exit(1)
	}

	slog.Info("wrote excluded pending enrollment requests report", "csv", csvFileName, "count", len(excluded))
}
