package openapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SegmentsService manages contact segments: saved, named contact filters.
type SegmentsService struct{ c *Client }

// Segment is a saved contact filter. Filters uses the same structure as the
// filters body of Contacts.Count.
type Segment struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Filters []map[string]any `json:"filters"`
	// Count is the contact count cached at the last refresh and may be stale;
	// use RecipientsCount for a live number.
	Count int `json:"count"`
	// CreatedAt and UpdatedAt are Unix millisecond timestamps.
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

// CreateSegmentRequest is the payload of Create. Name must be unique within
// the account; the API normalizes it (each word capitalized, `[ ] \ "`
// removed). Filters is required: pass an empty, non-nil slice to match every
// contact.
type CreateSegmentRequest struct {
	Name    string           `json:"name"`
	Filters []map[string]any `json:"filters"`
}

func (r *CreateSegmentRequest) validate() error {
	if r == nil || strings.TrimSpace(r.Name) == "" {
		return errors.New("segments create requires a Name")
	}
	// A nil slice marshals to null, which the API rejects.
	if r.Filters == nil {
		return errors.New("segments create requires Filters — pass an empty " +
			"slice to match every contact")
	}
	return nil
}

// UpdateSegmentRequest is the payload of Update. Omitted fields keep their
// current value, and at least one must be set. A nil Filters is omitted; an
// empty, non-nil Filters replaces the filters so the segment matches every
// contact.
type UpdateSegmentRequest struct {
	Name    string
	Filters []map[string]any
}

// MarshalJSON distinguishes a nil Filters (omitted) from an empty one (sent as
// []), which a plain omitempty tag cannot do.
func (r UpdateSegmentRequest) MarshalJSON() ([]byte, error) {
	body := map[string]any{}
	if r.Name != "" {
		body["name"] = r.Name
	}
	if r.Filters != nil {
		body["filters"] = r.Filters
	}
	return json.Marshal(body)
}

func (r *UpdateSegmentRequest) validate() error {
	if r == nil || (r.Name == "" && r.Filters == nil) {
		return errors.New("segments update requires at least one of Name or Filters")
	}
	return nil
}

// ListSegmentsParams holds the optional query parameters of List.
type ListSegmentsParams struct {
	Page    int
	PerPage int
	// Search matches the segment name, case-insensitively.
	Search string
	// Sort is one of name, createdAt, updatedAt or count; prefix with "-" for
	// descending order.
	Sort string
}

func (p *ListSegmentsParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	if p.Sort != "" {
		values.Set("sort", p.Sort)
	}
	return values
}

// segmentRecipientsCount matches the {"data": {"count": N}} body of
// RecipientsCount.
type segmentRecipientsCount struct {
	Count int `json:"count"`
}

// List returns a page of segments.
func (s *SegmentsService) List(params *ListSegmentsParams) (*Page[Segment], error) {
	return paged[Segment](s.c, request{
		capability: "Segments.List",
		method:     http.MethodGet,
		path:       "/segments",
		query:      params.values(),
	})
}

// Get returns one segment by ID.
func (s *SegmentsService) Get(id string) (*Segment, error) {
	return object[Segment](s.c, request{
		capability: "Segments.Get",
		method:     http.MethodGet,
		path:       "/segments/" + url.PathEscape(id),
	})
}

// Create saves a new segment.
func (s *SegmentsService) Create(req *CreateSegmentRequest) (*Segment, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[Segment](s.c, request{
		capability: "Segments.Create",
		method:     http.MethodPost,
		path:       "/segments",
		body:       req,
	})
}

// Update renames a segment and/or replaces its filters.
func (s *SegmentsService) Update(id string, req *UpdateSegmentRequest) (*Segment, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[Segment](s.c, request{
		capability: "Segments.Update",
		method:     http.MethodPut,
		path:       "/segments/" + url.PathEscape(id),
		body:       req,
	})
}

// Delete permanently removes a segment and untags every contact tagged with
// it.
func (s *SegmentsService) Delete(id string) error {
	return s.c.do(request{
		capability: "Segments.Delete",
		method:     http.MethodDelete,
		path:       "/segments/" + url.PathEscape(id),
	})
}

// RecipientsCount returns how many active contacts currently match a segment,
// i.e. who would receive a campaign sent to it.
func (s *SegmentsService) RecipientsCount(id string) (int, error) {
	out, err := object[segmentRecipientsCount](s.c, request{
		capability: "Segments.RecipientsCount",
		method:     http.MethodGet,
		path:       "/segments/" + url.PathEscape(id) + "/recipients-count",
	})
	if err != nil {
		return 0, err
	}
	return out.Count, nil
}
