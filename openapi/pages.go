package openapi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Page statuses accepted by ListPagesParams.Status.
const (
	PageStatusPublished = "published"
	PageStatusDraft     = "draft"
	// PageStatusDefault matches the domain's default page or its 404 page.
	PageStatusDefault = "default"
)

// Page categories, as found in LandingPage.Category and accepted by
// ListPagesParams.Category.
const (
	PageCategoryLanding  = "landing"
	PageCategoryThankYou = "thank_you"
)

// Sort orders accepted by ListPagesParams.Sort. A leading "-" sorts
// descending; the API default is PageSortCreatedAtDesc.
const (
	PageSortCreatedAt     = "createdAt"
	PageSortCreatedAtDesc = "-createdAt"
	PageSortName          = "name"
	PageSortNameDesc      = "-name"
	PageSortUpdatedAt     = "updatedAt"
	PageSortUpdatedAtDesc = "-updatedAt"
)

// PagesService manages landing pages. The models are named LandingPage* to
// keep them apart from the generic paginated result Page[T].
type PagesService struct{ c *Client }

// LandingPage is a landing page as returned by List and Duplicate.
type LandingPage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Category is PageCategoryLanding or PageCategoryThankYou.
	Category     string   `json:"category"`
	Labels       []string `json:"labels"`
	ContainsForm bool     `json:"containsForm"`
	// Published reports whether the page is live at URL.
	Published bool `json:"published"`
	// IsDefault reports whether the page is served at the company's domain
	// root.
	IsDefault bool   `json:"default"`
	Thumbnail string `json:"thumbnail"`
	// URL is where the page is (or will be, once published) served.
	URL              string `json:"url"`
	VisitCount       int64  `json:"visitCount"`
	UniqueVisitCount int64  `json:"uniqueVisitCount"`
	// CreatedAt and UpdatedAt are Unix millisecond timestamps.
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

// LandingPageFormRef is a form embedded in a landing page, as listed in
// LandingPageDetail.Forms.
type LandingPageFormRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LandingPageDetail is a landing page with its URL and embedded forms, as
// returned by Get, Publish, Unpublish and SetDefault. The page builder
// content is not included.
type LandingPageDetail struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Category is PageCategoryLanding or PageCategoryThankYou.
	Category     string   `json:"category"`
	Labels       []string `json:"labels"`
	ContainsForm bool     `json:"containsForm"`
	Published    bool     `json:"published"`
	IsDefault    bool     `json:"default"`
	Thumbnail    string   `json:"thumbnail"`
	// URL is filled for drafts too but only serves once the page is
	// published; the default page is served at the domain root.
	URL string `json:"url"`
	// Permalink is the page's path on the domain; the default page is served
	// at the domain root instead.
	Permalink string               `json:"permalink"`
	Forms     []LandingPageFormRef `json:"forms"`
	// UpdatedAt is a Unix millisecond timestamp.
	UpdatedAt int64 `json:"updatedAt"`
}

// LandingPageForm is a form embedded in a landing page, as returned by Forms.
// Its ID works with the FormsService methods.
type LandingPageForm struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Labels          []string `json:"labels"`
	SubmissionCount int64    `json:"submissionCount"`
	// CreatedAt and UpdatedAt are Unix millisecond timestamps.
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

// ListPagesParams holds the optional query parameters of List.
type ListPagesParams struct {
	Page    int
	PerPage int
	// Search matches the page name, case-insensitively.
	Search string
	// Status keeps only pages in any of these states (PageStatusPublished,
	// PageStatusDraft, PageStatusDefault).
	Status []string
	// Category is PageCategoryLanding or PageCategoryThankYou.
	Category string
	// Sort is one of the PageSort* constants; the API default is
	// PageSortCreatedAtDesc.
	Sort string
}

func (p *ListPagesParams) values() url.Values {
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
	if len(p.Status) > 0 {
		values.Set("status", strings.Join(p.Status, ","))
	}
	if p.Category != "" {
		values.Set("category", p.Category)
	}
	if p.Sort != "" {
		values.Set("sort", p.Sort)
	}
	return values
}

// List returns a page of landing pages, newest first unless Sort says
// otherwise. Page templates are not included.
func (s *PagesService) List(params *ListPagesParams) (*Page[LandingPage], error) {
	return paged[LandingPage](s.c, request{
		capability: "Pages.List",
		method:     http.MethodGet,
		path:       "/pages",
		query:      params.values(),
	})
}

// Get returns one landing page with its URL and the forms embedded in it.
// URL is filled for drafts too but only serves once the page is published;
// the default page is served at the domain root.
func (s *PagesService) Get(id string) (*LandingPageDetail, error) {
	return object[LandingPageDetail](s.c, request{
		capability: "Pages.Get",
		method:     http.MethodGet,
		path:       "/pages/" + url.PathEscape(id),
	})
}

// Delete permanently removes a landing page. The default page cannot be
// deleted (409); make another page the default first. Forms embedded in the
// page are not deleted.
//
// The API answering `data: false` is reported as an *Error with Code
// ErrCodeOperationFailed.
func (s *PagesService) Delete(id string) error {
	return confirm(s.c, request{
		capability: "Pages.Delete",
		method:     http.MethodDelete,
		path:       "/pages/" + url.PathEscape(id),
	})
}

// Duplicate copies a landing page with its content into a new, unpublished
// page named "<name> (copy)". Embedded forms are shared with the original,
// not copied.
func (s *PagesService) Duplicate(id string) (*LandingPage, error) {
	return object[LandingPage](s.c, request{
		capability: "Pages.Duplicate",
		method:     http.MethodPost,
		path:       "/pages/" + url.PathEscape(id) + "/duplicate",
	})
}

// Forms returns the forms embedded in a landing page. A form embedded more
// than once is listed once. Each ID works with the FormsService methods,
// which need the forms.read / forms.write scopes. This endpoint is not
// paginated.
func (s *PagesService) Forms(id string) ([]LandingPageForm, error) {
	return objects[LandingPageForm](s.c, request{
		capability: "Pages.Forms",
		method:     http.MethodGet,
		path:       "/pages/" + url.PathEscape(id) + "/forms",
	})
}

// Publish makes the current version of a landing page live at its URL.
// Publishing again pushes later edits live. Publishing a draft counts against
// the plan's published-page quota (422 when it is used up).
func (s *PagesService) Publish(id string) (*LandingPageDetail, error) {
	return object[LandingPageDetail](s.c, request{
		capability: "Pages.Publish",
		method:     http.MethodPost,
		path:       "/pages/" + url.PathEscape(id) + "/publish",
	})
}

// Unpublish takes a landing page offline: its URL stops serving. The default
// page cannot be unpublished (409); make another page the default first.
func (s *PagesService) Unpublish(id string) (*LandingPageDetail, error) {
	return object[LandingPageDetail](s.c, request{
		capability: "Pages.Unpublish",
		method:     http.MethodPost,
		path:       "/pages/" + url.PathEscape(id) + "/unpublish",
	})
}

// SetDefault serves a published landing page at the company's domain root,
// replacing the current default page, which stays published at its own URL.
// The page must be published (409 otherwise). Setting the current default
// page again is a no-op.
func (s *PagesService) SetDefault(id string) (*LandingPageDetail, error) {
	return object[LandingPageDetail](s.c, request{
		capability: "Pages.SetDefault",
		method:     http.MethodPost,
		path:       "/pages/" + url.PathEscape(id) + "/set-default",
	})
}
