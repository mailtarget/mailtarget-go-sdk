package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Form statuses accepted by ListFormsParams.Status.
const (
	FormStatusPublished = "published"
	FormStatusDraft     = "draft"
	FormStatusArchived  = "archived"
	// FormStatusUsed matches forms embedded in a landing page or automation.
	FormStatusUsed = "used"
)

// Form usage types, as found in FormUsage.Type.
const (
	FormUsagePage       = "page"
	FormUsageAutomation = "automation"
)

// FormsService manages sign-up forms and their submissions.
type FormsService struct{ c *Client }

// Form is a form as returned by List and Duplicate.
type Form struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Labels           []string `json:"labels"`
	Published        bool     `json:"published"`
	Archived         bool     `json:"archived"`
	Captcha          bool     `json:"captcha"`
	DoubleOptIn      bool     `json:"doubleOptIn"`
	AllowUpdate      bool     `json:"allowUpdate"`
	SubmissionCount  int      `json:"submissionCount"`
	UsedInAutomation bool     `json:"usedInAutomation"`
	UsedInPage       bool     `json:"usedInPage"`
	// CreatedAt and UpdatedAt are Unix millisecond timestamps.
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

// FormFieldOption is one choice of a select, radio or checkbox field.
type FormFieldOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FormField is one field of a form. Name is the key to submit its value under
// in FormSubmissionRequest.Values.
type FormField struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	// Type is the field type, e.g. "inputText".
	Type     string `json:"type"`
	Required bool   `json:"required"`
	// MultiValue fields take an array of option names when submitted.
	MultiValue   bool              `json:"multiValue"`
	ContactField bool              `json:"contactField"`
	Options      []FormFieldOption `json:"options"`
}

// FormUsage is a landing page or automation a form is used in.
type FormUsage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type is FormUsagePage or FormUsageAutomation.
	Type      string `json:"type"`
	Published bool   `json:"published"`
	URL       string `json:"url"`
}

// FormDetail is a form with its fields, as returned by Get, Publish and
// Unpublish.
//
// Fields describes the current (editable) version, while submissions are
// validated against the last published one: after editing a published form,
// publish it again before submitting new or renamed fields.
type FormDetail struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Labels           []string `json:"labels"`
	Published        bool     `json:"published"`
	Archived         bool     `json:"archived"`
	Captcha          bool     `json:"captcha"`
	DoubleOptIn      bool     `json:"doubleOptIn"`
	AllowUpdate      bool     `json:"allowUpdate"`
	SaveToContact    bool     `json:"saveToContact"`
	SubmissionCount  int      `json:"submissionCount"`
	SubmitLimit      int      `json:"submitLimit"`
	UsedInAutomation bool     `json:"usedInAutomation"`
	UsedInPage       bool     `json:"usedInPage"`
	// Fields are in display order.
	Fields []FormField `json:"fields"`
	UsedIn []FormUsage `json:"usedIn"`
	// PublicURL is the hosted form page and EmbedHTML an iframe snippet for
	// it. Both are set once the form has been published and serve the last
	// published version.
	PublicURL string `json:"publicUrl"`
	EmbedHTML string `json:"embedHtml"`
	// CreatedAt and UpdatedAt are Unix millisecond timestamps.
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

// FormSubmission is one stored response to a form.
type FormSubmission struct {
	ID string `json:"id"`
	// Values holds every current field of the form keyed by field name, with
	// "" for a field left empty; several values of one field are joined with
	// ", ".
	Values map[string]string `json:"values"`
	// RemovedFields holds non-empty answers to fields since removed from the
	// form, keyed by their old title.
	RemovedFields map[string]string `json:"removedFields"`
	// CheckedIn is only set for event check-in forms.
	CheckedIn *bool `json:"checkedIn,omitempty"`
	// SubmittedAt is a Unix millisecond timestamp.
	SubmittedAt int64 `json:"submittedAt"`
}

// FormSubmissionRequest is the payload of Submit. Values is keyed by
// FormField.Name; a MultiValue field takes a []string of option names, every
// other field a string. For example:
//
//	map[string]any{
//		"email":     "sarah@example.com",
//		"firstname": "Sarah",
//		"interest":  []string{"Promo", "Newsletter"},
//	}
type FormSubmissionRequest struct {
	Values map[string]any `json:"values"`
}

func (r *FormSubmissionRequest) validate() error {
	// A nil map marshals to null, which the API rejects.
	if r == nil || r.Values == nil {
		return errors.New("forms submit requires Values")
	}
	return nil
}

// ListFormsParams holds the optional query parameters of List.
type ListFormsParams struct {
	Page    int
	PerPage int
	// Search matches the form name, case-insensitively.
	Search string
	// Status keeps only forms in any of these states (FormStatusPublished,
	// FormStatusDraft, FormStatusArchived, FormStatusUsed). Archived forms are
	// excluded unless archived or used is requested.
	Status []string
	// Labels keeps only forms carrying any of these labels.
	Labels []string
}

