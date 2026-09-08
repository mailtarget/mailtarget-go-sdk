package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":"1","email":"a@b.co","name":"A","labels":["vip"]}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	page, err := c.Contacts.List(&ListContactsParams{Page: 1, PerPage: 20, Search: "a@b.co"})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=20&search=a%40b.co", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "a@b.co", page.Items[0].Email)
	assert.Equal(t, []string{"vip"}, page.Items[0].Labels)
	assert.Equal(t, 1, page.Meta.Total)
}

func TestContactsService_ListWithoutParamsSendsNoQuery(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.RawQuery
			return httpmock.NewStringResponse(200, `{"data":[]}`), nil
		})

	page, err := c.Contacts.List(nil)

	require.NoError(t, err)
	assert.Empty(t, gotQuery)
	assert.Empty(t, page.Items)
}

func TestContactsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/abc123",
		httpmock.NewStringResponder(200, `{"data":{"id":"abc123","email":"a@b.co"}}`))

	contact, err := c.Contacts.Get("abc123")

	require.NoError(t, err)
	assert.Equal(t, "abc123", contact.ID)
	assert.Equal(t, "a@b.co", contact.Email)
}

// Path parameters must be escaped so they cannot break out of their segment.
func TestContactsService_GetEscapesPathParameter(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/a%2Fb",
		httpmock.NewStringResponder(200, `{"data":{"id":"a/b"}}`))

	contact, err := c.Contacts.Get("a/b")

	require.NoError(t, err)
	assert.Equal(t, "a/b", contact.ID)
}

func TestContactsService_GetByEmail(t *testing.T) {
	c := withSecret(t)
	// A plus sign is legal inside a path segment, so it travels unescaped.
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/email/a+plus@b.co",
		httpmock.NewStringResponder(200, `{"data":{"id":"1","email":"a+plus@b.co"}}`))

	contact, err := c.Contacts.GetByEmail("a+plus@b.co")

	require.NoError(t, err)
	assert.Equal(t, "a+plus@b.co", contact.Email)
}

func TestContactsService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/contacts",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":"new1","email":"a@b.co"}}`), nil
		})

	contact, err := c.Contacts.Create(&CreateContactRequest{Email: "a@b.co", Firstname: "A"})

	require.NoError(t, err)
	assert.Equal(t, "new1", contact.ID)
	// Unset optional fields must be omitted rather than sent as empty strings.
	assert.JSONEq(t, `{"email":"a@b.co","firstname":"A"}`, gotBody)
}

func TestContactsService_Update(t *testing.T) {
	c := withSecret(t)

	var gotMethod, gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/contacts/abc123",
		func(req *http.Request) (*http.Response, error) {
			gotMethod = req.Method
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"id":"abc123","note":"vip"}}`), nil
		})

	contact, err := c.Contacts.Update("abc123", &UpdateContactRequest{Note: "vip"})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.JSONEq(t, `{"note":"vip"}`, gotBody)
	assert.Equal(t, "vip", contact.Note)
}

func TestContactsService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/contacts/abc123",
		httpmock.NewStringResponder(200, `{"message":"contact deleted"}`))

	assert.NoError(t, c.Contacts.Delete("abc123"))
}

func TestContactsService_Count(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/contacts/count",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"count":42}`), nil
		})

	count, err := c.Contacts.Count(nil)

	require.NoError(t, err)
	assert.Equal(t, 42, count)
	assert.JSONEq(t, `{}`, gotBody)
}

func TestContactsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List":       func() error { _, err := c.Contacts.List(nil); return err },
		"Get":        func() error { _, err := c.Contacts.Get("1"); return err },
		"GetByEmail": func() error { _, err := c.Contacts.GetByEmail("a@b.co"); return err },
		"Create":     func() error { _, err := c.Contacts.Create(&CreateContactRequest{}); return err },
		"Update":     func() error { _, err := c.Contacts.Update("1", &UpdateContactRequest{}); return err },
		"Delete":     func() error { return c.Contacts.Delete("1") },
		"Count":      func() error { _, err := c.Contacts.Count(nil); return err },
	}

	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Contacts.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount(), "no guarded call may reach the network")
}
