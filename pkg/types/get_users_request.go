package types

// GetUsersRequestPayload represents the request payload used by the GetUsers endpoint.
// Fields are pointers when zero values are ambiguous so the builder can choose
// whether to emit them or not.
type GetUsersRequestPayload struct {
	Offset      *int           `json:"offset"`
	Limit       *int           `json:"limit"`
	Orders      map[string]int `json:"orders"`
	Search      *string        `json:"search"`
	Filters     []any          `json:"filters"`
	ExtraParams map[string]any `json:"extra_params"`
}
