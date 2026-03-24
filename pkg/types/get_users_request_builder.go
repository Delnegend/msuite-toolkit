package types

// GetUsersRequestPayloadBuilder constructs GetUsersRequestPayload with sensible defaults.
type GetUsersRequestPayloadBuilder struct {
	offsetSet      bool
	offset         int
	limitSet       bool
	limit          int
	orders         map[string]int
	searchSet      bool
	search         string
	filtersSet     bool
	filters        []any
	extraParamsSet bool
	extraParams    map[string]any
}

// NewGetUsersRequestBuilder returns a fresh builder.
func NewGetUsersRequestBuilder() *GetUsersRequestPayloadBuilder {
	return &GetUsersRequestPayloadBuilder{}
}

// Backwards-compatible helper name to match previous API
func NewGetUsersRequest() *GetUsersRequestPayloadBuilder { return NewGetUsersRequestBuilder() }

func (b *GetUsersRequestPayloadBuilder) WithOffset(v int) *GetUsersRequestPayloadBuilder {
	b.offsetSet = true
	b.offset = v
	return b
}

func (b *GetUsersRequestPayloadBuilder) WithLimit(v int) *GetUsersRequestPayloadBuilder {
	b.limitSet = true
	b.limit = v
	return b
}

func (b *GetUsersRequestPayloadBuilder) WithOrders(m map[string]int) *GetUsersRequestPayloadBuilder {
	b.orders = m
	return b
}

func (b *GetUsersRequestPayloadBuilder) WithSearch(s string) *GetUsersRequestPayloadBuilder {
	b.searchSet = true
	b.search = s
	return b
}

func (b *GetUsersRequestPayloadBuilder) WithFilters(f []any) *GetUsersRequestPayloadBuilder {
	b.filtersSet = true
	b.filters = f
	return b
}

func (b *GetUsersRequestPayloadBuilder) WithFilterByOrgUnitID(ouID string) *GetUsersRequestPayloadBuilder {
	if ouID == "" {
		return b
	}
	filter := map[string]any{
		"key":      "IdentityOwnerInfo.OrganizationUnitInfos.OrganizationUnitId",
		"operator": "equal_to",
		"value":    ouID,
		"origin": map[string]any{
			"key":      "organization_unit_id",
			"operator": "custom",
			"value":    ouID,
		},
	}
	return b.WithFilters([]any{filter})
}

func (b *GetUsersRequestPayloadBuilder) WithExtraParams(p map[string]any) *GetUsersRequestPayloadBuilder {
	b.extraParamsSet = true
	b.extraParams = p
	return b
}

// Build returns a fully-initialized GetUsersRequestPayload with defaults filled
// for any fields the caller didn't set.
func (b *GetUsersRequestPayloadBuilder) Build() GetUsersRequestPayload {
	// defaults
	defOffset := 0
	defLimit := 100
	defSearch := ""
	defOrders := map[string]int{"created_time": 1}
	var filters []any
	var extraParams map[string]any

	p := GetUsersRequestPayload{}
	if b.offsetSet {
		p.Offset = &b.offset
	} else {
		p.Offset = &defOffset
	}
	if b.limitSet {
		p.Limit = &b.limit
	} else {
		p.Limit = &defLimit
	}
	if b.orders != nil {
		p.Orders = b.orders
	} else {
		p.Orders = defOrders
	}
	if b.searchSet {
		p.Search = &b.search
	} else {
		p.Search = &defSearch
	}
	if b.filtersSet {
		p.Filters = b.filters
	} else {
		p.Filters = filters
	}
	if b.extraParamsSet {
		p.ExtraParams = b.extraParams
	} else {
		p.ExtraParams = extraParams
	}
	return p
}
