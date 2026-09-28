package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendingDomainsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/domain/sending",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":5,"domain":"mtarget.co","status":"verified","types":["sending","tracking"],"default":true}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	page, err := c.SendingDomains.List(&ListSendingDomainsParams{Page: 1, PerPage: 20})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=20", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "mtarget.co", page.Items[0].Domain)
	assert.Equal(t, []string{"sending", "tracking"}, page.Items[0].Types)
	assert.True(t, page.Items[0].Default)
}

func TestSendingDomainsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/domain/sending/5",
		httpmock.NewStringResponder(200,
			`{"data":{"id":5,"domain":"mtarget.co","subAccountID":3,"createdAt":1750000000}}`))

	domain, err := c.SendingDomains.Get(5)

	require.NoError(t, err)
	assert.Equal(t, 5, domain.ID)
	assert.Equal(t, 3, domain.SubAccountID)
	assert.Equal(t, int64(1750000000), domain.CreatedAt)
}

func TestSendingDomainsService_VerifyTXT(t *testing.T) {
	c := withSecret(t)

	var gotMethod string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/domain/sending/5/verify-txt",
		func(req *http.Request) (*http.Response, error) {
			gotMethod = req.Method
			return httpmock.NewStringResponse(200, `{"data":{"id":5,"status":"verified"}}`), nil
		})

	domain, err := c.SendingDomains.VerifyTXT(5)

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "verified", domain.Status)
}

func TestSendingDomainsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List":      func() error { _, err := c.SendingDomains.List(nil); return err },
		"Get":       func() error { _, err := c.SendingDomains.Get(1); return err },
		"VerifyTXT": func() error { _, err := c.SendingDomains.VerifyTXT(1); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "SendingDomains.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
