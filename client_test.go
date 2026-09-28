package layang

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mailtarget/mailtarget-go-sdk/openapi"
	"github.com/stretchr/testify/assert"
)

const openAPISecretKey = "OPEN_API_SECRET"

func TestNewMailtargetClient_Defaults(t *testing.T) {
	client := NewMailtargetClient(privateAPIKey)

	// Sending works with only the Mailtarget API key.
	assert.NotNil(t, client.Layang)
	assert.NoError(t, client.Layang.IsValid())
	// The Open API resources are wired but locked until a secret key is supplied.
	assert.NotNil(t, client.Contacts)
	assert.NotNil(t, client.Analytics)
	assert.False(t, client.HasOpenAPIAccess())
}

func TestNewMailtargetClient_OpenAPIOptions(t *testing.T) {
	client := NewMailtargetClient(
		privateAPIKey,
		WithOpenAPISecretKey(openAPISecretKey),
		WithOpenAPIBaseURL("https://openapi.test"),
	)

	assert.True(t, client.HasOpenAPIAccess())

	httpmock.ActivateNonDefault(client.Layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://openapi.test/v1/contacts",
		httpmock.NewStringResponder(200, `{"data":[{"id":"1"}],"meta":{"total":1}}`))

	page, err := client.Contacts.List(nil)

	assert.NoError(t, err)
	assert.Len(t, page.Items, 1)
}

func TestMailtargetClient_SetOpenAPISecretKey(t *testing.T) {
	client := NewMailtargetClient(privateAPIKey)
	assert.False(t, client.HasOpenAPIAccess())

	client.SetOpenAPISecretKey(openAPISecretKey)

	assert.True(t, client.HasOpenAPIAccess())
	assert.True(t, client.HasAccess()) // promoted from the embedded openapi.Client
}

func TestMailtargetClient_OpenAPICallGuardedWithoutSecretKey(t *testing.T) {
	client := NewMailtargetClient(privateAPIKey, WithOpenAPIBaseURL("https://openapi.test"))

	httpmock.ActivateNonDefault(client.Layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://openapi.test/v1/contacts",
		httpmock.NewStringResponder(200, `{"data":[]}`))

	_, err := client.Contacts.List(nil)

	var configErr *openapi.ConfigError
	assert.ErrorAs(t, err, &configErr)
	assert.Equal(t, 0, httpmock.GetTotalCallCount(), "guard must fire before any HTTP call")
}

func TestMailtargetClient_StringDoesNotLeakCredentials(t *testing.T) {
	client := NewMailtargetClient("mt_super_secret_key", WithOpenAPISecretKey("oa_super_secret_key"))

	dump := client.String()

	assert.NotContains(t, dump, "mt_super_secret_key")
	assert.NotContains(t, dump, "oa_super_secret_key")
	assert.Contains(t, dump, "openAPIAccess:true")
}

// The Open API secret key must never reach the transmission host, and
// configuring it must not change how sending behaves.
func TestMailtargetClient_SendUsesTransmissionHostWithAPIKey(t *testing.T) {
	client := NewMailtargetClient(privateAPIKey, WithOpenAPISecretKey(openAPISecretKey))
	message := client.Layang.NewMessage(subject, body, html, sender, to)

	httpmock.ActivateNonDefault(client.Layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()

	var gotURL, gotAuth string
	httpmock.RegisterResponder("POST", "https://apiconfig.mailtarget.co/v1/layang/transmissions",
		func(req *http.Request) (*http.Response, error) {
			gotURL = req.URL.String()
			gotAuth = req.Header.Get("Authorization")
			return httpmock.NewJsonResponse(200, successResponse)
		})

	actual, _, err := client.Send(message)

	assert.NoError(t, err)
	assert.Equal(t, &successResponse, actual)
	assert.Equal(t, "https://apiconfig.mailtarget.co/v1/layang/transmissions", gotURL)
	assert.Equal(t, "Bearer "+privateAPIKey, gotAuth)
	assert.NotContains(t, gotAuth, openAPISecretKey)
	assert.False(t, strings.Contains(gotURL, "api.mailtarget.co/v1"))
}
