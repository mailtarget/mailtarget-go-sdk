package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The guard must reject the call locally: no request may leave the process.
func TestDo_WithoutSecretKey_MakesNoRequest(t *testing.T) {
	c := newTestClient(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts",
		httpmock.NewStringResponder(200, `{"data":[]}`))

	err := c.do(request{capability: "Contacts.List", method: http.MethodGet, path: "/contacts"})

	var configErr *ConfigError
	require.ErrorAs(t, err, &configErr)
	assert.Equal(t, "Contacts.List", configErr.Capability)
	assert.Contains(t, err.Error(), "Contacts.List")
	assert.Contains(t, err.Error(), "WithOpenAPISecretKey")
	assert.Equal(t, 0, httpmock.GetTotalCallCount(), "guard must fire before any HTTP call")
}

func TestDo_SignsRequest(t *testing.T) {
	c := withSecret(t)

	var gotAuth, gotAccept string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		func(req *http.Request) (*http.Response, error) {
			gotAuth = req.Header.Get("Authorization")
			gotAccept = req.Header.Get("Accept")
			return httpmock.NewStringResponse(200, `{"data":{}}`), nil
		})

	require.NoError(t, c.do(request{capability: "Things.Get", method: http.MethodGet, path: "/things"}))
	assert.Equal(t, "Bearer "+testSecretKey, gotAuth)
	assert.Equal(t, "application/json", gotAccept)
}

func TestDo_MapsErrorEnvelope(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		httpmock.NewStringResponder(403, `{"error":"forbidden","message":"missing contacts.read scope"}`))

	err := c.do(request{capability: "Things.Get", method: http.MethodGet, path: "/things"})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 403, apiErr.StatusCode)
	assert.Equal(t, "forbidden", apiErr.Code)
	assert.Equal(t, "missing contacts.read scope", apiErr.Message)
	assert.Contains(t, err.Error(), "http 403")
}

func TestDo_ErrorWithUnparsableBodyStillReports(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		httpmock.NewStringResponder(500, `<html>gateway blew up</html>`))

	err := c.do(request{capability: "Things.Get", method: http.MethodGet, path: "/things"})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 500, apiErr.StatusCode)
	assert.NotEmpty(t, apiErr.Message)
}

