package get_user_devices

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

// GetUserDevices fetches device basic info for the given user and returns parsed DeviceInfo entries.
func GetUserDevices(as *types.AppState, userID string) ([]types.DeviceInfo, error) {
	endpoint := fmt.Sprintf("https://%s/device-api/v1/domains/default/devices/user/%s/info/basic", as.AdminPortalAddress, userID)

	reqPayloadBytes, err := json.Marshal(struct {
		Limit int `json:"limit"`
	}{Limit: -1})
	if err != nil {
		slog.Error("marshalling request payload failed", "err", err)
		return nil, fmt.Errorf("marshalling request payload failed: %w", err)
	}

	values := url.Values{}
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// parse only the fields we care about
	var respPayload struct {
		Data []struct {
			DeviceID    string `json:"device_id"`
			DeviceName  string `json:"device_name"`
			UpdatedTime int64  `json:"updated_time"`
			MetaData    struct {
				OS            string `json:"os"`
				OSFamily      string `json:"os_family"`
				ProductName   string `json:"product_name"`
				ProductVendor string `json:"product_vendor"`
			} `json:"meta_data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &respPayload); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	devices := make([]types.DeviceInfo, 0, len(respPayload.Data))
	for _, d := range respPayload.Data {
		devices = append(devices, types.DeviceInfo{
			DeviceID:      d.DeviceID,
			DeviceName:    d.DeviceName,
			UpdatedTime:   d.UpdatedTime,
			ProductName:   d.MetaData.ProductName,
			ProductVendor: d.MetaData.ProductVendor,
			OS:            d.MetaData.OS,
			OSFamily:      d.MetaData.OSFamily,
		})
	}

	return devices, nil
}
