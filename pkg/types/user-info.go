package types

import (
	"strconv"
	"time"
)

type UserInfo struct {
	UserID         string            `json:"user_id"`
	Name           string            `json:"name"`
	DisplayName    string            `json:"display_name"`
	UserEmail      string            `json:"email"`
	PhoneNumber    string            `json:"phone_number"`
	Roles          []string          `json:"roles"`
	OwnerInfo      IdentityOwnerInfo `json:"owner_info"`
	Meta           any               `json:"meta"`
	Attributes     map[string]string `json:"attributes"`
	LoginSessionID string            `json:"login_session_id"`
	IsLocked       bool              `json:"is_locked"`
	Description    string            `json:"description"`
	Type           string            `json:"type"`
	Language       string            `json:"language"`
	CreatedTime    int64             `json:"created_time"`
}

// LastLoginTime attempts to extract the last login time from the user's attributes.
func (u UserInfo) LastLoginTime() *time.Time {
	if lastLoginStr, ok := u.Attributes["last_login_time"]; ok {
		if lastLoginInt, err := strconv.ParseInt(lastLoginStr, 10, 64); err == nil {
			t := time.Unix(lastLoginInt, 0)
			return &t
		}
	}
	return nil
}
