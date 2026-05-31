package types

// AccessRuleResponse represents the response for access rules listing endpoints.
type AccessRuleResponse struct {
	Data  []AccessRule `json:"data"`
	Count int          `json:"count"`
}

type AccessRule struct {
	DomainID                           string                  `json:"domain_id"`
	AccessRuleID                       string                  `json:"access_rule_id"`
	CreatedTime                        int64                   `json:"created_time"`
	UpdatedTime                        int64                   `json:"updated_time"`
	Disabled                           bool                    `json:"disabled"`
	Name                               string                  `json:"name"`
	Description                        string                  `json:"description"`
	Key                                string                  `json:"key"`
	Sources                            []AccessRuleSource      `json:"sources"`
	Destinations                       []AccessRuleDestination `json:"destinations"`
	Action                             string                  `json:"action"`
	AccessDuration                     int                     `json:"access_duration"`
	StartTime                          int64                   `json:"start_time"`
	EndTime                            int64                   `json:"end_time"`
	HoursOfDay                         int                     `json:"hours_of_day"`
	DaysOfWeek                         int                     `json:"days_of_week"`
	Priority                           int                     `json:"priority"`
	AccessPolicyIDs                    []any                   `json:"access_policy_ids"`
	DefaultSourceCondition             DefaultSourceCondition  `json:"default_source_condition"`
	BypassEnforcedAccessPolicies       bool                    `json:"bypass_enforced_access_policies"`
	BypassEnforcedDeviceBaselines      bool                    `json:"bypass_enforced_device_baselines"`
	AccessPolicies                     []any                   `json:"access_policies"`
	IdentityOwnerInfo                  any                     `json:"identity_owner_info"`
	SourceUserEmails                   []any                   `json:"source_user_emails"`
	SourceOrganizationUnitDisplayNames []any                   `json:"source_organization_unit_display_names"`
	SourceGroupDisplayNames            []any                   `json:"source_group_display_names"`
	SourceOthers                       []any                   `json:"source_others"`
	DestinationAppNames                []any                   `json:"destination_app_names"`
	DestinationResourceNames           []any                   `json:"destination_resource_names"`
	DestinationOthers                  []any                   `json:"destination_others"`
	AccessPolicyNames                  []any                   `json:"access_policy_names"`
	GroupAccessRequestsByUser          bool                    `json:"group_access_requests_by_user"`
	InactiveDuration                   int                     `json:"inactive_duration"`
	LastUsedTime                       int64                   `json:"last_used_time"`
	Deactivated                        bool                    `json:"deactivated"`
}

type AccessRuleSource struct {
	SourceID           string   `json:"source_id"`
	UserID             string   `json:"user_id"`
	OrganizationUnitID string   `json:"organization_unit_id"`
	GroupIDs           []any    `json:"group_ids"`
	DeviceID           string   `json:"device_id"`
	DeviceHardwareID   string   `json:"device_hardware_id"`
	SdpClientID        int      `json:"sdp_client_id"`
	IPFilter           any      `json:"ip_filter"`
	MFAs               []any    `json:"mfas"`
	UserAttributes     any      `json:"user_attributes"`
	ExtraAttributes    any      `json:"extra_attributes"`
	Condition          any      `json:"condition"`
	CustomExpression   string   `json:"custom_expression"`
	ExtraRules         []any    `json:"extra_rules"`
	User               UserInfo `json:"user"`
	OrganizationUnit   any      `json:"organization_unit"`
	Groups             []any    `json:"groups"`
}

type AccessRuleDestination struct {
	DestinationID string `json:"destination_id"`
	Type          string `json:"type"`
	ID            string `json:"id"`
	SdpApp        SdpApp `json:"sdp_app"`
	Resource      any    `json:"resource"`
}

type SdpApp struct {
	AppID              int                `json:"app_id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	ThumbnailID        string             `json:"thumbnail_id"`
	DestinationSetting DestinationSetting `json:"destination_setting"`
}

type DestinationSetting struct {
	Type        string `json:"type"`
	IPDef       string `json:"ip_def"`
	PortDef     string `json:"port_def"`
	DynamicMode bool   `json:"dynamic_mode"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	HTTPEnabled bool   `json:"http_enabled"`
	HTTPSSL     bool   `json:"http_ssl"`
}

type DefaultSourceCondition struct {
	AllowedDeviceTypes              []any `json:"allowed_device_types"`
	RequiredDeviceCompliedBaselines []any `json:"required_device_complied_baselines"`
	MinimumTrustPoint               int   `json:"minimum_trust_point"`
}
