package openapi

import "encoding/json"

// Meta carries the pagination block returned by list endpoints.
type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Total   int `json:"total"`
}

// Page is one page of results together with its pagination metadata. Every
// paginated endpoint returns this, so no resource declares a list type of its
// own.
type Page[T any] struct {
	Items []T
	Meta  Meta
}

// dataEnvelope matches the `{data, meta}` wrapper the Open API puts around
// every successful response.
type dataEnvelope struct {
	Data json.RawMessage `json:"data"`
	Meta *Meta           `json:"meta,omitempty"`
}

// errorEnvelope matches dto.ErrorEnvelope from the Open API spec.
type errorEnvelope struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
