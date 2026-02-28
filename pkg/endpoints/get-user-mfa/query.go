package get_user_mfa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"msuite-toolkit/pkg/httpclient"
	"msuite-toolkit/pkg/types"
	"net/http"
)

type UserMFAInfo struct {
	TOTP               bool
	EmailOTP           bool
	SMSOTP             bool
	RadiusOTP          bool
	KeystrokeBioAuthen bool
}

func GetUserMFA(as *types.AppState, userID types.UserID) (UserMFAInfo, error) {
	endpoint := fmt.Sprintf("https://%s/identity-authen-api/v1/domains/default/users/%s/mfa", as.AdminPortalAddress, userID)

	client := httpclient.GetHTTPClient()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("creating request failed", "err", err)
		return UserMFAInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+as.BearerToken)

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "err", err)
		return UserMFAInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body failed", "err", err)
		return UserMFAInfo{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("non-2xx response", "status", resp.StatusCode, "body", string(body))
		return UserMFAInfo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respPayload struct {
		Data map[string]struct {
			Disabled bool `json:"disabled"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &respPayload); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return UserMFAInfo{}, err
	}

	info := UserMFAInfo{}
	if v, ok := respPayload.Data["Totp"]; ok {
		info.TOTP = !v.Disabled
	}
	if v, ok := respPayload.Data["SmsOtp"]; ok {
		info.SMSOTP = !v.Disabled
	}
	if v, ok := respPayload.Data["RadiusOtp"]; ok {
		info.RadiusOTP = !v.Disabled
	}
	if v, ok := respPayload.Data["Biometric"]; ok {
		info.KeystrokeBioAuthen = !v.Disabled
	}
	// EmailOtp may not be present in response; default false if missing
	if v, ok := respPayload.Data["EmailOtp"]; ok {
		info.EmailOTP = !v.Disabled
	}

	return info, nil
}
