package get_user_device_last_ip

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

// GetUserDeviceLastIP searches the event_logs endpoint for the given userID and deviceID
// and returns the first found actor.ip. It pages through results one page at a time
// (limit set to 10) and only requests the next page if no IP is found in the current page.
// Returns empty string and nil error if no IP is found.
func GetUserDeviceLastIP(as *types.AppState, userID types.UserID, deviceID types.DeviceID) (types.IPAddress, error) {
	endpoint := fmt.Sprintf("https://%s/event-log-api/v1/domains/default/event_logs", as.AdminPortalAddress)

	limit := 10
	offset := 0

	client := httpclient.GetHTTPClient()

	for {
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
				map[string]any{"key": "actor.device.device_id", "operator": "equal_to", "value": deviceID},
			},
			ExtraParams: map[string]any{},
		})
		if err != nil {
			slog.Error("marshalling request payload failed", "err", err)
			return "", err
		}

		values := url.Values{}
		values.Set("ctx.user_id", as.AdminUserID)
		values.Set("request_payload", string(reqPayloadBytes))
		reqURL := endpoint + "?" + values.Encode()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
		if err != nil {
			slog.Error("creating request failed", "err", err)
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+as.BearerToken)

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("request failed", "err", err)
			return "", err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			slog.Error("reading response body failed", "err", err)
			return "", err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
			return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var respPayload struct {
			Data []struct {
				Actor struct {
					IP string `json:"ip"`
				} `json:"actor"`
			} `json:"data"`
			Count int `json:"count"`
		}
		if err := json.Unmarshal(body, &respPayload); err != nil {
			slog.Error("unmarshalling response failed", "err", err)
			return "", err
		}

		// scan current page for any actor.ip
		for _, e := range respPayload.Data {
			if e.Actor.IP != "" {
				return e.Actor.IP, nil
			}
		}

		// no IP found in current page; check if there are more pages
		offset += limit
		if respPayload.Count == 0 || offset >= respPayload.Count {
			break
		}
		// otherwise loop and fetch next page
	}

	// not found
	return "", nil
}
