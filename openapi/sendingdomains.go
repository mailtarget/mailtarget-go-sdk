package openapi

import (
	"net/http"
	"net/url"
	"strconv"
)

// SendingDomainsService inspects and verifies sending domains.
type SendingDomainsService struct{ c *Client }

// SendingDomain is a domain configured for sending.
type SendingDomain struct {
	ID             int      `json:"id"`
	Domain         string   `json:"domain"`
	Status         string   `json:"status"`
	Assignment     string   `json:"assignment"`
	Types          []string `json:"types"`
	Default        bool     `json:"default"`
	CompanyID      int      `json:"companyId"`
	SubAccountID   int      `json:"subAccountID"`
	SubAccountName string   `json:"subAccountName"`
	CreatedAt      int64    `json:"createdAt"`
}

// ListSendingDomainsParams holds the optional query parameters of List.
type ListSendingDomainsParams struct {
	Page    int
	PerPage int
}

func (p *ListSendingDomainsParams) values() url.Values {
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
	return values
}

// List returns a page of sending domains.
func (s *SendingDomainsService) List(params *ListSendingDomainsParams) (*Page[SendingDomain], error) {
	return paged[SendingDomain](s.c, request{
		capability: "SendingDomains.List",
		method:     http.MethodGet,
		path:       "/domain/sending",
		query:      params.values(),
	})
}

// Get returns one sending domain by ID.
func (s *SendingDomainsService) Get(id int) (*SendingDomain, error) {
	return object[SendingDomain](s.c, request{
		capability: "SendingDomains.Get",
		method:     http.MethodGet,
		path:       "/domain/sending/" + strconv.Itoa(id),
	})
}

// VerifyTXT asks Mailtarget to re-check the domain's TXT records and returns
// the refreshed domain.
func (s *SendingDomainsService) VerifyTXT(id int) (*SendingDomain, error) {
	return object[SendingDomain](s.c, request{
		capability: "SendingDomains.VerifyTXT",
		method:     http.MethodPut,
		path:       "/domain/sending/" + strconv.Itoa(id) + "/verify-txt",
	})
}
