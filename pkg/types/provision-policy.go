package types

type ProvisionPolicyResponse struct {
	Data  []ProvisionPolicy `json:"data"`
	Count int               `json:"count"`
}

type ProvisionPolicy struct {
	DomainID                  string             `json:"domain_id"`
	ProvisionPolicyID         string             `json:"provision_policy_id"`
	CreatedTime               int64              `json:"created_time"`
	UpdatedTime               int64              `json:"updated_time"`
	Disabled                  bool               `json:"disabled"`
	Name                      string             `json:"name"`
	Description               string             `json:"description"`
	Key                       string             `json:"key"`
	Any                       bool               `json:"any"`
	UserIDs                   []string           `json:"user_ids"`
	OrganizationUnitIDs       []any              `json:"organization_unit_ids"`
	OrganizationUnitID        string             `json:"organization_unit_id"`
	GroupIDs                  []string           `json:"group_ids"`
	DeviceIDs                 []any              `json:"device_ids"`
	Subnets                   []any              `json:"subnets"`
	ExtraRules                []any              `json:"extra_rules"`
	Action                    string             `json:"action"`
	PerUserQuota              int                `json:"per_user_quota"`
	ExpiredTime               int64              `json:"expired_time"`
	Priority                  int                `json:"priority"`
	Users                     []UserInfo         `json:"users"`
	OrganizationUnits         []any              `json:"organization_units"`
	OrganizationUnit          any                `json:"organization_unit"`
	Groups                    []any              `json:"groups"`
	Devices                   []any              `json:"devices"`
	IdentityOwnerInfo         *IdentityOwnerInfo `json:"identity_owner_info,omitempty"`
	RequiredExtraAction       string             `json:"required_extra_action"`
	RequiredExtraActionParams string             `json:"required_extra_action_params"`
	CheckDeviceMac            bool               `json:"check_device_mac"`
	DeviceTypes               []string           `json:"device_types"`
}

type IdentityOwnerInfo struct {
	UserIDs               []any                  `json:"user_ids"`
	GroupIDs              []string               `json:"group_ids"`
	OrganizationUnitIDs   []any                  `json:"organization_unit_ids"`
	OrganizationUnitInfos []OrganizationUnitInfo `json:"organization_unit_infos"`
	Meta                  any                    `json:"meta"`
	IdentityVersion       int                    `json:"identity_version"`
	UpdatedTime           int64                  `json:"updated_time"`
}

type OrganizationUnitInfo struct {
	OrganizationUnitID string `json:"organization_unit_id"`
	Version            int    `json:"version"`
	Left               int    `json:"left"`
	Right              int    `json:"right"`
	Level              int    `json:"level"`
	Meta               any    `json:"meta"`
	UpdatedTime        int64  `json:"updated_time"`
}
