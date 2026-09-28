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

func TestSubAccountsService_GuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	_, err := c.SubAccounts.List(nil)

	var configErr *ConfigError
	assert.ErrorAs(t, err, &configErr)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
