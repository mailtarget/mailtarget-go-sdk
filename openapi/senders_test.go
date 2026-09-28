package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This endpoint returns {"data": [...]} with no meta, so it is not paginated.
func TestSendersService_List(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/domain/senders",
		httpmock.NewStringResponder(200, `{"data":[
			{"id":"snd1","email":"no-reply@mtarget.co","spf":true,"dkim":true,"permitted":true},
			{"id":"snd2","email":"hello@mtarget.co","permitted":false}
		]}`))

	senders, err := c.Senders.List()

	require.NoError(t, err)
	require.Len(t, senders, 2)
	assert.Equal(t, "no-reply@mtarget.co", senders[0].Email)
	assert.True(t, senders[0].DKIM)
	assert.False(t, senders[1].Permitted)
}

func TestSendersService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/domain/senders/snd1",
		httpmock.NewStringResponder(200,
			`{"data":{"id":"snd1","email":"no-reply@mtarget.co","domainStatus":"verified"}}`))

	sender, err := c.Senders.Get("snd1")

	require.NoError(t, err)
	assert.Equal(t, "snd1", sender.ID)
	assert.Equal(t, "verified", sender.DomainStatus)
}

func TestSendersService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/domain/senders",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":"snd3","email":"new@mtarget.co"}}`), nil
		})

	sender, err := c.Senders.Create(&SenderRequest{Email: "new@mtarget.co", Name: "New"})

	require.NoError(t, err)
	assert.Equal(t, "snd3", sender.ID)
	assert.JSONEq(t, `{"email":"new@mtarget.co","name":"New"}`, gotBody)
}

func TestSendersService_Update(t *testing.T) {
	c := withSecret(t)

	var gotMethod, gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/domain/senders/snd1",
		func(req *http.Request) (*http.Response, error) {
			gotMethod, gotBody = req.Method, readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"id":"snd1","name":"Renamed"}}`), nil
		})

	sender, err := c.Senders.Update("snd1", &SenderRequest{Name: "Renamed"})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.JSONEq(t, `{"name":"Renamed"}`, gotBody)
	assert.Equal(t, "Renamed", sender.Name)
}

func TestSendersService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/domain/senders/snd1",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Senders.Delete("snd1"))
}

func TestSendersService_CheckDomain(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/domain/senders/check-domain",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{
				"id":"dom1","spf":true,"dkim":true,"dmarc":false,
				"hasSendingDomain":true,"permitted":true,"lastValidate":1750000000
			}}`), nil
		})

	status, err := c.Senders.CheckDomain("no-reply@mtarget.co")

	require.NoError(t, err)
	assert.JSONEq(t, `{"email":"no-reply@mtarget.co"}`, gotBody)
	assert.True(t, status.HasSendingDomain)
	assert.False(t, status.DMARC)
	assert.Equal(t, int64(1750000000), status.LastValidate)
}

func TestSendersService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List":        func() error { _, err := c.Senders.List(); return err },
		"Get":         func() error { _, err := c.Senders.Get("1"); return err },
		"Create":      func() error { _, err := c.Senders.Create(&SenderRequest{}); return err },
		"Update":      func() error { _, err := c.Senders.Update("1", &SenderRequest{}); return err },
		"Delete":      func() error { return c.Senders.Delete("1") },
		"CheckDomain": func() error { _, err := c.Senders.CheckDomain("a@b.co"); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Senders.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
