package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegmentsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/segments",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":"sg1","name":"Vip Buyers","count":12,
					"filters":[{"field":"labels","operator":"in","value":["vip"]}],
					"createdAt":1750000000000}],
				"meta":{"page":1,"perPage":10,"total":1}
			}`), nil
		})

	page, err := c.Segments.List(&ListSegmentsParams{Page: 1, PerPage: 10, Search: "vip", Sort: "-count"})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=10&search=vip&sort=-count", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "Vip Buyers", page.Items[0].Name)
	assert.Equal(t, 12, page.Items[0].Count)
	assert.Equal(t, "labels", page.Items[0].Filters[0]["field"])
	assert.Equal(t, 1, page.Meta.Total)
}

func TestSegmentsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/segments/sg1",
		httpmock.NewStringResponder(200, `{"data":{"id":"sg1","name":"Everyone","filters":[]}}`))

	segment, err := c.Segments.Get("sg1")

	require.NoError(t, err)
	assert.Equal(t, "sg1", segment.ID)
	assert.Empty(t, segment.Filters)
}

func TestSegmentsService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/segments",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":"sg9","name":"Vip"}}`), nil
		})

	segment, err := c.Segments.Create(&CreateSegmentRequest{
		Name:    "vip",
		Filters: []map[string]any{{"field": "labels", "operator": "in", "value": []string{"vip"}}},
	})

	require.NoError(t, err)
	assert.Equal(t, "sg9", segment.ID)
	assert.JSONEq(t,
		`{"name":"vip","filters":[{"field":"labels","operator":"in","value":["vip"]}]}`,
		gotBody)
}

// An empty, non-nil Filters is a deliberate "match every contact" and must be
// sent as [], not dropped.
func TestSegmentsService_CreateMatchAll(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/segments",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":"sg9"}}`), nil
		})

	_, err := c.Segments.Create(&CreateSegmentRequest{Name: "everyone", Filters: []map[string]any{}})

	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"everyone","filters":[]}`, gotBody)
}

func TestSegmentsService_CreateValidates(t *testing.T) {
	c := withSecret(t)

	_, errNil := c.Segments.Create(nil)
	_, errName := c.Segments.Create(&CreateSegmentRequest{Filters: []map[string]any{}})
	_, errFilters := c.Segments.Create(&CreateSegmentRequest{Name: "vip"})

	assert.ErrorContains(t, errNil, "requires a Name")
	assert.ErrorContains(t, errName, "requires a Name")
	assert.ErrorContains(t, errFilters, "requires Filters")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSegmentsService_Update(t *testing.T) {
	c := withSecret(t)

	var bodies []string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/segments/sg1",
		func(req *http.Request) (*http.Response, error) {
			bodies = append(bodies, readBody(t, req))
			return httpmock.NewStringResponse(200, `{"data":{"id":"sg1","name":"Renamed"}}`), nil
		})

	segment, err := c.Segments.Update("sg1", &UpdateSegmentRequest{Name: "renamed"})
	require.NoError(t, err)
	assert.Equal(t, "Renamed", segment.Name)

	_, err = c.Segments.Update("sg1", &UpdateSegmentRequest{Filters: []map[string]any{}})
	require.NoError(t, err)

	require.Len(t, bodies, 2)
	// A nil Filters keeps the current filters, so it must not be sent.
	assert.JSONEq(t, `{"name":"renamed"}`, bodies[0])
	// An empty Filters replaces them with "match every contact".
	assert.JSONEq(t, `{"filters":[]}`, bodies[1])
}

func TestSegmentsService_UpdateRequiresAField(t *testing.T) {
	c := withSecret(t)

	_, errNil := c.Segments.Update("sg1", nil)
	_, errEmpty := c.Segments.Update("sg1", &UpdateSegmentRequest{})

	assert.ErrorContains(t, errNil, "at least one of Name or Filters")
	assert.ErrorContains(t, errEmpty, "at least one of Name or Filters")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSegmentsService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/segments/sg1",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Segments.Delete("sg1"))
}

func TestSegmentsService_RecipientsCount(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/segments/sg1/recipients-count",
		httpmock.NewStringResponder(200, `{"data":{"count":321}}`))

	count, err := c.Segments.RecipientsCount("sg1")

	require.NoError(t, err)
	assert.Equal(t, 321, count)
}

func TestSegmentsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)
	create := &CreateSegmentRequest{Name: "vip", Filters: []map[string]any{}}
	update := &UpdateSegmentRequest{Name: "vip"}

	calls := map[string]func() error{
		"List":            func() error { _, err := c.Segments.List(nil); return err },
		"Get":             func() error { _, err := c.Segments.Get("1"); return err },
		"Create":          func() error { _, err := c.Segments.Create(create); return err },
		"Update":          func() error { _, err := c.Segments.Update("1", update); return err },
		"Delete":          func() error { return c.Segments.Delete("1") },
		"RecipientsCount": func() error { _, err := c.Segments.RecipientsCount("1"); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Segments.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
