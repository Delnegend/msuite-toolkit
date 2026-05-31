package get_access_rules

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"msuite-toolkit/pkg/httpclient"
	"msuite-toolkit/pkg/types"
	"net/http"
	"net/url"
)

// GetAccessRules fetches access rules using a QueryRequestPayload for parameters.
func GetAccessRules(as *types.AppState, payload types.QueryRequestPayload) (*types.AccessRuleResponse, error) {
	endpoint := fmt.Sprintf("https://%s/access-api/v1/domains/default/access_rules", as.AdminPortalAddress)

	reqPayloadBytes, err := json.Marshal(payload)
	if err != nil {
		slog.Error("marshalling request payload failed", "err", err)
		return nil, fmt.Errorf("marshalling request payload failed: %w", err)
	}

	values := url.Values{}
	values.Set("ctx.user_id", as.AdminUserID)
	values.Set("request_payload", string(reqPayloadBytes))
	reqURL := endpoint + "?" + values.Encode()

	client := httpclient.GetHTTPClient()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
	if err != nil {
		slog.Error("creating request failed", "err", err)
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+as.BearerToken)

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "err", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body failed", "err", err)
		return nil, fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response types.AccessRuleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("unmarshalling response body failed", "err", err)
		return nil, fmt.Errorf("unmarshalling response body failed: %w", err)
	}

	return &response, nil
}
