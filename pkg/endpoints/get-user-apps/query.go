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

type AuthorizedApp struct {
	App struct {
		DomainID    string `json:"domain_id"`
		AppID       int64  `json:"app_id"`
		CreatedTime int64  `json:"created_time"`
		UpdatedTime int64  `json:"updated_time"`
		Disabled    bool   `json:"disabled"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Key         string `json:"key"`
		ThumbnailID string `json:"thumbnail_id"`
		GatewayID   int64  `json:"gateway_id"`

		DestinationSetting struct {
			Type        string `json:"type"`
			IPDef       string `json:"ip_def"`
			PortDef     string `json:"port_def"`
			DynamicMode bool   `json:"dynamic_mode"`
			IP          string `json:"ip"`
			Port        int    `json:"port"`
			DynamicIPs  []struct {
				Value      string `json:"value"`
				RangeStart string `json:"range_start"`
				RangeEnd   string `json:"range_end"`
			} `json:"dynamic_ips"`
			DynamicPorts []struct {
				Value      int `json:"value"`
				RangeStart int `json:"range_start"`
				RangeEnd   int `json:"range_end"`
			} `json:"dynamic_ports"`
			UDPSetting *struct {
				Timeout   int `json:"timeout"`
				LocalPort int `json:"local_port"`
				Meta      any `json:"meta"`
			} `json:"udp_setting,omitempty"`
			TCPSetting *struct {
				Timeout int `json:"timeout"`
				Meta    any `json:"meta"`
			} `json:"tcp_setting,omitempty"`
			HTTPSetting struct {
				Enabled   bool     `json:"enabled"`
				SSL       bool     `json:"ssl"`
				SSMMix    bool     `json:"ssl_mix"`
				HostNames []string `json:"host_names"`
				Meta      any      `json:"meta"`
				Params    []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
					Meta  any    `json:"meta"`
				} `json:"params"`
			} `json:"http_setting"`
			ForwardSetting struct {
				Host string `json:"host"`
				Port int    `json:"port"`
				Meta any    `json:"meta"`
			} `json:"forward_setting"`
			ConnectAgentID string `json:"connect_agent_id"`
			Meta           any    `json:"meta"`
			RangeMode      bool   `json:"range_mode"`
			StartIP        string `json:"start_ip"`
			EndIP          string `json:"end_ip"`
			StartPort      int    `json:"start_port"`
			EndPort        int    `json:"end_port"`
			ProxyID        string `json:"proxy_id"`
			Proxy          any    `json:"proxy"`
		} `json:"destination_setting"`

		Meta              any   `json:"meta"`
		SubApps           []any `json:"sub_apps"`
		DependedAppIDs    []any `json:"depended_app_ids"`
		AutoStartDisabled bool  `json:"auto_start_disabled"`
		Gateway           any   `json:"gateway"`
		DependedApps      []any `json:"depended_apps"`
		IdentityOwnerInfo *struct {
			UserIDs               []any `json:"user_ids"`
			GroupIDs              []any `json:"group_ids"`
			OrganizationUnitIDs   []any `json:"organization_unit_ids"`
			OrganizationUnitInfos []any `json:"organization_unit_infos"`
			Meta                  any   `json:"meta"`
			IdentityVersion       int64 `json:"identity_version"`
			UpdatedTime           int64 `json:"updated_time"`
		} `json:"identity_owner_info,omitempty"`
	} `json:"app"`

	AccessRuleInfo struct {
		AccessRuleID   string `json:"access_rule_id"`
		AccessRuleName string `json:"access_rule_name"`
	} `json:"access_rule_info"`

	AccessPolicyInfos []any `json:"access_policy_infos"`
}

// GetUserApps fetches authorized apps for a user and returns a list of AppInfo.
func GetUserApps(as *types.AppState, userID string) ([]AuthorizedApp, error) {
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

	var root struct {
		Data []AuthorizedApp `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	return root.Data, nil
}
