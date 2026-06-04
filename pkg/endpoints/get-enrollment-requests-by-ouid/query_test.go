package get_enrollment_requests_by_ouid

import (
	"testing"

	"msuite-toolkit/pkg/types"

	"github.com/BurntSushi/toml"
)

func TestGetAllEnrollmentRequestsByOuid(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	if appState.OrganizationalUnitID == "" {
		t.Skip("OrganizationalUnitID is empty in config.test.toml, skipping test")
	}

	payload := types.NewQueryRequestBuilder().WithOrders(map[string]int{"created_time": 1, "enrollment_request_id": 1}).Build()
	requests, err := GetAllEnrollmentRequestsByOuid(&appState, payload, appState.OrganizationalUnitID, nil)
	if err != nil {
		t.Fatalf("GetAllEnrollmentRequestsByOuid failed: %v", err)
	}

	t.Logf("Total enrollment requests fetched by OU: %d", len(requests))
	for i, req := range requests {
		t.Logf("Request %d: ID=%s, UserID=%s, Status=%s", i, req.EnrollmentRequestID, req.UserID, req.Status)
	}
}

func TestGetEnrollmentRequestsByOuidWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	if appState.OrganizationalUnitID == "" {
		t.Skip("OrganizationalUnitID is empty in config.test.toml, skipping test")
	}

	payload := types.NewQueryRequestBuilder().WithOrders(map[string]int{"created_time": 1, "enrollment_request_id": 1}).Build()
	requests, err := GetEnrollmentRequestsByOuidWithProgress(&appState, payload, appState.OrganizationalUnitID)
	if err != nil {
		t.Fatalf("GetEnrollmentRequestsByOuidWithProgress failed: %v", err)
	}

	t.Logf("Total enrollment requests fetched by OU with progress: %d", len(requests))
}
