package get_enrollment_requests

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetEnrollmentRequests(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	payload := types.NewQueryRequestBuilder().WithOrders(map[string]int{"created_time": 1, "enrollment_request_id": 1}).Build()
	count, requests, err := GetEnrollmentRequests(&appState, payload)
	if err != nil {
		t.Fatalf("GetEnrollmentRequests failed: %v", err)
	}

	t.Logf("Enrollment requests count in batch: %d (total: %d)", len(requests), count)
}

func TestGetEnrollmentRequestsWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	payload := types.NewQueryRequestBuilder().WithOrders(map[string]int{"created_time": 1, "enrollment_request_id": 1}).Build()
	requests, err := GetEnrollmentRequestsWithProgress(&appState, payload)
	if err != nil {
		t.Fatalf("GetEnrollmentRequestsWithProgress failed: %v", err)
	}
	t.Logf("Total enrollment requests fetched with progress: %d", len(requests))
}
