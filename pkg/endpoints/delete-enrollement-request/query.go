package delete_enrollement_request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"msuite-toolkit/pkg/httpclient"
	"msuite-toolkit/pkg/types"
	"net/http"
)

var getHTTPClient = httpclient.GetHTTPClient

type deleteEnrollmentRequestsPayload struct {
	Value []string `json:"value"`
}

// DeleteEnrollmentRequests deletes the provided enrollment requests through the bulk delete endpoint.
func DeleteEnrollmentRequests(as *types.AppState, enrollmentRequestIDs []string) error {
	if len(enrollmentRequestIDs) == 0 {
		return fmt.Errorf("no enrollment request ids provided")
	}

	endpoint := fmt.Sprintf("https://%s/enrollment-api/v1/domains/default/enrollment_requests/@all/delete", as.AdminPortalAddress)
	payloadBytes, err := json.Marshal(deleteEnrollmentRequestsPayload{Value: enrollmentRequestIDs})
	if err != nil {
		slog.Error("marshalling request payload failed", "err", err)
		return fmt.Errorf("marshalling request payload failed: %w", err)
	}

	client := getHTTPClient()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body failed", "err", err)
		return fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
