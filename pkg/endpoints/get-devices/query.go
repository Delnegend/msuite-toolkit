package get_devices

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
	"strings"
	"sync"
	"sync/atomic"

	"github.com/alitto/pond/v2"
)

// GetDevices fetches a batch of devices using the provided request payload.
func GetDevices(as *types.AppState, reqPayload types.QueryRequestPayload) (int, []types.DeviceInfo, error) {
	endpoint := fmt.Sprintf("https://%s/device-api/v1/domains/default/devices/info", as.AdminPortalAddress)
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

	var respPayload struct {
		Data  []types.DeviceInfo `json:"data"`
		Count int                `json:"count"`
	}
	if err := json.Unmarshal(body, &respPayload); err != nil {
		slog.Error("unmarshalling response failed", "err", err)
		return 0, nil, fmt.Errorf("unmarshalling response failed: %w", err)
	}

	return respPayload.Count, respPayload.Data, nil
}

// GetAllDevices fetches all devices by making paginated requests.
func GetAllDevices(as *types.AppState, basePayload types.QueryRequestPayload, progressPercentChan chan<- int) ([]types.DeviceInfo, error) {
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
	total, firstBatch, err := GetDevices(as, p)
	if err != nil {
		slog.Error("initial fetch of devices failed", "err", err)
		return nil, err
	}

	devices := make([]types.DeviceInfo, 0, total)
	mutex := &sync.Mutex{}
	devices = append(devices, firstBatch...)

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
		return devices, nil
	}

	for page := 1; page < pages; page++ {
		offset := page * limit
		task := pool.SubmitErr(func() error {
			pp := basePayload
			pp.Offset = &offset
			pp.Limit = &limit
			_, batch, err := GetDevices(as, pp)
			if err != nil {
				slog.Error("fetching devices failed", "err", err, "offset", offset)
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
				devices = append(devices, batch...)
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
		return devices, fmt.Errorf("encountered %d errors: %s", len(errs), b.String())
	}
	return devices, nil
}
