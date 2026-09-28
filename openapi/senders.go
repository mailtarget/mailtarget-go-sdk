package openapi

import (
	"net/http"
	"net/url"
)

// SendersService manages sender identities.
type SendersService struct{ c *Client }

// Sender is a sender identity and its domain authentication state.
type Sender struct {
	ID                 string `json:"id"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Domain             string `json:"domain"`
	DomainID           int    `json:"domainId"`
	DomainStatus       string `json:"domainStatus"`
	DomainAssignment   string `json:"domainAssignment"`
	Message            string `json:"message"`
	SPF                bool   `json:"spf"`
	DKIM               bool   `json:"dkim"`
	DMARC              bool   `json:"dmarc"`
	Validate           bool   `json:"validate"`
	Permitted          bool   `json:"permitted"`
	Used               bool   `json:"used"`
	Deleted            bool   `json:"deleted"`
	UseSuggestionValue bool   `json:"useSuggestionValue"`
	AMPStatus          any    `json:"ampStatus"`
	LastValidate       any    `json:"lastValidate"`
	LastUpdateConfig   any    `json:"lastUpdateConfig"`
}

// SenderRequest is the payload of Create and Update. The API uses the same
// shape for both.
type SenderRequest struct {
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	Assignment string `json:"assignment,omitempty"`
}

// SenderDomainStatus reports the authentication state of a sender's domain.
type SenderDomainStatus struct {
	ID               string `json:"id"`
	SPF              bool   `json:"spf"`
	DKIM             bool   `json:"dkim"`
	DMARC            bool   `json:"dmarc"`
	Validate         bool   `json:"validate"`
	Permitted        bool   `json:"permitted"`
	HasSendingDomain bool   `json:"hasSendingDomain"`
	FirstValidate    int64  `json:"firstValidate"`
	LastValidate     int64  `json:"lastValidate"`
	CreatedAt        int64  `json:"createdAt"`
}

// List returns every sender identity. This endpoint is not paginated.
func (s *SendersService) List() ([]Sender, error) {
	return objects[Sender](s.c, request{
		capability: "Senders.List",
		method:     http.MethodGet,
		path:       "/domain/senders",
	})
}

// Get returns one sender identity by ID.
func (s *SendersService) Get(id string) (*Sender, error) {
	return object[Sender](s.c, request{
		capability: "Senders.Get",
		method:     http.MethodGet,
		path:       "/domain/senders/" + url.PathEscape(id),
	})
}

// Create adds a sender identity.
func (s *SendersService) Create(req *SenderRequest) (*Sender, error) {
	return object[Sender](s.c, request{
		capability: "Senders.Create",
		method:     http.MethodPost,
		path:       "/domain/senders",
		body:       req,
	})
}

// Update changes a sender identity.
func (s *SendersService) Update(id string, req *SenderRequest) (*Sender, error) {
	return object[Sender](s.c, request{
		capability: "Senders.Update",
		method:     http.MethodPut,
		path:       "/domain/senders/" + url.PathEscape(id),
		body:       req,
	})
}

// Delete removes a sender identity.
func (s *SendersService) Delete(id string) error {
	return s.c.do(request{
		capability: "Senders.Delete",
		method:     http.MethodDelete,
		path:       "/domain/senders/" + url.PathEscape(id),
	})
}

// CheckDomain reports whether the domain behind an email address is set up for
// sending.
func (s *SendersService) CheckDomain(email string) (*SenderDomainStatus, error) {
	return object[SenderDomainStatus](s.c, request{
		capability: "Senders.CheckDomain",
		method:     http.MethodPost,
		path:       "/domain/senders/check-domain",
		body:       map[string]string{"email": email},
	})
}
