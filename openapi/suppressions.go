package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SuppressionsService manages suppressed email addresses, which no longer
// receive email.
type SuppressionsService struct{ c *Client }

// Suppression is one suppressed email address.
type Suppression struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	Source string `json:"source"`
	// SubAccountID is zero for an account-wide suppression.
	SubAccountID int `json:"subAccountId"`
	// SubAccountName is filled in by the list endpoints only.
	SubAccountName string `json:"subAccountName"`
	Description    string `json:"description"`
	// CreatedAt is a Unix millisecond timestamp.
	CreatedAt int64 `json:"createdAt"`
}

// CreateSuppressionRequest is the payload of Create. Email, Type and Source are
// required. Leave SubAccountID zero to suppress the address account-wide.
type CreateSuppressionRequest struct {
	Email        string `json:"email"`
	Type         string `json:"type"`
	Source       string `json:"source"`
	SubAccountID int    `json:"subAccountId,omitempty"`
	Description  string `json:"description,omitempty"`
}

func (r *CreateSuppressionRequest) validate() error {
	if r == nil || strings.TrimSpace(r.Email) == "" ||
		strings.TrimSpace(r.Type) == "" || strings.TrimSpace(r.Source) == "" {
		return errors.New("suppressions create requires Email, Type and Source")
	}
	return nil
}

// ListSuppressionsParams holds the optional query parameters of List, Bounces
// and Unsubscribes.
type ListSuppressionsParams struct {
	Page    int
	PerPage int
	// Email matches one address exactly.
	Email string
	Type  string
	// Source filters by source; separate several with commas to match any of
	// them. Bounces and Unsubscribes preset the source and ignore this field.
	Source       string
	SubAccountID int
}

func (p *ListSuppressionsParams) values() url.Values {
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
	if p.Email != "" {
		values.Set("email", p.Email)
	}
	if p.Type != "" {
		values.Set("type", p.Type)
	}
	if p.Source != "" {
		values.Set("source", p.Source)
	}
	if p.SubAccountID > 0 {
		values.Set("subAccountId", strconv.Itoa(p.SubAccountID))
	}
	return values
}

// presetSourceValues is values() for the endpoints that fix the source
// server side.
func (p *ListSuppressionsParams) presetSourceValues() url.Values {
	values := p.values()
	values.Del("source")
	return values
}

// LookupSuppressionsParams holds the parameters of Lookup. Email is required.
type LookupSuppressionsParams struct {
	// Email is matched exactly and case-sensitively.
	Email string
	// SubAccountID narrows the result to one sub-account.
	SubAccountID int
	Page         int
	PerPage      int
}

func (p *LookupSuppressionsParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	values.Set("email", p.Email)
	if p.SubAccountID > 0 {
		values.Set("subAccountId", strconv.Itoa(p.SubAccountID))
	}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("perPage", strconv.Itoa(p.PerPage))
	}
	return values
}

func (p *LookupSuppressionsParams) validate() error {
	if p == nil || strings.TrimSpace(p.Email) == "" {
		return errors.New("suppressions lookup requires an Email")
	}
	return nil
}

// List returns a page of suppressions, newest first.
func (s *SuppressionsService) List(params *ListSuppressionsParams) (*Page[Suppression], error) {
	return paged[Suppression](s.c, request{
		capability: "Suppressions.List",
		method:     http.MethodGet,
		path:       "/suppressions",
		query:      params.values(),
	})
}

// Bounces returns a page of addresses suppressed by a bounce rule, newest
// first.
func (s *SuppressionsService) Bounces(params *ListSuppressionsParams) (*Page[Suppression], error) {
	return paged[Suppression](s.c, request{
		capability: "Suppressions.Bounces",
		method:     http.MethodGet,
		path:       "/suppressions/bounces",
		query:      params.presetSourceValues(),
	})
}

// Unsubscribes returns a page of addresses suppressed because they
// unsubscribed, through the List-Unsubscribe header or an unsubscribe link,
// newest first.
func (s *SuppressionsService) Unsubscribes(params *ListSuppressionsParams) (*Page[Suppression], error) {
	return paged[Suppression](s.c, request{
		capability: "Suppressions.Unsubscribes",
		method:     http.MethodGet,
		path:       "/suppressions/unsubscribes",
		query:      params.presetSourceValues(),
	})
}

// Lookup returns every suppression recorded for one address, newest first.
// An address can be suppressed once per sub-account plus account-wide, so
// the result is a list.
func (s *SuppressionsService) Lookup(params *LookupSuppressionsParams) (*Page[Suppression], error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	return paged[Suppression](s.c, request{
		capability: "Suppressions.Lookup",
		method:     http.MethodGet,
		path:       "/suppressions/lookup",
		query:      params.values(),
	})
}

// Create suppresses an address. Creating a suppression that already exists
// returns the existing one.
func (s *SuppressionsService) Create(req *CreateSuppressionRequest) (*Suppression, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[Suppression](s.c, request{
		capability: "Suppressions.Create",
		method:     http.MethodPost,
		path:       "/suppressions",
		body:       req,
	})
}

// Delete removes a suppression so the address can receive email again.
func (s *SuppressionsService) Delete(id int) error {
	return s.c.do(request{
		capability: "Suppressions.Delete",
		method:     http.MethodDelete,
		path:       "/suppressions/" + strconv.Itoa(id),
	})
}
