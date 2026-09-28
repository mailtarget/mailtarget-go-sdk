package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/labels",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"_id":"lbl1","name":"vip","contactCount":12}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	page, err := c.Labels.List(&ListLabelsParams{Page: 1, PerPage: 20, Search: "vip"})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=20&search=vip", gotQuery)
	require.Len(t, page.Items, 1)
	// The API returns the identifier as `_id`.
	assert.Equal(t, "lbl1", page.Items[0].ID)
	assert.Equal(t, 12, page.Items[0].ContactCount)
}

func TestLabelsService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/labels",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"_id":"lbl2","name":"newsletter"}}`), nil
		})

	label, err := c.Labels.Create("newsletter")

	require.NoError(t, err)
	assert.Equal(t, "newsletter", label.Name)
	assert.JSONEq(t, `{"name":"newsletter"}`, gotBody)
}

func TestLabelsService_Rename(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/labels/old%20name",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"_id":"lbl1","name":"new name"}}`), nil
		})

	label, err := c.Labels.Rename("old name", "new name")

	require.NoError(t, err)
	assert.Equal(t, "new name", label.Name)
	assert.JSONEq(t, `{"name":"new name"}`, gotBody)
}

// The endpoint takes the names as one comma separated path segment.
func TestLabelsService_DeleteJoinsNames(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/labels/vip,cold%20lead",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Labels.Delete("vip", "cold lead"))
}

func TestLabelsService_DeleteRejectsUnaddressableInput(t *testing.T) {
	c := withSecret(t)

	assert.ErrorContains(t, c.Labels.Delete(), "at least one label name")
	assert.ErrorContains(t, c.Labels.Delete("a,b"), "contains a comma")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestLabelsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List":   func() error { _, err := c.Labels.List(nil); return err },
		"Create": func() error { _, err := c.Labels.Create("x"); return err },
		"Rename": func() error { _, err := c.Labels.Rename("a", "b"); return err },
		"Delete": func() error { return c.Labels.Delete("x") },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Labels.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
