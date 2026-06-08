package delete_enrollement_request

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"msuite-toolkit/pkg/types"
)

func TestDeleteEnrollmentRequests(t *testing.T) {
	const serverUserID = "admin-user-id"
	const bearerToken = "test-bearer-token"

	var seenMethod string
	var seenPath string
	var seenAuth string
	var seenContentType string
	var seenBody []byte

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenMethod = r.Method
		seenPath = r.URL.Path
		seenAuth = r.Header.Get("Authorization")
		seenContentType = r.Header.Get("Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body failed: %v", err)
		}
		seenBody = body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	originalGetHTTPClient := getHTTPClient
	getHTTPClient = func() *http.Client {
		return server.Client()
	}
	defer func() {
		getHTTPClient = originalGetHTTPClient
	}()

	as := &types.AppState{
		AdminPortalAddress: strings.TrimPrefix(server.URL, "https://"),
		AdminUserID:        serverUserID,
		BearerToken:        bearerToken,
	}

	ids := []string{"69e20d7bd7b24bc736073c84", "..."}
	if err := DeleteEnrollmentRequests(as, ids, nil); err != nil {
		t.Fatalf("DeleteEnrollmentRequests failed: %v", err)
	}

	if seenMethod != http.MethodPost {
		t.Fatalf("unexpected method: %s", seenMethod)
	}
	if seenPath != "/enrollment-api/v1/domains/default/enrollment_requests/@all/delete" {
		t.Fatalf("unexpected path: %s", seenPath)
	}
	if seenAuth != "Bearer "+bearerToken {
		t.Fatalf("unexpected authorization header: %s", seenAuth)
	}
	if seenContentType != "application/json" {
		t.Fatalf("unexpected content-type header: %s", seenContentType)
	}

	var payload struct {
		Value []string `json:"value"`
	}
	if err := json.Unmarshal(seenBody, &payload); err != nil {
		t.Fatalf("unmarshalling request body failed: %v", err)
	}
	if len(payload.Value) != len(ids) {
		t.Fatalf("unexpected number of ids: %d", len(payload.Value))
	}
	for i, id := range ids {
		if payload.Value[i] != id {
			t.Fatalf("unexpected id at index %d: %s", i, payload.Value[i])
		}
	}
}
