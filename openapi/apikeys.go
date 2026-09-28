package openapi

import (
	"net/http"
	"net/url"
	"strconv"
)

// APIKeysService manages the Mailtarget API keys used for sending email.
//
// Note the two credentials involved: these calls authenticate with the Open API
// secret key, while the keys they return are the ones used to send email.
type APIKeysService struct{ c *Client }

// APIKeyPermission is one permission granted to an API key.
type APIKeyPermission struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

// APIKey is an API key as returned by the list endpoint.
type APIKey struct {
	ID             int    `json:"id"`
	Key            string `json:"key"`
	Name           string `json:"name"`
	PermissionType string `json:"permissionType"`
	SubAccountID   int    `json:"subAccountId"`
	SubAccountName string `json:"subAccountName"`
	CreatedAt      int64  `json:"createdAt"`
}

// APIKeyDetail is the fuller record returned when fetching a single key.
type APIKeyDetail struct {
	ID             int                `json:"id"`
	Key            string             `json:"key"`
	Name           string             `json:"name"`
	PermissionType string             `json:"permissionType"`
	Permissions    []APIKeyPermission `json:"permissions"`
	AllowedIP      string             `json:"allowedIp"`
	CompanyID      int                `json:"companyId"`
	CompanyName    string             `json:"companyName"`
	SubAccountID   int                `json:"subAccountId"`
	SubAccountName string             `json:"subAccountName"`
	CreatedAt      int64              `json:"createdAt"`
}

// CreateAPIKeyRequest is the payload of Create. PermissionIDs is required.
type CreateAPIKeyRequest struct {
	Name           string `json:"name,omitempty"`
	PermissionIDs  []int  `json:"permissionIds"`
	PermissionType string `json:"permissionType,omitempty"`
	AllowedIP      string `json:"allowedIp,omitempty"`
	SubAccountID   int    `json:"subAccountId,omitempty"`
}

// UpdateAPIKeyRequest is the payload of Update. Only the fields you set change.
type UpdateAPIKeyRequest struct {
	Name           string `json:"name,omitempty"`
	PermissionIDs  []int  `json:"permissionIds,omitempty"`
	PermissionType string `json:"permissionType,omitempty"`
	AllowedIP      string `json:"allowedIp,omitempty"`
}

// ListAPIKeysParams holds the optional query parameters of List.
type ListAPIKeysParams struct {
	Page           int
	PerPage        int
	Search         string
	Sort           string
	SubAccountID   int
	PermissionType string
}

func (p *ListAPIKeysParams) values() url.Values {
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
	if p.Sort != "" {
		values.Set("sort", p.Sort)
	}
	if p.SubAccountID > 0 {
		values.Set("subAccountId", strconv.Itoa(p.SubAccountID))
	}
	if p.PermissionType != "" {
		values.Set("permissionType", p.PermissionType)
	}
	return values
}

// List returns a page of API keys.
func (s *APIKeysService) List(params *ListAPIKeysParams) (*Page[APIKey], error) {
	return paged[APIKey](s.c, request{
		capability: "APIKeys.List",
		method:     http.MethodGet,
		path:       "/api-keys",
		query:      params.values(),
	})
}

// Get returns one API key by ID.
func (s *APIKeysService) Get(id int) (*APIKeyDetail, error) {
	return object[APIKeyDetail](s.c, request{
		capability: "APIKeys.Get",
		method:     http.MethodGet,
		path:       "/api-keys/" + strconv.Itoa(id),
	})
}

// Create issues a new API key.
func (s *APIKeysService) Create(req *CreateAPIKeyRequest) (*APIKeyDetail, error) {
	return object[APIKeyDetail](s.c, request{
		capability: "APIKeys.Create",
		method:     http.MethodPost,
		path:       "/api-keys",
		body:       req,
	})
}

// Update changes an existing API key.
func (s *APIKeysService) Update(id int, req *UpdateAPIKeyRequest) (*APIKeyDetail, error) {
	return object[APIKeyDetail](s.c, request{
		capability: "APIKeys.Update",
		method:     http.MethodPut,
		path:       "/api-keys/" + strconv.Itoa(id),
		body:       req,
	})
}

// Delete revokes an API key immediately.
func (s *APIKeysService) Delete(id int) error {
	return s.c.do(request{
		capability: "APIKeys.Delete",
		method:     http.MethodDelete,
		path:       "/api-keys/" + strconv.Itoa(id),
	})
}
