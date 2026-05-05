package types

// QueryRequestPayloadBuilder constructs QueryRequestPayload with sensible defaults.
type QueryRequestPayloadBuilder struct {
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

// NewQueryRequestBuilder returns a fresh builder.
func NewQueryRequestBuilder() *QueryRequestPayloadBuilder {
	return &QueryRequestPayloadBuilder{}
}

// Backwards-compatible helper name to match previous API.
func NewGetDevicesRequestBuilder() *QueryRequestPayloadBuilder { return NewQueryRequestBuilder() }

// Backwards-compatible helper name to match previous API.
func NewGetDevicesRequest() *QueryRequestPayloadBuilder { return NewQueryRequestBuilder() }

func (b *QueryRequestPayloadBuilder) WithOffset(v int) *QueryRequestPayloadBuilder {
	b.offsetSet = true
	b.offset = v
	return b
}

func (b *QueryRequestPayloadBuilder) WithLimit(v int) *QueryRequestPayloadBuilder {
	b.limitSet = true
	b.limit = v
	return b
}

func (b *QueryRequestPayloadBuilder) WithOrders(m map[string]int) *QueryRequestPayloadBuilder {
	b.orders = m
	return b
}

func (b *QueryRequestPayloadBuilder) WithSearch(s string) *QueryRequestPayloadBuilder {
	b.searchSet = true
	b.search = s
	return b
}

func (b *QueryRequestPayloadBuilder) WithFilters(f []any) *QueryRequestPayloadBuilder {
	b.filtersSet = true
	b.filters = f
	return b
}

func (b *QueryRequestPayloadBuilder) WithFilterByOrgUnitID(ouID string) *QueryRequestPayloadBuilder {
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

func (b *QueryRequestPayloadBuilder) WithExtraParams(p map[string]any) *QueryRequestPayloadBuilder {
	b.extraParamsSet = true
	b.extraParams = p
	return b
}

// Build returns a fully-initialized QueryRequestPayload with defaults filled for any fields the caller didn't set.
func (b *QueryRequestPayloadBuilder) Build() QueryRequestPayload {
	defOffset := 0
	defLimit := 100
	defSearch := ""
	defOrders := map[string]int{"updated_time": 1}
	var filters []any
	var extraParams map[string]any

	p := QueryRequestPayload{}
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
