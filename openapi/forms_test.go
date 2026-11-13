package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/forms",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":"f1","name":"Newsletter","labels":["vip"],"published":true,
					"submissionCount":42,"createdAt":1750000000000,"updatedAt":1760000000000}],
				"meta":{"page":1,"perPage":10,"total":1}
			}`), nil
		})

	page, err := c.Forms.List(&ListFormsParams{
		Page: 1, PerPage: 10, Search: "news",
		Status: []string{FormStatusPublished, FormStatusUsed},
		Labels: []string{"vip", "cold lead"},
	})

	require.NoError(t, err)
	assert.Equal(t,
		"labels=vip%2Ccold+lead&page=1&perPage=10&search=news&status=published%2Cused",
		gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "Newsletter", page.Items[0].Name)
	assert.True(t, page.Items[0].Published)
	assert.Equal(t, 42, page.Items[0].SubmissionCount)
	assert.Equal(t, int64(1760000000000), page.Items[0].UpdatedAt)
	assert.Equal(t, 1, page.Meta.Total)
}

func TestFormsService_ListWithoutParamsSendsNoQuery(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/forms",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.RawQuery
			return httpmock.NewStringResponse(200, `{"data":[],"meta":{"page":1,"perPage":10,"total":0}}`), nil
		})

	_, err := c.Forms.List(nil)

	require.NoError(t, err)
	assert.Empty(t, gotQuery)
}

func TestFormsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/forms/f1",
		httpmock.NewStringResponder(200, `{"data":{
			"id":"f1","name":"Newsletter","published":true,"saveToContact":true,
			"publicUrl":"https://forms.test/f1","embedHtml":"<iframe></iframe>",
			"fields":[{"id":"fl1","name":"interest","title":"Interest","type":"checkbox",
				"multiValue":true,"required":true,
				"options":[{"id":"o1","name":"Promo"},{"id":"o2","name":"Newsletter"}]}],
			"usedIn":[{"id":"p1","name":"Landing","type":"page","published":true,"url":"https://p.test"}]
		}}`))

	form, err := c.Forms.Get("f1")

	require.NoError(t, err)
	assert.Equal(t, "https://forms.test/f1", form.PublicURL)
	assert.Equal(t, "<iframe></iframe>", form.EmbedHTML)
	assert.True(t, form.SaveToContact)
	require.Len(t, form.Fields, 1)
	assert.Equal(t, "interest", form.Fields[0].Name)
	assert.True(t, form.Fields[0].MultiValue)
	assert.Equal(t, "Newsletter", form.Fields[0].Options[1].Name)
	require.Len(t, form.UsedIn, 1)
	assert.Equal(t, FormUsagePage, form.UsedIn[0].Type)
}

func TestFormsService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/forms/f1",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Forms.Delete("f1"))
}

func TestFormsService_DeleteDataFalseFails(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/forms/f1",
		httpmock.NewStringResponder(200, `{"data":false}`))

	err := c.Forms.Delete("f1")

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, ErrCodeOperationFailed, apiErr.Code)
	assert.Equal(t, 200, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Forms.Delete")
}

func TestFormsService_Duplicate(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/forms/f1/duplicate",
		httpmock.NewStringResponder(201, `{"data":{"id":"f2","name":"Newsletter (Duplicate)","published":false}}`))

	form, err := c.Forms.Duplicate("f1")

	require.NoError(t, err)
	assert.Equal(t, "f2", form.ID)
	assert.Equal(t, "Newsletter (Duplicate)", form.Name)
}

func TestFormsService_PublishAndUnpublish(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/forms/f1/publish",
		httpmock.NewStringResponder(200, `{"data":{"id":"f1","published":true,"publicUrl":"https://forms.test/f1"}}`))
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/forms/f1/unpublish",
		httpmock.NewStringResponder(200, `{"data":{"id":"f1","published":false}}`))

	published, err := c.Forms.Publish("f1")
	require.NoError(t, err)
	assert.True(t, published.Published)
	assert.Equal(t, "https://forms.test/f1", published.PublicURL)

	unpublished, err := c.Forms.Unpublish("f1")
	require.NoError(t, err)
	assert.False(t, unpublished.Published)
}

func TestFormsService_Submissions(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/forms/f1/submissions",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[
					{"id":"s1","values":{"email":"a@b.co","interest":"Promo, Newsletter"},
						"removedFields":{"Old":"x"},"checkedIn":true,"submittedAt":1760000000123},
					{"id":"s2","values":{"email":"c@d.co"},"submittedAt":1760000000456}
				],
				"meta":{"page":2,"perPage":5,"total":7}
			}`), nil
		})

	page, err := c.Forms.Submissions("f1", &ListFormSubmissionsParams{
		Page: 2, PerPage: 5, Search: "a@b", From: 1750000000000, To: 1760000000000,
	})

	require.NoError(t, err)
	assert.Equal(t, "from=1750000000000&page=2&perPage=5&search=a%40b&to=1760000000000", gotQuery)
	require.Len(t, page.Items, 2)
	assert.Equal(t, "Promo, Newsletter", page.Items[0].Values["interest"])
	assert.Equal(t, "x", page.Items[0].RemovedFields["Old"])
	require.NotNil(t, page.Items[0].CheckedIn)
	assert.True(t, *page.Items[0].CheckedIn)
	assert.Equal(t, int64(1760000000123), page.Items[0].SubmittedAt)
	// checkedIn is only present for event check-in forms.
	assert.Nil(t, page.Items[1].CheckedIn)
	assert.Equal(t, 7, page.Meta.Total)
}

