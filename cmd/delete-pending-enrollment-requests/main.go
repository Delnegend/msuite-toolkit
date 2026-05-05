package main

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	delete_enrollement_request "msuite-toolkit/pkg/endpoints/delete-enrollement-request"
	get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
	"msuite-toolkit/pkg/types"
	"os"
)

func main() {
	app.Init("")

	as := &app.AppState
	requests, err := get_enrollment_requests.GetEnrollmentRequestsWithProgress(
		as,
		types.NewQueryRequestBuilder().
			WithFilters([]any{map[string]any{
				"key":   "status",
				"value": get_enrollment_requests.Pending,
			}}).
			Build(),
	)
	if err != nil {
		slog.Error("failed to get enrollment requests", "err", err)
		os.Exit(1)
	}

	deletionCandidates := filterPendingEnrollmentRequests(requests)
	if len(deletionCandidates) == 0 {
		slog.Info("no pending enrollment requests found")
		return
	}

	requestIDs := make([]string, 0, len(deletionCandidates))
	for _, request := range deletionCandidates {
		requestIDs = append(requestIDs, request.EnrollmentRequestID)
	}

	if err := delete_enrollement_request.DeleteEnrollmentRequests(as, requestIDs); err != nil {
		slog.Error("deleting pending enrollment requests failed", "err", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted %d pending enrollment request(s)\n", len(requestIDs))
}

func filterPendingEnrollmentRequests(requests []get_enrollment_requests.EnrollmentRequest) []get_enrollment_requests.EnrollmentRequest {
	deletionCandidates := make([]get_enrollment_requests.EnrollmentRequest, 0, len(requests))
	for _, request := range requests {
		if request.Status == get_enrollment_requests.Pending {
			deletionCandidates = append(deletionCandidates, request)
		}
	}
	return deletionCandidates
}
