package openapi

import (
	"net/http"
	"net/url"
	"strconv"
)

// TemplatesService reads the account's email templates. A template ID can be
// passed as CampaignRequest.TemplateID to build a campaign from it.
type TemplatesService struct{ c *Client }

// Template is a template as returned by the list endpoint.
type Template struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	EmailType string   `json:"emailType"`
	CSSTheme  string   `json:"cssTheme"`
	Labels    []string `json:"labels"`
	OwnerName string   `json:"ownerName"`
	Thumbnail string   `json:"thumbnail"`
	UpdatedAt string   `json:"updatedAt"`
}

// TemplateDetail is the full template, including its content.
type TemplateDetail struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	EmailType string   `json:"emailType"`
	Content   string   `json:"content"`
	Body      string   `json:"body"`
	CSS       string   `json:"css"`
	CSSTheme  string   `json:"cssTheme"`
	Labels    []string `json:"labels"`
	OwnerName string   `json:"ownerName"`
	Thumbnail string   `json:"thumbnail"`
	UpdatedAt string   `json:"updatedAt"`
}

// ListTemplatesParams holds the optional query parameters of List.
type ListTemplatesParams struct {
	Page      int
	PerPage   int
	Search    string
	Type      string
	Sort      string
	Category  string
	EmailType string
}

func (p *ListTemplatesParams) values() url.Values {
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
	for key, value := range map[string]string{
		"search": p.Search, "type": p.Type, "sort": p.Sort,
		"category": p.Category, "emailType": p.EmailType,
	} {
		if value != "" {
			values.Set(key, value)
		}
	}
	return values
}

// List returns a page of templates.
func (s *TemplatesService) List(params *ListTemplatesParams) (*Page[Template], error) {
	return paged[Template](s.c, request{
		capability: "Templates.List",
		method:     http.MethodGet,
		path:       "/templates",
		query:      params.values(),
	})
}

// Get returns one template with its content.
func (s *TemplatesService) Get(id string) (*TemplateDetail, error) {
	return object[TemplateDetail](s.c, request{
		capability: "Templates.Get",
		method:     http.MethodGet,
		path:       "/templates/" + url.PathEscape(id),
	})
}
