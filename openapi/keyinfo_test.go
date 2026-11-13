package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyInfoService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/open-api/key-info",
		httpmock.NewStringResponder(200, `{"data":{
			"apiKeyId":"k1","companyId":12,"emailApiId":34,
			"scopes":["contacts.read","forms.write"],"unifyAccountId":"u1"
		}}`))

	info, err := c.KeyInfo.Get()

	require.NoError(t, err)
	assert.Equal(t, "k1", info.APIKeyID)
	assert.Equal(t, 12, info.CompanyID)
	assert.Equal(t, 34, info.EmailAPIID)
	assert.Equal(t, []string{"contacts.read", "forms.write"}, info.Scopes)
	assert.Equal(t, "u1", info.UnifyAccountID)
}

func TestKeyInfoService_CheckTransmissionKey(t *testing.T) {
	c := withSecret(t)

	var bodies []string
	responses := []string{`{"data":{"belongsToCompany":true}}`, `{"data":{"belongsToCompany":false}}`}
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/open-api/key-info/transmission-key",
		func(req *http.Request) (*http.Response, error) {
			bodies = append(bodies, readBody(t, req))
			return httpmock.NewStringResponse(200, responses[len(bodies)-1]), nil
		})

	ours, err := c.KeyInfo.CheckTransmissionKey("tx_ours")
	require.NoError(t, err)
	assert.True(t, ours)

	theirs, err := c.KeyInfo.CheckTransmissionKey("tx_theirs")
	require.NoError(t, err)
	assert.False(t, theirs)

	require.Len(t, bodies, 2)
	assert.JSONEq(t, `{"apiKey":"tx_ours"}`, bodies[0])
}

func TestKeyInfoService_CheckTransmissionKeyRequiresKey(t *testing.T) {
	c := withSecret(t)

	_, err := c.KeyInfo.CheckTransmissionKey("  ")

	assert.ErrorContains(t, err, "requires an API key")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestKeyInfoService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"Get":                  func() error { _, err := c.KeyInfo.Get(); return err },
		"CheckTransmissionKey": func() error { _, err := c.KeyInfo.CheckTransmissionKey("tx"); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "KeyInfo.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
