package types

import (
	"time"
)

type DeviceInfo struct {
	DomainID     string `json:"domain_id"`
	DeviceID     string `json:"device_id"`
	Type         string `json:"type"`
	ActiveUserID string `json:"active_user_id"`
	DeviceName   string `json:"device_name"`
	UserInfos    []struct {
		LastUpdateTime int64    `json:"last_update_time"`
		UserInfo       UserInfo `json:"user_info"`
	} `json:"user_infos"`
	AgentVersion string `json:"agent_version"`
	DeviceStatus string `json:"device_status"`
	PolicyID     string `json:"policy_id"`
	CreatedTime  int64  `json:"created_time"`
	UpdatedTime  int64  `json:"updated_time"`
	MetaData     struct {
		OSVersion     string `json:"os_version"`
		OS            string `json:"os"`
		OSFamily      string `json:"os_family"`
		ProductName   string `json:"product_name"`
		ProductVendor string `json:"product_vendor"`
	} `json:"meta_data"`
	Extra             any `json:"extra"`
	IdentityOwnerInfo struct {
		UserIDs               []any    `json:"user_ids"`
		GroupIDs              []string `json:"group_ids"`
		OrganizationUnitIDs   []string `json:"organization_unit_ids"`
		OrganizationUnitInfos []struct {
			OrganizationUnitID string `json:"organization_unit_id"`
			Version            int    `json:"version"`
			Left               int    `json:"left"`
			Right              int    `json:"right"`
			Level              int    `json:"level"`
			Meta               any    `json:"meta"`
			UpdatedTime        int64  `json:"updated_time"`
		} `json:"organization_unit_infos"`
		Meta            any   `json:"meta"`
		IdentityVersion int   `json:"identity_version"`
		UpdatedTime     int64 `json:"updated_time"`
	} `json:"identity_owner_info"`
}

func (d *DeviceInfo) UpdatedTimeString() string {
	if d.UpdatedTime == 0 {
		return "Never"
	}
	t := time.Unix(d.UpdatedTime, 0)
	return t.Format("2006-01-02 15:04:05")
}
