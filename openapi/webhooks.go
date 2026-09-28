package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Webhook authentication methods, used as AuthenticationID. Leaving
// AuthenticationID zero on a request also means no authentication.
const (
	WebhookAuthNone  = 1
	WebhookAuthBasic = 2
	WebhookAuthOAuth = 3
)

// WebhooksService manages webhook subscriptions.
type WebhooksService struct{ c *Client }

// WebhookEvent is an event type a webhook can subscribe to.
type WebhookEvent struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Detail string `json:"detail"`
}

// Webhook is a webhook as returned by Get, Create and Update. Password and
// ClientSecret are write-only and never returned.
type Webhook struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TargetURL string `json:"targetUrl"`
	// SubAccountID is zero for a webhook shared by every sub-account.
	SubAccountID     int            `json:"subAccountId"`
	SubAccountName   string         `json:"subAccountName"`
	Events           []WebhookEvent `json:"events"`
	AuthenticationID int            `json:"authenticationId"`
	Username         string         `json:"username"`
	ClientID         string         `json:"clientId"`
	TokenURL         string         `json:"tokenUrl"`
	IsActive         bool           `json:"isActive"`
	LastStatus       string         `json:"lastStatus"`
	CreatedAt        int64          `json:"createdAt"`
}

// WebhookListItem is a webhook as returned by the list endpoint.
type WebhookListItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TargetURL   string `json:"targetUrl"`
	IsActive    bool   `json:"isActive"`
	LastSuccess int64  `json:"lastSuccess"`
	LastFailure int64  `json:"lastFailure"`
}

// CreateWebhookRequest is the payload of Create. Name, TargetURL and EventIDs
// are required, and TargetURL must be a public http(s) URL. New webhooks are
// always active.
//
// Leave SubAccountID zero to share the webhook with every sub-account. With
// WebhookAuthBasic set Username and Password; with WebhookAuthOAuth set
// ClientID, ClientSecret and TokenURL.
type CreateWebhookRequest struct {
	Name             string `json:"name"`
	TargetURL        string `json:"targetUrl"`
	EventIDs         []int  `json:"eventIds"`
	SubAccountID     int    `json:"subAccountId,omitempty"`
	AuthenticationID int    `json:"authenticationId,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	ClientID         string `json:"clientId,omitempty"`
	ClientSecret     string `json:"clientSecret,omitempty"`
	TokenURL         string `json:"tokenUrl,omitempty"`
}

func (r *CreateWebhookRequest) validate() error {
	if r == nil {
		return errors.New("webhooks create requires a request")
	}
	return validateWebhookRequired("create", r.Name, r.TargetURL, r.EventIDs)
}

// UpdateWebhookRequest is the payload of Update. It replaces the webhook's
// whole configuration, so Name, TargetURL and EventIDs are required just as on
// Create.
//
// Password and IsActive are the exceptions: leave them nil to keep the
// current stored password and active state.
type UpdateWebhookRequest struct {
	Name             string  `json:"name"`
	TargetURL        string  `json:"targetUrl"`
	EventIDs         []int   `json:"eventIds"`
	SubAccountID     int     `json:"subAccountId,omitempty"`
	AuthenticationID int     `json:"authenticationId,omitempty"`
	Username         string  `json:"username,omitempty"`
	Password         *string `json:"password,omitempty"`
	ClientID         string  `json:"clientId,omitempty"`
	ClientSecret     string  `json:"clientSecret,omitempty"`
	TokenURL         string  `json:"tokenUrl,omitempty"`
	IsActive         *bool   `json:"isActive,omitempty"`
}

func (r *UpdateWebhookRequest) validate() error {
	if r == nil {
		return errors.New("webhooks update requires a request")
	}
	return validateWebhookRequired("update", r.Name, r.TargetURL, r.EventIDs)
}

func validateWebhookRequired(op, name, targetURL string, eventIDs []int) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(targetURL) == "" || len(eventIDs) == 0 {
		return errors.New("webhooks " + op + " requires Name, TargetURL and at least one EventID")
	}
	return nil
}

// ListWebhooksParams holds the optional query parameters of List.
type ListWebhooksParams struct {
	Page    int
	PerPage int
	// Search matches the webhook name.
	Search string
	// IsActive filters by active state; nil lists both.
	IsActive *bool
}

func (p *ListWebhooksParams) values() url.Values {
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
	if p.IsActive != nil {
		values.Set("isActive", strconv.FormatBool(*p.IsActive))
	}
	return values
}

// List returns a page of webhooks.
func (s *WebhooksService) List(params *ListWebhooksParams) (*Page[WebhookListItem], error) {
	return paged[WebhookListItem](s.c, request{
		capability: "Webhooks.List",
		method:     http.MethodGet,
		path:       "/webhooks",
		query:      params.values(),
	})
}

// Get returns one webhook by ID.
func (s *WebhooksService) Get(id int) (*Webhook, error) {
	return object[Webhook](s.c, request{
		capability: "Webhooks.Get",
		method:     http.MethodGet,
		path:       "/webhooks/" + strconv.Itoa(id),
	})
}

// Create adds a webhook subscription.
func (s *WebhooksService) Create(req *CreateWebhookRequest) (*Webhook, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[Webhook](s.c, request{
		capability: "Webhooks.Create",
		method:     http.MethodPost,
		path:       "/webhooks",
		body:       req,
	})
}

// Update replaces a webhook's configuration.
func (s *WebhooksService) Update(id int, req *UpdateWebhookRequest) (*Webhook, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[Webhook](s.c, request{
		capability: "Webhooks.Update",
		method:     http.MethodPut,
		path:       "/webhooks/" + strconv.Itoa(id),
		body:       req,
	})
}

// Delete permanently removes a webhook.
func (s *WebhooksService) Delete(id int) error {
	return s.c.do(request{
		capability: "Webhooks.Delete",
		method:     http.MethodDelete,
		path:       "/webhooks/" + strconv.Itoa(id),
	})
}
