package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const suppressionPage = `{
	"data":[{"id":4,"email":"a@b.co","type":"Non-transactional","source":"Bounce Rule",
		"subAccountId":3,"subAccountName":"marketing","createdAt":1750000000000}],
	"meta":{"page":1,"perPage":10,"total":1}
}`

func TestSuppressionsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/suppressions",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, suppressionPage), nil
		})

	page, err := c.Suppressions.List(&ListSuppressionsParams{
		Page: 1, PerPage: 10, Email: "a@b.co", Type: "Non-transactional",
		Source: "Manual,Bounce Rule", SubAccountID: 3,
	})

	require.NoError(t, err)
	assert.Equal(t,
		"email=a%40b.co&page=1&perPage=10&source=Manual%2CBounce+Rule&subAccountId=3&type=Non-transactional",
		gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "marketing", page.Items[0].SubAccountName)
	assert.Equal(t, int64(1750000000000), page.Items[0].CreatedAt)
}

// Bounces and Unsubscribes preset the source server side, so a Source set on
// the shared params must not leak into their query.
func TestSuppressionsService_PresetSourceEndpoints(t *testing.T) {
	c := withSecret(t)

	queries := map[string]string{}
	for _, path := range []string{"/suppressions/bounces", "/suppressions/unsubscribes"} {
		path := path
		httpmock.RegisterResponder("GET", testBaseURL+"/v1"+path,
			func(req *http.Request) (*http.Response, error) {
				queries[path] = req.URL.Query().Encode()
				return httpmock.NewStringResponse(200, suppressionPage), nil
			})
	}
	params := &ListSuppressionsParams{Page: 2, Source: "Manual"}

	bounces, err := c.Suppressions.Bounces(params)
	require.NoError(t, err)
	unsubscribes, err := c.Suppressions.Unsubscribes(params)
	require.NoError(t, err)

	assert.Equal(t, "page=2", queries["/suppressions/bounces"])
	assert.Equal(t, "page=2", queries["/suppressions/unsubscribes"])
	assert.Len(t, bounces.Items, 1)
	assert.Len(t, unsubscribes.Items, 1)
}

func TestSuppressionsService_Lookup(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/suppressions/lookup",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, suppressionPage), nil
		})

	page, err := c.Suppressions.Lookup(&LookupSuppressionsParams{Email: "a@b.co", SubAccountID: 3})

	require.NoError(t, err)
	assert.Equal(t, "email=a%40b.co&subAccountId=3", gotQuery)
	assert.Len(t, page.Items, 1)
}

func TestSuppressionsService_LookupRequiresEmail(t *testing.T) {
	c := withSecret(t)

	_, errNil := c.Suppressions.Lookup(nil)
	_, errEmpty := c.Suppressions.Lookup(&LookupSuppressionsParams{SubAccountID: 3})

	assert.ErrorContains(t, errNil, "requires an Email")
	assert.ErrorContains(t, errEmpty, "requires an Email")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSuppressionsService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/suppressions",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201,
				`{"data":{"id":11,"email":"a@b.co","type":"Transactional","source":"Manual"}}`), nil
		})

	sup, err := c.Suppressions.Create(&CreateSuppressionRequest{
		Email: "a@b.co", Type: "Transactional", Source: "Manual",
	})

	require.NoError(t, err)
	// An unset sub-account is omitted, making the suppression account-wide.
	assert.JSONEq(t, `{"email":"a@b.co","type":"Transactional","source":"Manual"}`, gotBody)
	assert.Equal(t, 11, sup.ID)
	assert.Equal(t, 0, sup.SubAccountID)
}

func TestSuppressionsService_CreateRequiresFields(t *testing.T) {
	c := withSecret(t)

	requests := map[string]*CreateSuppressionRequest{
		"nil":       nil,
		"no email":  {Type: "Transactional", Source: "Manual"},
		"no type":   {Email: "a@b.co", Source: "Manual"},
		"no source": {Email: "a@b.co", Type: "Transactional"},
	}
	for name, req := range requests {
		_, err := c.Suppressions.Create(req)
		assert.ErrorContains(t, err, "Email, Type and Source", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSuppressionsService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/suppressions/11",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Suppressions.Delete(11))
}

func TestSuppressionsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)
	create := &CreateSuppressionRequest{Email: "a@b.co", Type: "Transactional", Source: "Manual"}

	calls := map[string]func() error{
		"List":         func() error { _, err := c.Suppressions.List(nil); return err },
		"Bounces":      func() error { _, err := c.Suppressions.Bounces(nil); return err },
		"Unsubscribes": func() error { _, err := c.Suppressions.Unsubscribes(nil); return err },
		"Lookup": func() error {
			_, err := c.Suppressions.Lookup(&LookupSuppressionsParams{Email: "a@b.co"})
			return err
		},
		"Create": func() error { _, err := c.Suppressions.Create(create); return err },
		"Delete": func() error { return c.Suppressions.Delete(1) },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Suppressions.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