func TestDo_SendsQueryAndBody(t *testing.T) {
	c := withSecret(t)

	var gotQuery, gotBody, gotContentType string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/things",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.RawQuery
			gotContentType = req.Header.Get("Content-Type")
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{}}`), nil
		})

	err := c.do(request{
		capability: "Things.Create", method: http.MethodPost, path: "/things",
		query: map[string][]string{"page": {"3"}},
		body:  map[string]string{"email": "a@b.co"},
	})

	require.NoError(t, err)
	assert.Equal(t, "page=3", gotQuery)
	assert.Equal(t, "application/json", gotContentType)
	assert.JSONEq(t, `{"email":"a@b.co"}`, gotBody)
}

// The generic helpers below are what every resource method is built from.

func TestObject_DecodesData(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		httpmock.NewStringResponder(200, `{"data":{"id":"42"}}`))

	out, err := object[struct {
		ID string `json:"id"`
	}](c, request{capability: "Things.Get", method: http.MethodGet, path: "/things"})

	require.NoError(t, err)
	assert.Equal(t, "42", out.ID)
}

func TestObjects_DecodesDataSlice(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		httpmock.NewStringResponder(200, `{"data":[{"id":"1"},{"id":"2"}]}`))

	out, err := objects[struct {
		ID string `json:"id"`
	}](c, request{capability: "Things.List", method: http.MethodGet, path: "/things"})

	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "2", out[1].ID)
}

func TestPaged_DecodesDataAndMeta(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		httpmock.NewStringResponder(200,
			`{"data":[{"id":"1"}],"meta":{"page":2,"perPage":20,"total":41}}`))

	page, err := paged[struct {
		ID string `json:"id"`
	}](c, request{capability: "Things.List", method: http.MethodGet, path: "/things"})

	require.NoError(t, err)
	assert.Len(t, page.Items, 1)
	assert.Equal(t, Meta{Page: 2, PerPage: 20, Total: 41}, page.Meta)
}

// Endpoints that answer with a bare object instead of {"data": ...}.
func TestBare_DecodesWholeBody(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/contacts/count",
		httpmock.NewStringResponder(200, `{"count":17}`))

	out, err := bare[contactCount](c, request{
		capability: "Contacts.Count", method: http.MethodPost, path: "/contacts/count",
	})

	require.NoError(t, err)
	assert.Equal(t, 17, out.Count)
}

// dto.ContactImportErrorEnvelope adds an untyped `details` member, which is
// kept verbatim on the error.
func TestDo_KeepsErrorDetails(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/contacts/import",
		httpmock.NewStringResponder(409, `{"error":"IMPORT_IN_PROGRESS","message":"an import is running",
			"details":{"imports":[{"id":"imp1","state":"RUNNING"}]}}`))

	err := c.do(request{capability: "Contacts.Import", method: http.MethodPost, path: "/contacts/import"})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "IMPORT_IN_PROGRESS", apiErr.Code)
	assert.Equal(t, "an import is running", apiErr.Message)
	assert.JSONEq(t, `{"imports":[{"id":"imp1","state":"RUNNING"}]}`, string(apiErr.Details))
}

func TestDo_ErrorWithoutDetailsLeavesThemNil(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/things",
		httpmock.NewStringResponder(404, `{"error":"not_found","message":"nope","details":null}`))

	err := c.do(request{capability: "Things.Get", method: http.MethodGet, path: "/things"})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Nil(t, apiErr.Details)
}

// Every dto.BoolEnvelope operation must turn `data: false` into an
// OPERATION_FAILED error carrying the HTTP status, and accept `data: true`.
func TestBoolEnvelopeOperations_DataFalseFails(t *testing.T) {
	type op struct {
		method, path string
		call         func(c *Client) error
	}
	ops := map[string]op{
		"APIKeys.Delete":      {"DELETE", "/v1/api-keys/7", func(c *Client) error { return c.APIKeys.Delete(7) }},
		"Campaigns.Send":      {"POST", "/v1/campaigns/cmp1/send", func(c *Client) error { return c.Campaigns.Send("cmp1") }},
		"Campaigns.SendTest":  {"POST", "/v1/campaigns/cmp1/send-test", func(c *Client) error { return c.Campaigns.SendTest("cmp1", "a@b.co") }},
		"Senders.Delete":      {"DELETE", "/v1/domain/senders/snd1", func(c *Client) error { return c.Senders.Delete("snd1") }},
		"Forms.Delete":        {"DELETE", "/v1/forms/f1", func(c *Client) error { return c.Forms.Delete("f1") }},
		"Forms.Submit":        {"POST", "/v1/forms/f1/submissions", func(c *Client) error { return c.Forms.Submit("f1", &FormSubmissionRequest{Values: map[string]any{}}) }},
		"Pages.Delete":        {"DELETE", "/v1/pages/p1", func(c *Client) error { return c.Pages.Delete("p1") }},
		"Files.Delete":        {"DELETE", "/v1/files/f1", func(c *Client) error { return c.Files.Delete("f1") }},
		"Files.DeleteLabel":   {"DELETE", "/v1/files/labels/Promo", func(c *Client) error { return c.Files.DeleteLabel("Promo") }},
		"Labels.Delete":       {"DELETE", "/v1/labels/vip", func(c *Client) error { return c.Labels.Delete("vip") }},
		"Segments.Delete":     {"DELETE", "/v1/segments/sg1", func(c *Client) error { return c.Segments.Delete("sg1") }},
		"Suppressions.Delete": {"DELETE", "/v1/suppressions/11", func(c *Client) error { return c.Suppressions.Delete(11) }},
		"Webhooks.Delete":     {"DELETE", "/v1/webhooks/6", func(c *Client) error { return c.Webhooks.Delete(6) }},
	}
	for name, o := range ops {
		t.Run(name, func(t *testing.T) {
			c := withSecret(t)
			httpmock.RegisterResponder(o.method, testBaseURL+o.path,
				httpmock.NewStringResponder(200, `{"data":false}`))

			err := o.call(c)

			var apiErr *Error
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, ErrCodeOperationFailed, apiErr.Code)
			assert.Equal(t, 200, apiErr.StatusCode)
			assert.Equal(t, name+" returned data=false", apiErr.Message)

			httpmock.RegisterResponder(o.method, testBaseURL+o.path,
				httpmock.NewStringResponder(200, `{"data":true}`))
			assert.NoError(t, o.call(c))
		})
	}
}

// A body without `data` keeps the pre-existing lenient behavior.
func TestConfirm_MissingDataIsSuccess(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/things/1",
		httpmock.NewStringResponder(200, `{}`))

	assert.NoError(t, confirm(c, request{capability: "Things.Delete", method: http.MethodDelete, path: "/things/1"}))
}
