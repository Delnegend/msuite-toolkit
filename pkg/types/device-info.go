package types

import (
	"time"
)

type DeviceInfo struct {
	DeviceID    string
	DeviceName  string
	UpdatedTime int64

	OS            string
	OSFamily      string
	ProductName   string
	ProductVendor string
}

func (d *DeviceInfo) UpdatedTimeString() string {
	if d.UpdatedTime == 0 {
		return "Never"
	}
	t := time.Unix(d.UpdatedTime, 0)
	return t.Format("2006-01-02 15:04:05")
}