// The API requires from and to together, so the SDK refuses half a range.
func TestFormsService_SubmissionsRequiresBothBounds(t *testing.T) {
	c := withSecret(t)

	_, errFrom := c.Forms.Submissions("f1", &ListFormSubmissionsParams{From: 1})
	_, errTo := c.Forms.Submissions("f1", &ListFormSubmissionsParams{To: 1})

	assert.ErrorContains(t, errFrom, "both From and To")
	assert.ErrorContains(t, errTo, "both From and To")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestFormsService_Submit(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/forms/f1/submissions",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":true}`), nil
		})

	err := c.Forms.Submit("f1", &FormSubmissionRequest{Values: map[string]any{
		"email":    "sarah@example.com",
		"interest": []string{"Promo", "Newsletter"},
	}})

	require.NoError(t, err)
	assert.JSONEq(t,
		`{"values":{"email":"sarah@example.com","interest":["Promo","Newsletter"]}}`,
		gotBody)
}

func TestFormsService_SubmitDataFalseFails(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/forms/f1/submissions",
		httpmock.NewStringResponder(201, `{"data":false}`))

	err := c.Forms.Submit("f1", &FormSubmissionRequest{Values: map[string]any{"email": "a@b.co"}})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, ErrCodeOperationFailed, apiErr.Code)
	assert.Equal(t, 201, apiErr.StatusCode)
}

func TestFormsService_SubmitValidates(t *testing.T) {
	c := withSecret(t)

	errNil := c.Forms.Submit("f1", nil)
	errValues := c.Forms.Submit("f1", &FormSubmissionRequest{})

	assert.ErrorContains(t, errNil, "requires Values")
	assert.ErrorContains(t, errValues, "requires Values")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestFormsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)
	submit := &FormSubmissionRequest{Values: map[string]any{"email": "a@b.co"}}

	calls := map[string]func() error{
		"List":        func() error { _, err := c.Forms.List(nil); return err },
		"Get":         func() error { _, err := c.Forms.Get("f1"); return err },
		"Delete":      func() error { return c.Forms.Delete("f1") },
		"Duplicate":   func() error { _, err := c.Forms.Duplicate("f1"); return err },
		"Publish":     func() error { _, err := c.Forms.Publish("f1"); return err },
		"Unpublish":   func() error { _, err := c.Forms.Unpublish("f1"); return err },
		"Submissions": func() error { _, err := c.Forms.Submissions("f1", nil); return err },
		"Submit":      func() error { return c.Forms.Submit("f1", submit) },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Forms.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
