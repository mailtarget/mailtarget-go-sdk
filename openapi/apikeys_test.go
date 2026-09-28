package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIKeysService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/api-keys",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":7,"key":"mt_abc","name":"prod","permissionType":"full","subAccountId":3}],
				"meta":{"page":1,"perPage":10,"total":1}
			}`), nil
		})

	page, err := c.APIKeys.List(&ListAPIKeysParams{
		Page: 1, PerPage: 10, Search: "prod", Sort: "createdAt", SubAccountID: 3, PermissionType: "full",
	})

	require.NoError(t, err)
	assert.Equal(t,
		"page=1&perPage=10&permissionType=full&search=prod&sort=createdAt&subAccountId=3", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, 7, page.Items[0].ID)
	assert.Equal(t, "mt_abc", page.Items[0].Key)
	assert.Equal(t, 1, page.Meta.Total)
}

func TestAPIKeysService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/api-keys/7",
		httpmock.NewStringResponder(200, `{"data":{
			"id":7,"name":"prod","allowedIp":"10.0.0.1","companyId":11,
			"permissions":[{"id":1,"name":"contacts.read","detail":"read contacts"}]
		}}`))

	key, err := c.APIKeys.Get(7)

	require.NoError(t, err)
	assert.Equal(t, 7, key.ID)
	assert.Equal(t, "10.0.0.1", key.AllowedIP)
	require.Len(t, key.Permissions, 1)
	assert.Equal(t, "contacts.read", key.Permissions[0].Name)
}

func TestAPIKeysService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/api-keys",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":9,"key":"mt_new"}}`), nil
		})

	key, err := c.APIKeys.Create(&CreateAPIKeyRequest{Name: "ci", PermissionIDs: []int{1, 2}})

	require.NoError(t, err)
	assert.Equal(t, "mt_new", key.Key)
	// permissionIds is required by the API, so it is always sent.
	assert.JSONEq(t, `{"name":"ci","permissionIds":[1,2]}`, gotBody)
}

func TestAPIKeysService_Update(t *testing.T) {
	c := withSecret(t)

	var gotMethod, gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/api-keys/7",
		func(req *http.Request) (*http.Response, error) {
			gotMethod, gotBody = req.Method, readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"id":7,"name":"renamed"}}`), nil
		})

	key, err := c.APIKeys.Update(7, &UpdateAPIKeyRequest{Name: "renamed"})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.JSONEq(t, `{"name":"renamed"}`, gotBody)
	assert.Equal(t, "renamed", key.Name)
}

func TestAPIKeysService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/api-keys/7",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.APIKeys.Delete(7))
}

func TestAPIKeysService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List":   func() error { _, err := c.APIKeys.List(nil); return err },
		"Get":    func() error { _, err := c.APIKeys.Get(1); return err },
		"Create": func() error { _, err := c.APIKeys.Create(&CreateAPIKeyRequest{}); return err },
		"Update": func() error { _, err := c.APIKeys.Update(1, &UpdateAPIKeyRequest{}); return err },
		"Delete": func() error { return c.APIKeys.Delete(1) },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "APIKeys.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