func (p *ListFormsParams) values() url.Values {
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
	if len(p.Labels) > 0 {
		values.Set("labels", strings.Join(p.Labels, ","))
	}
	return values
}

// ListFormSubmissionsParams holds the optional query parameters of
// Submissions. From and To must be set together or not at all.
type ListFormSubmissionsParams struct {
	Page    int
	PerPage int
	// Search matches the answers to the form's current fields,
	// case-sensitively.
	Search string
	// From and To bound SubmittedAt, inclusive, as Unix millisecond
	// timestamps.
	From int64
	To   int64
}

func (p *ListFormSubmissionsParams) validate() error {
	if p != nil && (p.From > 0) != (p.To > 0) {
		return errors.New("forms submissions requires both From and To, or neither")
	}
	return nil
}

func (p *ListFormSubmissionsParams) values() url.Values {
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
	if p.From > 0 {
		values.Set("from", strconv.FormatInt(p.From, 10))
	}
	if p.To > 0 {
		values.Set("to", strconv.FormatInt(p.To, 10))
	}
	return values
}

// List returns a page of forms, newest first. AMP forms created by campaigns
// are not included.
func (s *FormsService) List(params *ListFormsParams) (*Page[Form], error) {
	return paged[Form](s.c, request{
		capability: "Forms.List",
		method:     http.MethodGet,
		path:       "/forms",
		query:      params.values(),
	})
}

// Get returns one form with its fields and, once published, its PublicURL and
// EmbedHTML.
func (s *FormsService) Get(id string) (*FormDetail, error) {
	return object[FormDetail](s.c, request{
		capability: "Forms.Get",
		method:     http.MethodGet,
		path:       "/forms/" + url.PathEscape(id),
	})
}

// Delete permanently removes a form and its submissions. A published form
// must be unpublished first, and a form used by an automation cannot be
// deleted (both 409).
//
// The API answering `data: false` is reported as an *Error with Code
// ErrCodeOperationFailed.
func (s *FormsService) Delete(id string) error {
	return confirm(s.c, request{
		capability: "Forms.Delete",
		method:     http.MethodDelete,
		path:       "/forms/" + url.PathEscape(id),
	})
}

// Duplicate copies a form into a new, unpublished form named
// "<name> (Duplicate)". Webhook and event check-in settings are not copied.
func (s *FormsService) Duplicate(id string) (*Form, error) {
	return object[Form](s.c, request{
		capability: "Forms.Duplicate",
		method:     http.MethodPost,
		path:       "/forms/" + url.PathEscape(id) + "/duplicate",
	})
}

// Publish makes the current version of a form live at its PublicURL and opens
// it to submissions. Publishing again pushes later edits live.
func (s *FormsService) Publish(id string) (*FormDetail, error) {
	return object[FormDetail](s.c, request{
		capability: "Forms.Publish",
		method:     http.MethodPost,
		path:       "/forms/" + url.PathEscape(id) + "/publish",
	})
}

// Unpublish takes a form offline: its PublicURL stops serving and submissions
// are refused.
func (s *FormsService) Unpublish(id string) (*FormDetail, error) {
	return object[FormDetail](s.c, request{
		capability: "Forms.Unpublish",
		method:     http.MethodPost,
		path:       "/forms/" + url.PathEscape(id) + "/unpublish",
	})
}

// Submissions returns a page of a form's submissions, newest first.
func (s *FormsService) Submissions(id string, params *ListFormSubmissionsParams) (*Page[FormSubmission], error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	return paged[FormSubmission](s.c, request{
		capability: "Forms.Submissions",
		method:     http.MethodGet,
		path:       "/forms/" + url.PathEscape(id) + "/submissions",
		query:      params.values(),
	})
}

// Submit stores one response to a published form, as if a visitor filled it
// in, and creates or updates the contact when the form saves to contacts.
//
// Server-side only: this call needs the Open API secret key, which must never
// be shipped to a browser or mobile app. To collect responses from visitors,
// show the form itself through FormDetail.PublicURL or FormDetail.EmbedHTML.
//
// Values are validated against the last published version of the form, not
// edits made since. Forms with captcha or file upload fields are not
// supported (422); an unpublished or archived form is rejected with 409.
//
// The API answering `data: false` is reported as an *Error with Code
// ErrCodeOperationFailed.
func (s *FormsService) Submit(id string, req *FormSubmissionRequest) error {
	if err := req.validate(); err != nil {
		return err
	}
	return confirm(s.c, request{
		capability: "Forms.Submit",
		method:     http.MethodPost,
		path:       "/forms/" + url.PathEscape(id) + "/submissions",
		body:       req,
	})
}
