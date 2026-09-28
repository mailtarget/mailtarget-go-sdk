package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SubAccountsService manages the account's sub-accounts.
type SubAccountsService struct{ c *Client }

// SubAccount is one sub-account of the authenticated company.
type SubAccount struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}

// SubAccountAPIKey is an API key attached to a sub-account. The key value is
// not included; fetch it with APIKeys.Get.
type SubAccountAPIKey struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SubAccountDomain is a domain attached to a sub-account.
type SubAccountDomain struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
}

// SubAccountDetail is the fuller record returned when fetching, creating or
// updating one sub-account.
type SubAccountDetail struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Status     string             `json:"status"`
	IPPoolID   int                `json:"ipPoolId"`
	IPPoolName string             `json:"ipPoolName"`
	APIKeys    []SubAccountAPIKey `json:"apiKeys"`
	Domains    []SubAccountDomain `json:"domains"`
	CreatedAt  int64              `json:"createdAt"`
}

// CreateSubAccountRequest is the payload of Create. Name is required and must
// be unique within the account. Leave IPPoolID zero to use the default pool.
type CreateSubAccountRequest struct {
	Name     string `json:"name"`
	IPPoolID int    `json:"ipPoolId,omitempty"`
}

func (r *CreateSubAccountRequest) validate() error {
	if r == nil || strings.TrimSpace(r.Name) == "" {
		return errors.New("sub-accounts create requires a Name")
	}
	return nil
}

// UpdateSubAccountRequest is the payload of Update. Only the fields you set
// change, and at least one must be set. Status may be "Active" or
// "Suspended".
type UpdateSubAccountRequest struct {
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty"`
	IPPoolID int    `json:"ipPoolId,omitempty"`
}

func (r *UpdateSubAccountRequest) validate() error {
	if r == nil || (r.Name == "" && r.Status == "" && r.IPPoolID == 0) {
		return errors.New("sub-accounts update requires at least one of Name, Status or IPPoolID")
	}
	return nil
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

// Get returns one sub-account by ID, including its API keys and domains.
func (s *SubAccountsService) Get(id int) (*SubAccountDetail, error) {
	return object[SubAccountDetail](s.c, request{
		capability: "SubAccounts.Get",
		method:     http.MethodGet,
		path:       "/sub-accounts/" + strconv.Itoa(id),
	})
}

// Create adds a new, active sub-account. To give it an API key, call
// APIKeys.Create with the returned ID as SubAccountID.
func (s *SubAccountsService) Create(req *CreateSubAccountRequest) (*SubAccountDetail, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[SubAccountDetail](s.c, request{
		capability: "SubAccounts.Create",
		method:     http.MethodPost,
		path:       "/sub-accounts",
		body:       req,
	})
}

// Update changes a sub-account's name, status or IP pool. The primary
// sub-account and terminated sub-accounts cannot be updated.
func (s *SubAccountsService) Update(id int, req *UpdateSubAccountRequest) (*SubAccountDetail, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return object[SubAccountDetail](s.c, request{
		capability: "SubAccounts.Update",
		method:     http.MethodPut,
		path:       "/sub-accounts/" + strconv.Itoa(id),
		body:       req,
	})
}
