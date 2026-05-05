package get_enrollment_requests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	get_devices "msuite-toolkit/pkg/endpoints/get-devices"
	"msuite-toolkit/pkg/httpclient"
	"msuite-toolkit/pkg/types"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

type EnrollmentRequestResponse struct {
	Data  []EnrollmentRequest `json:"data"`
	Count int                 `json:"count"`
}

type EnrollmentRequestStatus string

const (
	Pending  EnrollmentRequestStatus = "Pending"
	Approved EnrollmentRequestStatus = "Approved"
	Rejected EnrollmentRequestStatus = "Rejected"
)

type EnrollmentRequest struct {
	DomainID                 string                  `json:"domain_id"`
	EnrollmentRequestID      string                  `json:"enrollment_request_id"`
	CreatedTime              int64                   `json:"created_time"`
	UpdatedTime              int64                   `json:"updated_time"`
	Disabled                 bool                    `json:"disabled"`
	Name                     string                  `json:"name"`
	Description              string                  `json:"description"`
	Status                   EnrollmentRequestStatus `json:"status"`
	IP                       string                  `json:"ip"`
	UserID                   string                  `json:"user_id"`
	DeviceID                 string                  `json:"device_id"`
	DeviceMACs               []string                `json:"device_macs"`
	DeviceName               string                  `json:"device_name"`
	DeviceType               string                  `json:"device_type"`
	DeviceOS                 string                  `json:"device_os"`
	DeviceOSVersion          string                  `json:"device_os_version"`
	DeviceModel              string                  `json:"device_model"`
	DeviceManufacturer       string                  `json:"device_manufacturer"`
	DeviceTLSCSR             string                  `json:"device_tls_csr"`
	SDPTLSCSR                string                  `json:"sdp_tls_csr"`
	DecisionMessage          string                  `json:"decision_message"`
	DecisionTime             int64                   `json:"decision_time"`
	DecisionAdminID          string                  `json:"decision_admin_id"`
	CredentialsEncryptionKey string                  `json:"credentials_encryption_key"`
	DeviceTLSCert            string                  `json:"device_tls_cert"`
	DeviceTLSKey             string                  `json:"device_tls_key"`
	SDPClientID              int64                   `json:"sdp_client_id"`
	ProvisionPolicyID        string                  `json:"provision_policy_id"`
	Meta                     struct {
		Baseboard                  *string `json:"baseboard,omitempty"`
		BIOS                       *string `json:"bios,omitempty"`
		BuildBranch                *string `json:"build_branch,omitempty"`
		BuildNumber                *string `json:"build_number,omitempty"`
		ConnectivityManagerEnabled *string `json:"connectivity_manager_enabled,omitempty"`
		CPU                        *string `json:"cpu,omitempty"`
		GPU                        *string `json:"gpu,omitempty"`
		Memory                     *string `json:"memory,omitempty"`
		Network                    *string `json:"network,omitempty"`
		OS                         *string `json:"os,omitempty"`
		Product                    *string `json:"product,omitempty"`
		RequestLanguageAgent       string  `json:"request_language_agent"`
		SelfProtectEnabled         *string `json:"self_protect_enabled,omitempty"`
		SimulatorID                *string `json:"simulator_id,omitempty"`
	} `json:"meta"`
	UserInfo          types.UserInfo `json:"user_info"`
	DecisionAdminInfo struct {
		UserID         string   `json:"user_id"`
		Name           string   `json:"name"`
		DisplayName    string   `json:"display_name"`
		Email          string   `json:"email"`
		PhoneNumber    string   `json:"phone_number"`
		Roles          []any    `json:"roles"`
		OwnerInfo      any      `json:"owner_info"`
		Meta           struct{} `json:"meta"`
		Attributes     struct{} `json:"attributes"`
		LoginSessionID string   `json:"login_session_id"`
		IsLocked       bool     `json:"is_locked"`
		Description    string   `json:"description"`
		Type           string   `json:"type"`
		Language       string   `json:"language"`
		CreatedTime    int64    `json:"created_time"`
	} `json:"decision_admin_info"`
	ProvisionPolicy any `json:"provision_policy"`
}

