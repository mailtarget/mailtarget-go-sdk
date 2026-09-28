package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// LabelsService manages contact labels.
type LabelsService struct{ c *Client }

// Label is a contact label. The API returns its identifier as `_id`.
type Label struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	ContactCount int    `json:"contactCount"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateLabelRequest is the payload of Create, and carries the new name for
// Rename.
type CreateLabelRequest struct {
	Name string `json:"name"`
}

// ListLabelsParams holds the optional query parameters of List.
type ListLabelsParams struct {
	Page    int
	PerPage int
	Search  string
}

func (p *ListLabelsParams) values() url.Values {
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
	return values
}

// List returns a page of labels.
func (s *LabelsService) List(params *ListLabelsParams) (*Page[Label], error) {
	return paged[Label](s.c, request{
		capability: "Labels.List",
		method:     http.MethodGet,
		path:       "/labels",
		query:      params.values(),
	})
}

// Create adds a new label.
func (s *LabelsService) Create(name string) (*Label, error) {
	return object[Label](s.c, request{
		capability: "Labels.Create",
		method:     http.MethodPost,
		path:       "/labels",
		body:       &CreateLabelRequest{Name: name},
	})
}

// Rename changes a label's name, addressing it by its current name.
func (s *LabelsService) Rename(oldName, newName string) (*Label, error) {
	return object[Label](s.c, request{
		capability: "Labels.Rename",
		method:     http.MethodPut,
		path:       "/labels/" + url.PathEscape(oldName),
		body:       &CreateLabelRequest{Name: newName},
	})
}

// Delete removes one or more labels by name. The endpoint takes them as a
// single comma separated path segment, so a name containing a comma cannot be
// addressed.
func (s *LabelsService) Delete(names ...string) error {
	if len(names) == 0 {
		return errors.New("labels delete requires at least one label name")
	}
	escaped := make([]string, 0, len(names))
	for _, name := range names {
		if strings.Contains(name, ",") {
			return errors.New("label name " + name + " contains a comma, which the endpoint cannot address")
		}
		escaped = append(escaped, url.PathEscape(name))
	}
	return s.c.do(request{
		capability: "Labels.Delete",
		method:     http.MethodDelete,
		path:       "/labels/" + strings.Join(escaped, ","),
	})
}
