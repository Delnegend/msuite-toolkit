package types

// QueryRequestPayload represents the request payload used by list-style endpoints.
type QueryRequestPayload struct {
	Offset      *int           `json:"offset"`
	Limit       *int           `json:"limit"`
	Orders      map[string]int `json:"orders"`
	Search      *string        `json:"search"`
	Filters     []any          `json:"filters"`
	ExtraParams map[string]any `json:"extra_params"`
}
