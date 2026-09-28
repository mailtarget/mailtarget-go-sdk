package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubAccountsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/sub-accounts",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":3,"name":"marketing","status":"active","createdAt":1750000000}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	hasKey := true
	page, err := c.SubAccounts.List(&ListSubAccountsParams{
		Page: 1, Size: 25, Search: "mark", Status: "active", HasAPIKey: &hasKey,
	})

	require.NoError(t, err)
	// This endpoint spells the page size "size", not "perPage".
	assert.Equal(t, "hasApiKey=true&page=1&search=mark&size=25&status=active", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, 3, page.Items[0].ID)
	assert.Equal(t, "marketing", page.Items[0].Name)
}

func TestSubAccountsService_ListHasAPIKeyFalseIsSent(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/sub-accounts",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.RawQuery
			return httpmock.NewStringResponse(200, `{"data":[]}`), nil
		})

	hasKey := false
	_, err := c.SubAccounts.List(&ListSubAccountsParams{HasAPIKey: &hasKey})

	require.NoError(t, err)
	// A pointer distinguishes "filter by false" from "do not filter".
	assert.Equal(t, "hasApiKey=false", gotQuery)
}

func TestSubAccountsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/sub-accounts/3",
		httpmock.NewStringResponder(200, `{"data":{
			"id":3,"name":"marketing","status":"Active","ipPoolId":2,"ipPoolName":"shared",
			"apiKeys":[{"id":7,"name":"ci"}],"domains":[{"id":5,"domain":"mtarget.co"}]
		}}`))

	sub, err := c.SubAccounts.Get(3)

	require.NoError(t, err)
	assert.Equal(t, "shared", sub.IPPoolName)
	assert.Equal(t, []SubAccountAPIKey{{ID: 7, Name: "ci"}}, sub.APIKeys)
	assert.Equal(t, []SubAccountDomain{{ID: 5, Domain: "mtarget.co"}}, sub.Domains)
}

func TestSubAccountsService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/sub-accounts",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":9,"name":"sales","status":"Active"}}`), nil
		})

	sub, err := c.SubAccounts.Create(&CreateSubAccountRequest{Name: "sales"})

	require.NoError(t, err)
	// An unset IP pool is omitted so the default pool applies.
	assert.JSONEq(t, `{"name":"sales"}`, gotBody)
	assert.Equal(t, 9, sub.ID)
}

func TestSubAccountsService_Update(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/sub-accounts/9",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"id":9,"name":"sales","status":"Suspended"}}`), nil
		})

	sub, err := c.SubAccounts.Update(9, &UpdateSubAccountRequest{Status: "Suspended"})

	require.NoError(t, err)
	assert.JSONEq(t, `{"status":"Suspended"}`, gotBody)
	assert.Equal(t, "Suspended", sub.Status)
}

func TestSubAccountsService_ValidatesRequests(t *testing.T) {
	c := withSecret(t)

	_, errCreateNil := c.SubAccounts.Create(nil)
	_, errCreateBlank := c.SubAccounts.Create(&CreateSubAccountRequest{Name: "  "})
	_, errUpdateNil := c.SubAccounts.Update(9, nil)
	_, errUpdateEmpty := c.SubAccounts.Update(9, &UpdateSubAccountRequest{})

	assert.ErrorContains(t, errCreateNil, "requires a Name")
	assert.ErrorContains(t, errCreateBlank, "requires a Name")
	assert.ErrorContains(t, errUpdateNil, "at least one of")
	assert.ErrorContains(t, errUpdateEmpty, "at least one of")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSubAccountsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List": func() error { _, err := c.SubAccounts.List(nil); return err },
		"Get":  func() error { _, err := c.SubAccounts.Get(1); return err },
		"Create": func() error {
			_, err := c.SubAccounts.Create(&CreateSubAccountRequest{Name: "a"})
			return err
		},
		"Update": func() error {
			_, err := c.SubAccounts.Update(1, &UpdateSubAccountRequest{Name: "a"})
			return err
		},
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "SubAccounts.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
