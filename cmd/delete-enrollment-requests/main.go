package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	delete_enrollement_request "msuite-toolkit/pkg/endpoints/delete-enrollment-request"
	finduserbyemail "msuite-toolkit/pkg/endpoints/find-user-by-email"
	get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
	get_enrollment_requests_by_ouid "msuite-toolkit/pkg/endpoints/get-enrollment-requests-by-ouid"
	"msuite-toolkit/pkg/types"
	"os"
	"sync"

	"github.com/alitto/pond/v2"
)

func main() {
	app.Init("")

	as := &app.AppState
	requireOrganizationalUnitID(as)

	requests := fetchAllEnrollmentRequests(as)

	deletionCandidates := computeIncludedRequests(requests, buildExcludeUserIDs(as))
	excludedRequests := computeExcludedRequests(requests, deletionCandidates)

	if len(deletionCandidates) == 0 {
		slog.Info("no enrollment requests to delete after exclusions")
		reportExcludedRequests(excludedRequests)
		return
	}

	reportIncludedRequests(as, deletionCandidates)
	reportExcludedRequests(excludedRequests)

	if as.DryRun {
		slog.Info("dry-run mode: enrollment requests were NOT deleted", "count", len(deletionCandidates))
		fmt.Printf("Would delete %d enrollment request(s) — see report CSVs\n", len(deletionCandidates))
		return
	}

	deleteRequests(as, deletionCandidates, len(requests))
	fmt.Printf("Deleted %d enrollment request(s) — see report CSVs\n", len(deletionCandidates))
}

// requireOrganizationalUnitID exits if the OU ID is not set in config.
func requireOrganizationalUnitID(as *types.AppState) {
	if as.OrganizationalUnitID == "" {
		slog.Error("organizational_unit_id is required in config")
		os.Exit(1)
	}
}

// fetchAllEnrollmentRequests retrieves all enrollment requests for the configured OU.
func fetchAllEnrollmentRequests(as *types.AppState) []get_enrollment_requests.EnrollmentRequest {
	requests, err := get_enrollment_requests_by_ouid.GetEnrollmentRequestsByOuidWithProgress(
		as,
		types.NewQueryRequestBuilder().Build(),
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

// deleteRequests performs the bulk deletion and logs progress.
func deleteRequests(as *types.AppState, candidates []get_enrollment_requests.EnrollmentRequest, totalFetched int) {
	requestIDs := extractRequestIDs(candidates)
	slog.Info("deleting enrollment requests", "total_fetched", totalFetched, "to_delete", len(requestIDs), "excluded", totalFetched-len(requestIDs))

	if err := delete_enrollement_request.DeleteEnrollmentRequestsWithProgress(as, requestIDs); err != nil {
		slog.Error("deleting enrollment requests failed", "err", err)
		os.Exit(1)
	}
}

// computeIncludedRequests filters out enrollment requests belonging to excluded
// users and returns only those that are candidates for deletion.
func computeIncludedRequests(
	allRequests []get_enrollment_requests.EnrollmentRequest,
	excludeUserIDs map[string]struct{},
) []get_enrollment_requests.EnrollmentRequest {
	candidates := make([]get_enrollment_requests.EnrollmentRequest, 0, len(allRequests))
	for _, request := range allRequests {
		isExplicitlyExcluded := false
		if _, existed := excludeUserIDs[request.UserID]; existed {
			isExplicitlyExcluded = true
		}
		if isExplicitlyExcluded {
			continue
		}

		if !request.IsDesktop() {
			continue
		}

		candidates = append(candidates, request)
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

// reportExcludedRequests writes a CSV with the enrollment requests that were
// excluded from deletion.
func reportExcludedRequests(excluded []get_enrollment_requests.EnrollmentRequest) {
	if len(excluded) == 0 {
		return
	}

	csvFileName := "excluded-enrollment-requests.csv"

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

	if err := w.Write([]string{
		"EnrollmentRequestID",
		"UserID",
		"Email",
		"Status",
		"DeviceName",
		"DeviceID",
		"Device OS",
		"CreatedTime",
	}); err != nil {
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
			request.DeviceOS,
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

	slog.Info("wrote excluded enrollment requests report", "csv", csvFileName, "count", len(excluded))
}

// reportIncludedRequests writes a CSV with the enrollment requests that are
// candidates for deletion (or were deleted). The file is named based on dry_run mode.
func reportIncludedRequests(
	as *types.AppState,
	candidates []get_enrollment_requests.EnrollmentRequest,
) {
	csvFileName := "deleted-enrollment-requests.csv"
	if as.DryRun {
		csvFileName = "to-be-deleted-enrollment-requests.csv"
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

	if err := w.Write([]string{
		"EnrollmentRequestID",
		"UserID",
		"Email",
		"Status",
		"DeviceName",
		"DeviceID",
		"Device OS",
		"CreatedTime",
	}); err != nil {
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
			request.DeviceOS,
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

	slog.Info("wrote included enrollment requests report", "csv", csvFileName, "count", len(candidates))
}

// buildExcludeUserIDs resolves each email in as.ExcludeEmails to a user ID in
// parallel and returns the set of user IDs to exclude from deletion.
func buildExcludeUserIDs(as *types.AppState) map[string]struct{} {
	pool := pond.NewPool(as.WorkerCount)
	var mu sync.Mutex
	excludeUserIDs := make(map[string]struct{}, len(as.ExcludeEmails))

	for _, email := range as.ExcludeEmails {
		pool.Submit(func() {
			user, err := finduserbyemail.FindUserByEmail(as, email)
			if err != nil {
				slog.Error("could not resolve email to user ID", "email", email, "err", err)
				return
			}
			mu.Lock()
			excludeUserIDs[user.UserID] = struct{}{}
			mu.Unlock()
			slog.Info("excluding enrollment requests for user", "user_id", user.UserID, "email", email)
		})
	}

	pool.StopAndWait()
	return excludeUserIDs
}
