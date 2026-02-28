package get_user_apps

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

type AppInfo struct {
	AppID int64
	Name  string
}

// GetUserApps fetches authorized apps for a user and returns a list of AppInfo.
func GetUserApps(as *types.AppState, userID string) ([]AppInfo, error) {
	endpoint := fmt.Sprintf("https://%s/sdp-api/v1/domains/default/public/users/%s/detailed_authorized_apps", as.AdminPortalAddress, userID)

	values := url.Values{}
	values.Set("limit", "-1")
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respPayload struct {
		Data []struct {
			App struct {
				AppID int64  `json:"app_id"`
				Name  string `json:"name"`
			} `json:"app"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &respPayload); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	apps := make([]AppInfo, 0, len(respPayload.Data))
	for _, d := range respPayload.Data {
		apps = append(apps, AppInfo{AppID: d.App.AppID, Name: d.App.Name})
	}

	return apps, nil
}
