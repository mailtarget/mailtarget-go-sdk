package openapi

import (
	"net/http"
	"net/url"
	"strconv"
)

// ContactsService exposes the contact endpoints.
type ContactsService struct{ c *Client }

// SubscribeSource describes where a contact subscribed from.
type SubscribeSource struct {
	Channel string `json:"channel,omitempty"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
}

// Contact is a contact record.
type Contact struct {
	ID                string           `json:"id"`
	Email             string           `json:"email"`
	Name              string           `json:"name"`
	Firstname         string           `json:"firstname"`
	Lastname          string           `json:"lastname"`
	Phone             string           `json:"phone"`
	Company           string           `json:"company"`
	City              string           `json:"city"`
	Country           string           `json:"country"`
	Gender            string           `json:"gender"`
	BirthDate         string           `json:"birthDate"`
	DayOfBirth        string           `json:"dayOfBirth"`
	Note              string           `json:"note"`
	State             string           `json:"state"`
	Funnel            string           `json:"funnel"`
	Index             string           `json:"index"`
	DeletedState      string           `json:"deletedState"`
	CreatedAt         string           `json:"createdAt"`
	Labels            []string         `json:"labels"`
	Segments          []string         `json:"segments"`
	Campaigns         []any            `json:"campaigns"`
	Scenarios         []any            `json:"scenarios"`
	CustomPeriod      []any            `json:"customPeriod"`
	CustomField       map[string]any   `json:"customField"`
	SubscribeSource   *SubscribeSource `json:"subscribeSource"`
	Tester            bool             `json:"tester"`
	UnsubscribeOrigin string           `json:"unsubscribeOrigin"`
	UnsubscribeReason string           `json:"unsubscribeReason"`
}

// CreateContactRequest is the payload of Create. Email is required by the API.
// Unset fields are omitted so a partially filled request never blanks out
// server side values.
type CreateContactRequest struct {
	Email             string           `json:"email"`
	Name              string           `json:"name,omitempty"`
	Firstname         string           `json:"firstname,omitempty"`
	Lastname          string           `json:"lastname,omitempty"`
	Phone             string           `json:"phone,omitempty"`
	Company           string           `json:"company,omitempty"`
	City              string           `json:"city,omitempty"`
	Country           string           `json:"country,omitempty"`
	Gender            string           `json:"gender,omitempty"`
	DayOfBirth        string           `json:"dayOfBirth,omitempty"`
	Note              string           `json:"note,omitempty"`
	State             string           `json:"state,omitempty"`
	Funnel            string           `json:"funnel,omitempty"`
	Index             string           `json:"index,omitempty"`
	DeletedState      string           `json:"deletedState,omitempty"`
	Labels            []string         `json:"labels,omitempty"`
	Segments          []string         `json:"segments,omitempty"`
	Campaigns         []any            `json:"campaigns,omitempty"`
	Scenarios         []any            `json:"scenarios,omitempty"`
	CustomPeriod      []any            `json:"customPeriod,omitempty"`
	CustomField       map[string]any   `json:"customField,omitempty"`
	SubscribeSource   *SubscribeSource `json:"subscribeSource,omitempty"`
	Tester            bool             `json:"tester,omitempty"`
	UnsubscribeOrigin string           `json:"unsubscribeOrigin,omitempty"`
	UnsubscribeReason string           `json:"unsubscribeReason,omitempty"`
}

// UpdateContactRequest is the payload of Update. Only the fields you set are
// changed.
type UpdateContactRequest struct {
	Email       string         `json:"email,omitempty"`
	Firstname   string         `json:"firstname,omitempty"`
	Lastname    string         `json:"lastname,omitempty"`
	Phone       string         `json:"phone,omitempty"`
	Company     string         `json:"company,omitempty"`
	City        string         `json:"city,omitempty"`
	Country     string         `json:"country,omitempty"`
	Gender      string         `json:"gender,omitempty"`
	DayOfBirth  string         `json:"dayOfBirth,omitempty"`
	Note        string         `json:"note,omitempty"`
	Labels      []string       `json:"labels,omitempty"`
	CustomField map[string]any `json:"customField,omitempty"`
}

// ListContactsParams holds the optional query parameters of List.
type ListContactsParams struct {
	// Page defaults to 1 server side.
	Page int
	// PerPage defaults to 20 server side and is capped at 100.
	PerPage int
	// Search matches against email or name.
	Search string
}

func (p *ListContactsParams) values() url.Values {
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

// contactCount matches the bare {"count": N} body of the count endpoint.
type contactCount struct {
	Count int `json:"count"`
}

// List returns a page of contacts.
func (s *ContactsService) List(params *ListContactsParams) (*Page[Contact], error) {
	return paged[Contact](s.c, request{
		capability: "Contacts.List",
		method:     http.MethodGet,
		path:       "/contacts",
		query:      params.values(),
	})
}

// Get returns a single contact by ID.
func (s *ContactsService) Get(id string) (*Contact, error) {
	return object[Contact](s.c, request{
		capability: "Contacts.Get",
		method:     http.MethodGet,
		path:       "/contacts/" + url.PathEscape(id),
	})
}

// GetByEmail returns a single contact looked up by email address.
func (s *ContactsService) GetByEmail(email string) (*Contact, error) {
	return object[Contact](s.c, request{
		capability: "Contacts.GetByEmail",
		method:     http.MethodGet,
		path:       "/contacts/email/" + url.PathEscape(email),
	})
}

// Create adds a new contact.
func (s *ContactsService) Create(req *CreateContactRequest) (*Contact, error) {
	return object[Contact](s.c, request{
		capability: "Contacts.Create",
		method:     http.MethodPost,
		path:       "/contacts",
		body:       req,
	})
}

// Update changes an existing contact.
func (s *ContactsService) Update(id string, req *UpdateContactRequest) (*Contact, error) {
	return object[Contact](s.c, request{
		capability: "Contacts.Update",
		method:     http.MethodPut,
		path:       "/contacts/" + url.PathEscape(id),
		body:       req,
	})
}

// Delete removes a contact by ID.
func (s *ContactsService) Delete(id string) error {
	return s.c.do(request{
		capability: "Contacts.Delete",
		method:     http.MethodDelete,
		path:       "/contacts/" + url.PathEscape(id),
	})
}

// Count returns how many contacts match filter. Pass nil to count everything.
func (s *ContactsService) Count(filter map[string]any) (int, error) {
	if filter == nil {
		filter = map[string]any{}
	}
	out, err := bare[contactCount](s.c, request{
		capability: "Contacts.Count",
		method:     http.MethodPost,
		path:       "/contacts/count",
		body:       filter,
	})
	if err != nil {
		return 0, err
	}
	return out.Count, nil
}
