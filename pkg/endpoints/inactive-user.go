package endpoints

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"net/http"
	"strings"
)

// InactiveUser locks the given user by calling the admin lock endpoint.
// It sends a JSON payload like: {"value":"<user-id>"} and returns an error
// if the server responds with a non-200 status code.
func InactiveUser(as *types.AppState, userID types.UserID) error {
	endpoint := fmt.Sprintf("https://%s/identity-api/v1/domains/default/users/%s/lock", as.AdminPortalAddress, userID)

	payload := fmt.Sprintf(`{"value":"%s"}`, userID)

	client := getHTTPClient()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, endpoint, strings.NewReader(payload))
	if err != nil {
		slog.Error("creating request failed", "err", err)
		return fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+as.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "err", err)
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("inactive user failed", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
