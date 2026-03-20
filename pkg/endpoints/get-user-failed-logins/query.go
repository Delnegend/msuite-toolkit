package get_user_failed_logins

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

type FailedLogin struct {
	CreatedTime int64  `json:"created_time"`
	IP          string `json:"ip"`
	DeviceID    string `json:"device_id"`
	LoginName   string `json:"login_name"`
	LoginType   string `json:"login_type"`
	Reason      string `json:"reason"`
}

// GetUserFailedLogins fetches a batch of failed logins for a specific user starting from the given offset with the specified limit.
// It returns the total count of failed logins, the list of FailedLogin, and any error encountered.
func GetUserFailedLogins(as *types.AppState, userID string, offset int, limit int) (int, []FailedLogin, error) {
	endpoint := fmt.Sprintf("https://%s/event-log-api/v1/domains/default/event_logs", as.AdminPortalAddress)

	reqPayloadBytes, err := json.Marshal(struct {
		Offset      int            `json:"offset"`
		Limit       int            `json:"limit"`
		Orders      map[string]int `json:"orders"`
		Search      string         `json:"search"`
		Filters     []any          `json:"filters"`
		ExtraParams map[string]any `json:"extra_params"`
	}{
		Offset: offset,
		Limit:  limit,
		Orders: map[string]int{"created_time": 1},
		Search: "",
		Filters: []any{
			map[string]any{"key": "actor.user.user_id", "operator": "equal_to", "value": userID},
			map[string]any{"key": "result.status", "operator": "equal_to", "value": false},
			map[string]any{"key": "action", "operator": "equal_to", "value": "AuthenLogin"},
		},
		ExtraParams: map[string]any{},
	})
	if err != nil {
		slog.Error("marshalling request payload failed", "err", err)
		return 0, nil, fmt.Errorf("marshalling request payload failed: %w", err)
	}

	values := url.Values{}
	values.Set("ctx.user_id", as.AdminUserID)
	values.Set("request_payload", string(reqPayloadBytes))
	reqURL := endpoint + "?" + values.Encode()

	client := httpclient.GetHTTPClient()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
	if err != nil {
		slog.Error("creating request failed", "err", err)
		return 0, nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+as.BearerToken)

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "err", err)
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body failed", "err", err)
		return 0, nil, fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
		return 0, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respPayload struct {
		Data []struct {
			CreatedTime int64 `json:"created_time"`
			Actor       struct {
				IP     string `json:"ip"`
				Device struct {
					DeviceID string `json:"device_id"`
				} `json:"device"`
				Meta struct {
					LoginName string `json:"login_name"`
					LoginType string `json:"login_type"`
				} `json:"meta"`
			} `json:"actor"`
			Result struct {
				Error struct {
					Data string `json:"data"`
				} `json:"error"`
			} `json:"result"`
		} `json:"data"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal(body, &respPayload); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return 0, nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	failedLogins := make([]FailedLogin, 0, len(respPayload.Data))
	for _, raw := range respPayload.Data {
		reason := raw.Result.Error.Data
		var errData struct {
			ErrorInfo struct {
				Message string `json:"message"`
			} `json:"error_info"`
		}
		if err := json.Unmarshal([]byte(raw.Result.Error.Data), &errData); err == nil {
			if errData.ErrorInfo.Message != "" {
				reason = errData.ErrorInfo.Message
			}
		}

		failedLogins = append(failedLogins, FailedLogin{
			CreatedTime: raw.CreatedTime,
			IP:          raw.Actor.IP,
			DeviceID:    raw.Actor.Device.DeviceID,
			LoginName:   raw.Actor.Meta.LoginName,
			LoginType:   raw.Actor.Meta.LoginType,
			Reason:      reason,
		})
	}

	return respPayload.Count, failedLogins, nil
}