func (er *EnrollmentRequest) DeviceInfo(as *types.AppState) (*types.DeviceInfo, error) {
	_, devicesInfo, err := get_devices.GetDevices(
		as,
		types.
			NewQueryRequestBuilder().
			WithFilters([]any{map[string]any{
				"key":      "device_id",
				"operator": "equal_to",
				"value":    er.DeviceID,
			}}).
			Build(),
	)
	if err != nil {
		slog.Error("fetching device info failed", "err", err, "device_id", er.DeviceID)
		return nil, fmt.Errorf("fetching device info failed: %w", err)
	}
	if len(devicesInfo) < 1 {
		return nil, fmt.Errorf("no device info found for device_id: %s", er.DeviceID)
	}

	return &devicesInfo[0], nil
}

// GetEnrollmentRequests fetches a batch of enrollment requests using the provided request payload.
func GetEnrollmentRequests(as *types.AppState, reqPayload types.QueryRequestPayload) (int, []EnrollmentRequest, error) {

	endpoint := fmt.Sprintf("https://%s/enrollment-api/v1/domains/default/enrollment_requests", as.AdminPortalAddress)
	reqPayloadBytes, err := json.Marshal(reqPayload)
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

	var response EnrollmentRequestResponse
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return 0, nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	return response.Count, response.Data, nil
}

// GetAllEnrollmentRequests fetches all enrollment requests by making paginated requests.
func GetAllEnrollmentRequests(as *types.AppState, basePayload types.QueryRequestPayload, progressPercentChan chan<- int) ([]EnrollmentRequest, error) {
	pool := pond.NewPool(as.WorkerCount)

	limit := 100

	if progressPercentChan != nil {
		select {
		case progressPercentChan <- 0:
		default:
		}
	}

	zero := 0
	p := basePayload
	p.Offset = &zero
	p.Limit = &limit
	total, firstBatch, err := GetEnrollmentRequests(as, p)
	if err != nil {
		slog.Error("initial fetch of enrollment requests failed", "err", err)
		return nil, err
	}

	requests := make([]EnrollmentRequest, 0, total)
	mutex := &sync.Mutex{}
	requests = append(requests, firstBatch...)

	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	var errs []error
	tasks := make([]pond.Task, 0, max(0, pages-1))

	var completed int32 = 1
	if progressPercentChan != nil {
		percent := int(atomic.LoadInt32(&completed)) * 100 / pages
		select {
		case progressPercentChan <- percent:
		default:
		}
	}

	if pages == 1 {
		if progressPercentChan != nil {
			select {
			case progressPercentChan <- 100:
			default:
			}
		}
		return requests, nil
	}

	for page := 1; page < pages; page++ {
		offset := page * limit
		task := pool.SubmitErr(func() error {
			pp := basePayload
			pp.Offset = &offset
			pp.Limit = &limit
			_, batch, err := GetEnrollmentRequests(as, pp)
			if err != nil {
				slog.Error("fetching enrollment requests failed", "err", err, "offset", offset)
				atomic.AddInt32(&completed, 1)
				if progressPercentChan != nil {
					percent := int(atomic.LoadInt32(&completed)) * 100 / pages
					select {
					case progressPercentChan <- percent:
					default:
					}
				}
				return err
			}
			if len(batch) != 0 {
				mutex.Lock()
				requests = append(requests, batch...)
				mutex.Unlock()
			}
			atomic.AddInt32(&completed, 1)
			if progressPercentChan != nil {
				percent := int(atomic.LoadInt32(&completed)) * 100 / pages
				select {
				case progressPercentChan <- percent:
				default:
				}
			}
			return nil
		})
		tasks = append(tasks, task)
	}

	pool.StopAndWait()

	for _, t := range tasks {
		if tErr := t.Wait(); tErr != nil {
			errs = append(errs, tErr)
		}
	}

	if progressPercentChan != nil {
		completedVal := int(atomic.LoadInt32(&completed))
		percent := completedVal * 100 / pages
		if completedVal >= pages {
			percent = 100
		}
		select {
		case progressPercentChan <- percent:
		default:
		}
	}

	if len(errs) > 0 {
		var b strings.Builder
		for i, e := range errs {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(e.Error())
		}
		return requests, fmt.Errorf("encountered %d errors: %s", len(errs), b.String())
	}
	return requests, nil
}
