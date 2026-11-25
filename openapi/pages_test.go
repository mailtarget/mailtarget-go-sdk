package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPagesService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/pages",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":"p1","name":"Promo","category":"landing","labels":["vip"],
					"containsForm":true,"published":true,"default":true,
					"thumbnail":"https://cdn.test/p1.png","url":"https://pages.test/promo",
					"visitCount":5000000000,"uniqueVisitCount":42,
					"createdAt":1750000000000,"updatedAt":1760000000000}],
				"meta":{"page":1,"perPage":10,"total":1}
			}`), nil
		})

	page, err := c.Pages.List(&ListPagesParams{
		Page: 1, PerPage: 10, Search: "promo",
		Status:   []string{PageStatusPublished, PageStatusDefault},
		Category: PageCategoryLanding,
		Sort:     PageSortNameDesc,
	})

	require.NoError(t, err)
	assert.Equal(t,
		"category=landing&page=1&perPage=10&search=promo&sort=-name&status=published%2Cdefault",
		gotQuery)
	require.Len(t, page.Items, 1)
	p := page.Items[0]
	assert.Equal(t, "Promo", p.Name)
	assert.Equal(t, PageCategoryLanding, p.Category)
	assert.Equal(t, []string{"vip"}, p.Labels)
	assert.True(t, p.ContainsForm)
	assert.True(t, p.Published)
	assert.True(t, p.IsDefault)
	assert.Equal(t, "https://cdn.test/p1.png", p.Thumbnail)
	assert.Equal(t, "https://pages.test/promo", p.URL)
	assert.Equal(t, int64(5000000000), p.VisitCount)
	assert.Equal(t, int64(42), p.UniqueVisitCount)
	assert.Equal(t, int64(1750000000000), p.CreatedAt)
	assert.Equal(t, int64(1760000000000), p.UpdatedAt)
	assert.Equal(t, 1, page.Meta.Total)
}

func TestPagesService_ListWithoutParamsSendsNoQuery(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/pages",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.RawQuery
			return httpmock.NewStringResponse(200, `{"data":[],"meta":{"page":1,"perPage":10,"total":0}}`), nil
		})

	_, err := c.Pages.List(nil)

	require.NoError(t, err)
	assert.Empty(t, gotQuery)
}

func TestPagesService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/pages/p1",
		httpmock.NewStringResponder(200, `{"data":{
			"id":"p1","name":"Thanks","category":"thank_you","published":false,"default":false,
			"containsForm":true,"url":"https://pages.test/thanks","permalink":"thanks",
			"forms":[{"id":"f1","name":"Newsletter"}],"updatedAt":1760000000000
		}}`))

	page, err := c.Pages.Get("p1")

	require.NoError(t, err)
	assert.Equal(t, PageCategoryThankYou, page.Category)
	assert.False(t, page.Published)
	assert.False(t, page.IsDefault)
	assert.Equal(t, "https://pages.test/thanks", page.URL)
	assert.Equal(t, "thanks", page.Permalink)
	require.Len(t, page.Forms, 1)
	assert.Equal(t, LandingPageFormRef{ID: "f1", Name: "Newsletter"}, page.Forms[0])
	assert.Equal(t, int64(1760000000000), page.UpdatedAt)
}

func TestPagesService_GetEscapesID(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/pages/a%2Fb",
		httpmock.NewStringResponder(200, `{"data":{"id":"a/b"}}`))

	page, err := c.Pages.Get("a/b")

	require.NoError(t, err)
	assert.Equal(t, "a/b", page.ID)
}

func TestPagesService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/pages/p1",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Pages.Delete("p1"))
}

func TestPagesService_DeleteDataFalseFails(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/pages/p1",
		httpmock.NewStringResponder(200, `{"data":false}`))

	err := c.Pages.Delete("p1")

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, ErrCodeOperationFailed, apiErr.Code)
	assert.Equal(t, 200, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Pages.Delete")
}

func TestPagesService_DeleteDefaultPageConflict(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/pages/p1",
		httpmock.NewStringResponder(409, `{"error":"CONFLICT","message":"default page cannot be deleted"}`))

	err := c.Pages.Delete("p1")

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 409, apiErr.StatusCode)
	assert.Equal(t, "CONFLICT", apiErr.Code)
}

func TestPagesService_Duplicate(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/pages/p1/duplicate",
		httpmock.NewStringResponder(201, `{"data":{"id":"p2","name":"Promo (copy)","published":false}}`))

	page, err := c.Pages.Duplicate("p1")

	require.NoError(t, err)
	assert.Equal(t, "p2", page.ID)
	assert.Equal(t, "Promo (copy)", page.Name)
	assert.False(t, page.Published)
}

// This endpoint returns {"data": [...]} with no meta, so it is not paginated.
func TestPagesService_Forms(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/pages/p1/forms",
		httpmock.NewStringResponder(200, `{"data":[
			{"id":"f1","name":"Newsletter","labels":["vip"],"submissionCount":12,
				"createdAt":1750000000000,"updatedAt":1760000000000},
			{"id":"f2","name":"Contact"}
		]}`))

	forms, err := c.Pages.Forms("p1")

	require.NoError(t, err)
	require.Len(t, forms, 2)
	assert.Equal(t, "Newsletter", forms[0].Name)
	assert.Equal(t, []string{"vip"}, forms[0].Labels)
	assert.Equal(t, int64(12), forms[0].SubmissionCount)
	assert.Equal(t, int64(1750000000000), forms[0].CreatedAt)
	assert.Equal(t, int64(1760000000000), forms[0].UpdatedAt)
	assert.Equal(t, "f2", forms[1].ID)
}

func TestPagesService_PublishUnpublishAndSetDefault(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/pages/p1/publish",
		httpmock.NewStringResponder(200, `{"data":{"id":"p1","published":true,"url":"https://pages.test/promo"}}`))
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/pages/p1/unpublish",
		httpmock.NewStringResponder(200, `{"data":{"id":"p1","published":false}}`))
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/pages/p1/set-default",
		httpmock.NewStringResponder(200, `{"data":{"id":"p1","published":true,"default":true}}`))

	published, err := c.Pages.Publish("p1")
	require.NoError(t, err)
	assert.True(t, published.Published)
	assert.Equal(t, "https://pages.test/promo", published.URL)

	unpublished, err := c.Pages.Unpublish("p1")
	require.NoError(t, err)
	assert.False(t, unpublished.Published)

	def, err := c.Pages.SetDefault("p1")
	require.NoError(t, err)
	assert.True(t, def.IsDefault)
	assert.True(t, def.Published)
}

func TestPagesService_PublishQuotaExceeded(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/pages/p1/publish",
		httpmock.NewStringResponder(422, `{"error":"QUOTA_EXCEEDED","message":"published page quota used up"}`))

	_, err := c.Pages.Publish("p1")

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 422, apiErr.StatusCode)
	assert.Equal(t, "QUOTA_EXCEEDED", apiErr.Code)
}

func TestPagesService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List":       func() error { _, err := c.Pages.List(nil); return err },
		"Get":        func() error { _, err := c.Pages.Get("p1"); return err },
		"Delete":     func() error { return c.Pages.Delete("p1") },
		"Duplicate":  func() error { _, err := c.Pages.Duplicate("p1"); return err },
		"Forms":      func() error { _, err := c.Pages.Forms("p1"); return err },
		"Publish":    func() error { _, err := c.Pages.Publish("p1"); return err },
		"Unpublish":  func() error { _, err := c.Pages.Unpublish("p1"); return err },
		"SetDefault": func() error { _, err := c.Pages.SetDefault("p1"); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Pages.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
