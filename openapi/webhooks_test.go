package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhooksService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/webhooks",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":5,"name":"crm","targetUrl":"https://crm.example.com/hook",
					"isActive":false,"lastFailure":1750000000000}],
				"meta":{"page":1,"perPage":10,"total":1}
			}`), nil
		})

	active := false
	page, err := c.Webhooks.List(&ListWebhooksParams{Page: 1, PerPage: 10, Search: "crm", IsActive: &active})

	require.NoError(t, err)
	// A pointer distinguishes "filter by false" from "do not filter".
	assert.Equal(t, "isActive=false&page=1&perPage=10&search=crm", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "https://crm.example.com/hook", page.Items[0].TargetURL)
	assert.Equal(t, int64(1750000000000), page.Items[0].LastFailure)
}

func TestWebhooksService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/webhooks/5",
		httpmock.NewStringResponder(200, `{"data":{
			"id":5,"name":"crm","targetUrl":"https://crm.example.com/hook","isActive":true,
			"authenticationId":2,"username":"bot",
			"events":[{"id":1,"name":"Delivered","slug":"delivered"}]
		}}`))

	hook, err := c.Webhooks.Get(5)

	require.NoError(t, err)
	assert.Equal(t, WebhookAuthBasic, hook.AuthenticationID)
	assert.Equal(t, "bot", hook.Username)
	require.Len(t, hook.Events, 1)
	assert.Equal(t, "delivered", hook.Events[0].Slug)
}

func TestWebhooksService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/webhooks",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":6,"name":"crm","isActive":true}}`), nil
		})

	hook, err := c.Webhooks.Create(&CreateWebhookRequest{
		Name:             "crm",
		TargetURL:        "https://crm.example.com/hook",
		EventIDs:         []int{1, 2},
		AuthenticationID: WebhookAuthBasic,
		Username:         "bot",
		Password:         "s3cret",
	})

	require.NoError(t, err)
	assert.Equal(t, 6, hook.ID)
	assert.JSONEq(t, `{
		"name":"crm","targetUrl":"https://crm.example.com/hook","eventIds":[1,2],
		"authenticationId":2,"username":"bot","password":"s3cret"
	}`, gotBody)
}

func TestWebhooksService_Update(t *testing.T) {
	c := withSecret(t)

	var gotMethod, gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/webhooks/6",
		func(req *http.Request) (*http.Response, error) {
			gotMethod, gotBody = req.Method, readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"id":6,"isActive":false}}`), nil
		})

	inactive := false
	hook, err := c.Webhooks.Update(6, &UpdateWebhookRequest{
		Name: "crm", TargetURL: "https://crm.example.com/hook", EventIDs: []int{1}, IsActive: &inactive,
	})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	// A nil Password keeps the stored one, so it must not be sent; a false
	// IsActive is a real change and must be.
	assert.JSONEq(t, `{
		"name":"crm","targetUrl":"https://crm.example.com/hook","eventIds":[1],"isActive":false
	}`, gotBody)
	assert.False(t, hook.IsActive)
}

func TestWebhooksService_ValidatesRequests(t *testing.T) {
	c := withSecret(t)

	_, errCreateNil := c.Webhooks.Create(nil)
	_, errCreate := c.Webhooks.Create(&CreateWebhookRequest{Name: "crm", TargetURL: "https://x.co"})
	_, errUpdateNil := c.Webhooks.Update(6, nil)
	_, errUpdate := c.Webhooks.Update(6, &UpdateWebhookRequest{TargetURL: "https://x.co", EventIDs: []int{1}})

	assert.ErrorContains(t, errCreateNil, "requires a request")
	assert.ErrorContains(t, errCreate, "Name, TargetURL and at least one EventID")
	assert.ErrorContains(t, errUpdateNil, "requires a request")
	assert.ErrorContains(t, errUpdate, "Name, TargetURL and at least one EventID")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestWebhooksService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/webhooks/6",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Webhooks.Delete(6))
}

func TestWebhooksService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)
	create := &CreateWebhookRequest{Name: "a", TargetURL: "https://x.co", EventIDs: []int{1}}
	update := &UpdateWebhookRequest{Name: "a", TargetURL: "https://x.co", EventIDs: []int{1}}

	calls := map[string]func() error{
		"List":   func() error { _, err := c.Webhooks.List(nil); return err },
		"Get":    func() error { _, err := c.Webhooks.Get(1); return err },
		"Create": func() error { _, err := c.Webhooks.Create(create); return err },
		"Update": func() error { _, err := c.Webhooks.Update(1, update); return err },
		"Delete": func() error { return c.Webhooks.Delete(1) },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Webhooks.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
