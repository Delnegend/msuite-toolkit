package add_users_to_group

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"msuite-toolkit/pkg/types"
)

func TestAddUsersToGroup(t *testing.T) {
	const groupID = "6a22b46a8ce03e572462784f"
	const bearerToken = "test-bearer-token"

	var mu sync.Mutex
	var requestCount int
	seenIDs := make([]string, 0)
	var seenMethod string
	var seenPath string
	var seenAuth string
	var seenContentType string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var payload addUsersToGroupPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("unmarshalling request body failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mu.Lock()
		requestCount++
		seenMethod = r.Method
		seenPath = r.URL.Path
		seenAuth = r.Header.Get("Authorization")
		seenContentType = r.Header.Get("Content-Type")
		seenIDs = append(seenIDs, payload.UserIDs...)
		batchSize := len(payload.UserIDs)
		mu.Unlock()

		if batchSize > addBatchSize {
			t.Errorf("batch larger than %d: %d", addBatchSize, batchSize)
		}
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
		BearerToken:        bearerToken,
		WorkerCount:        4,
	}

	// 23 ids -> 3 batches (10, 10, 3)
	ids := make([]string, 0, 23)
	for i := 0; i < 23; i++ {
		ids = append(ids, "user-"+string(rune('a'+i)))
	}

	if err := AddUsersToGroup(as, groupID, ids, nil); err != nil {
		t.Fatalf("AddUsersToGroup failed: %v", err)
	}

	if requestCount != 3 {
		t.Fatalf("expected 3 batched requests, got %d", requestCount)
	}
	if seenMethod != http.MethodPost {
		t.Fatalf("unexpected method: %s", seenMethod)
	}
	if seenPath != "/identity-api/v1/domains/default/groups/"+groupID+"/users/@/all/add" {
		t.Fatalf("unexpected path: %s", seenPath)
	}
	if seenAuth != "Bearer "+bearerToken {
		t.Fatalf("unexpected authorization header: %s", seenAuth)
	}
	if seenContentType != "application/json" {
		t.Fatalf("unexpected content-type header: %s", seenContentType)
	}

	if len(seenIDs) != len(ids) {
		t.Fatalf("expected %d ids across batches, got %d", len(ids), len(seenIDs))
	}
	got := make(map[string]struct{}, len(seenIDs))
	for _, id := range seenIDs {
		got[id] = struct{}{}
	}
	for _, id := range ids {
		if _, ok := got[id]; !ok {
			t.Fatalf("missing id in received batches: %s", id)
		}
	}
}

func TestAddUsersToGroupValidation(t *testing.T) {
	as := &types.AppState{WorkerCount: 1}
	if err := AddUsersToGroup(as, "", []string{"x"}, nil); err == nil {
		t.Fatal("expected error for empty group id")
	}
	if err := AddUsersToGroup(as, "group", nil, nil); err == nil {
		t.Fatal("expected error for empty user ids")
	}
}
