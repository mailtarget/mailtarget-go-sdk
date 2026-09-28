package openapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func importRequest() *ImportContactsRequest {
	return &ImportContactsRequest{
		File:         strings.NewReader("Email Address,Name\na@b.co,A\n"),
		Filename:     "contacts.csv",
		FieldMapping: map[string]string{"Email Address": "email", "Name": "firstname"},
		Labels:       []string{"vip", "newsletter"},
	}
}

func TestContactsService_Fields(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/fields",
		httpmock.NewStringResponder(200, `{"data":[
			{"fieldName":"email","displayType":"EMAIL","fieldGroup":"DEFAULT_MAIN"},
			{"fieldName":"tier","displayType":"CHOICE","fieldGroup":"CUSTOM","values":["gold","silver"]}
		]}`))

	fields, err := c.Contacts.Fields()

	require.NoError(t, err)
	require.Len(t, fields, 2)
	assert.Equal(t, "email", fields[0].FieldName)
	assert.Equal(t, "CUSTOM", fields[1].FieldGroup)
	assert.Equal(t, []any{"gold", "silver"}, fields[1].Values)
}

// The import is the only multipart endpoint: the file and the form fields
// must travel as separate parts, with the usual auth header.
func TestContactsService_ImportSendsMultipart(t *testing.T) {
	c := withSecret(t)

	var gotAuth, gotFile, gotFilename string
	var gotForm map[string][]string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/contacts/import",
		func(req *http.Request) (*http.Response, error) {
			gotAuth = req.Header.Get("Authorization")
			require.NoError(t, req.ParseMultipartForm(1<<20))
			gotForm = req.MultipartForm.Value
			file, header, err := req.FormFile("file")
			require.NoError(t, err)
			content, _ := io.ReadAll(file)
			gotFile, gotFilename = string(content), header.Filename
			return httpmock.NewStringResponse(202, `{
				"data":{"id":"imp1","state":"READY","count":1,"labels":["vip","newsletter"]},
				"validation":{"totalRows":1,"estimatedValid":1,"issues":[],"issuesTruncated":false}
			}`), nil
		})

	req := importRequest()
	req.ImportMode = ImportModeSkipExisting
	req.Strict = true
	result, err := c.Contacts.Import(req)

	require.NoError(t, err)
	assert.Equal(t, "Bearer "+testSecretKey, gotAuth)
	assert.Equal(t, "contacts.csv", gotFilename)
	assert.Equal(t, "Email Address,Name\na@b.co,A\n", gotFile)
	assert.Equal(t, []string{"vip,newsletter"}, gotForm["labels"])
	assert.Equal(t, []string{"SKIP_EXISTING"}, gotForm["importMode"])
	assert.Equal(t, []string{"true"}, gotForm["strict"])

	var mapping map[string]string
	require.NoError(t, json.Unmarshal([]byte(gotForm["fieldMapping"][0]), &mapping))
	assert.Equal(t, map[string]string{"Email Address": "email", "Name": "firstname"}, mapping)

	// The 202 body carries validation next to data; both are returned.
	assert.Equal(t, "imp1", result.Import.ID)
	assert.Equal(t, "READY", result.Import.State)
	require.NotNil(t, result.Validation)
	assert.Equal(t, 1, result.Validation.EstimatedValid)
}

func TestContactsService_ImportOmitsUnsetOptionalFields(t *testing.T) {
	c := withSecret(t)

	var gotForm map[string][]string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/contacts/import",
		func(req *http.Request) (*http.Response, error) {
			require.NoError(t, req.ParseMultipartForm(1<<20))
			gotForm = req.MultipartForm.Value
			return httpmock.NewStringResponse(202, `{"data":{"id":"imp1"}}`), nil
		})

	_, err := c.Contacts.Import(importRequest())

	require.NoError(t, err)
	assert.NotContains(t, gotForm, "importMode")
	assert.NotContains(t, gotForm, "strict")
}

func TestContactsService_ImportRequiresMandatoryFields(t *testing.T) {
	c := withSecret(t)

	noFile := importRequest()
	noFile.File = nil
	noName := importRequest()
	noName.Filename = " "
	noMapping := importRequest()
	noMapping.FieldMapping = nil
	noLabels := importRequest()
	noLabels.Labels = nil
	commaLabel := importRequest()
	commaLabel.Labels = []string{"a,b"}

	for name, tc := range map[string]struct {
		req  *ImportContactsRequest
		want string
	}{
		"nil request": {nil, "File and a Filename"},
		"no file":     {noFile, "File and a Filename"},
		"no filename": {noName, "File and a Filename"},
		"no mapping":  {noMapping, "FieldMapping"},
		"no labels":   {noLabels, "at least one label"},
		"comma label": {commaLabel, "contains a comma"},
	} {
		_, err := c.Contacts.Import(tc.req)
		assert.ErrorContains(t, err, tc.want, name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount(), "invalid imports must not reach the network")
}

func TestContactsService_Imports(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/import",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":"imp1","state":"FINISH","newContact":5}],
				"meta":{"page":1,"perPage":10,"total":1}
			}`), nil
		})

	page, err := c.Contacts.Imports(&ListContactImportsParams{Page: 1, PerPage: 10, State: "FINISH", Sort: "createdAt"})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=10&sort=createdAt&state=FINISH", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, 5, page.Items[0].NewContact)
	assert.Equal(t, 1, page.Meta.Total)
}

func TestContactsService_CurrentImport(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/import/current",
		httpmock.NewStringResponder(200,
			`{"data":{"stillImporting":true,"imports":[{"id":"imp1","state":"RUNNING"}]}}`))

	queue, err := c.Contacts.CurrentImport()

	require.NoError(t, err)
	assert.True(t, queue.StillImporting)
	require.Len(t, queue.Imports, 1)
	assert.Equal(t, "RUNNING", queue.Imports[0].State)
}

func TestContactsService_GetImport(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/contacts/import/imp1",
		httpmock.NewStringResponder(200,
			`{"data":{"id":"imp1","state":"INTERUPTED","message":"bad header"}}`))

	job, err := c.Contacts.GetImport("imp1")

	require.NoError(t, err)
	assert.Equal(t, "INTERUPTED", job.State)
	assert.Equal(t, "bad header", job.Message)
}

func TestContactsService_CancelImport(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/contacts/import/imp1",
		httpmock.NewStringResponder(200, `{"data":{"id":"imp1","cancelled":true,"state":"DELETED"}}`))

	out, err := c.Contacts.CancelImport("imp1")

	require.NoError(t, err)
	assert.True(t, out.Cancelled)
	assert.Equal(t, "DELETED", out.State)
}

func TestContactsService_ImportEndpointsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"Fields":        func() error { _, err := c.Contacts.Fields(); return err },
		"Import":        func() error { _, err := c.Contacts.Import(importRequest()); return err },
		"Imports":       func() error { _, err := c.Contacts.Imports(nil); return err },
		"CurrentImport": func() error { _, err := c.Contacts.CurrentImport(); return err },
		"GetImport":     func() error { _, err := c.Contacts.GetImport("imp1"); return err },
		"CancelImport":  func() error { _, err := c.Contacts.CancelImport("imp1"); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Contacts.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
