package main

import (
	"testing"

	get_enrollment_requests "msuite-toolkit/pkg/endpoints/get-enrollment-requests"
)

func TestComputeIncludedRequests(t *testing.T) {
	requests := []get_enrollment_requests.EnrollmentRequest{
		{EnrollmentRequestID: "1", Status: get_enrollment_requests.Pending},
		{EnrollmentRequestID: "2", Status: get_enrollment_requests.Approved},
		{EnrollmentRequestID: "3", Status: get_enrollment_requests.Rejected},
		{EnrollmentRequestID: "4", Status: get_enrollment_requests.Pending},
		{EnrollmentRequestID: "5", Status: "unknown"},
	}

	filtered := computeIncludedRequests(requests)
	if len(filtered) != 2 {
		t.Fatalf("unexpected filtered request count: got %d want 2", len(filtered))
	}
	if filtered[0].EnrollmentRequestID != "1" || filtered[1].EnrollmentRequestID != "4" {
		t.Fatalf("unexpected filtered requests: %#v", filtered)
	}
}

func TestComputeExcludedRequests(t *testing.T) {
	all := []get_enrollment_requests.EnrollmentRequest{
		{EnrollmentRequestID: "1", Status: get_enrollment_requests.Pending},
		{EnrollmentRequestID: "2", Status: get_enrollment_requests.Approved},
		{EnrollmentRequestID: "3", Status: get_enrollment_requests.Rejected},
		{EnrollmentRequestID: "4", Status: get_enrollment_requests.Pending},
	}
	included := []get_enrollment_requests.EnrollmentRequest{
		{EnrollmentRequestID: "1", Status: get_enrollment_requests.Pending},
		{EnrollmentRequestID: "4", Status: get_enrollment_requests.Pending},
	}

	excluded := computeExcludedRequests(all, included)
	if len(excluded) != 2 {
		t.Fatalf("unexpected excluded request count: got %d want 2", len(excluded))
	}
	if excluded[0].EnrollmentRequestID != "2" || excluded[1].EnrollmentRequestID != "3" {
		t.Fatalf("unexpected excluded requests: %#v", excluded)
	}
}
