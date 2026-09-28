package openapi

import (
	"net/http"
	"net/url"
	"strconv"
)

// SubAccountsService lists the account's sub-accounts.
type SubAccountsService struct{ c *Client }

// SubAccount is one sub-account of the authenticated company.
type SubAccount struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}

// ListSubAccountsParams holds the optional query parameters of List.
type ListSubAccountsParams struct {
	Page int
	// Size is the page size. This endpoint spells it "size" rather than the
	// "perPage" used elsewhere in the Open API.
	Size      int
	Search    string
	Status    string
	HasAPIKey *bool
}

func (p *ListSubAccountsParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.Size > 0 {
		values.Set("size", strconv.Itoa(p.Size))
	}
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	if p.Status != "" {
		values.Set("status", p.Status)
	}
	if p.HasAPIKey != nil {
		values.Set("hasApiKey", strconv.FormatBool(*p.HasAPIKey))
	}
	return values
}

// List returns a page of sub-accounts.
func (s *SubAccountsService) List(params *ListSubAccountsParams) (*Page[SubAccount], error) {
	return paged[SubAccount](s.c, request{
		capability: "SubAccounts.List",
		method:     http.MethodGet,
		path:       "/sub-accounts",
		query:      params.values(),
	})
}
