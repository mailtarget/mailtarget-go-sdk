package openapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sendRequest() *SendEmailRequest {
	return &SendEmailRequest{
		From:     &Address{Email: "no-reply@mtarget.co", Name: "Mailtarget"},
		Subject:  "Hello",
		To:       []Address{{Email: "recipient@example.com"}},
		BodyText: "hi",
	}
}

func TestTransmissionsService_Send(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey), WithTransmissionAPIKey("mt_api_key"))

	var gotBody, gotAuth string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/transmissions",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			gotAuth = req.Header.Get("Authorization")
			return httpmock.NewStringResponse(200, `{"data":{"transmissionId":"tx-1"}}`), nil
		})

	result, err := c.Transmissions.Send(sendRequest())

	require.NoError(t, err)
	assert.Equal(t, "tx-1", result.TransmissionID)
	// The secret key authenticates the call; the API key travels in the body.
	assert.Equal(t, "Bearer "+testSecretKey, gotAuth)

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(gotBody), &payload))
	assert.Equal(t, "mt_api_key", payload["apiKey"])
	assert.Equal(t, "Hello", payload["subject"])
	// Unset optional fields stay out of the payload.
	assert.NotContains(t, payload, "cc")
	assert.NotContains(t, payload, "templateId")
}

// Filling in the API key must not mutate the caller's request.
func TestTransmissionsService_SendDoesNotMutateCallerRequest(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey), WithTransmissionAPIKey("mt_api_key"))
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/transmissions",
		httpmock.NewStringResponder(200, `{"data":{"transmissionId":"tx-1"}}`))

	req := sendRequest()
	_, err := c.Transmissions.Send(req)

	require.NoError(t, err)
	assert.Empty(t, req.APIKey, "the caller's request must be left untouched")
}

// An explicit APIKey on the request wins over the client's.
func TestTransmissionsService_SendKeepsExplicitAPIKey(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey), WithTransmissionAPIKey("mt_client_key"))

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/transmissions",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"transmissionId":"tx-1"}}`), nil
		})

	req := sendRequest()
	req.APIKey = "mt_explicit_key"
	_, err := c.Transmissions.Send(req)

	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(gotBody), &payload))
	assert.Equal(t, "mt_explicit_key", payload["apiKey"])
}

// Standalone use without a transmission API key cannot satisfy the endpoint.
func TestTransmissionsService_SendRequiresAPIKey(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey))

	_, err := c.Transmissions.Send(sendRequest())

	assert.ErrorContains(t, err, "requires a Mailtarget API key")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

// from, subject and to are required by the spec, so a payload that cannot
// succeed is rejected before it costs a round trip.
func TestTransmissionsService_SendRequiresMandatoryFields(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey), WithTransmissionAPIKey("mt_api_key"))

	noFrom := sendRequest()
	noFrom.From = nil
	blankFrom := sendRequest()
	blankFrom.From = &Address{Name: "No email"}
	noSubject := sendRequest()
	noSubject.Subject = "   "
	noTo := sendRequest()
	noTo.To = nil

	_, errNoFrom := c.Transmissions.Send(noFrom)
	_, errBlankFrom := c.Transmissions.Send(blankFrom)
	_, errNoSubject := c.Transmissions.Send(noSubject)
	_, errNoTo := c.Transmissions.Send(noTo)

	assert.ErrorContains(t, errNoFrom, "requires From")
	assert.ErrorContains(t, errBlankFrom, "requires From")
	assert.ErrorContains(t, errNoSubject, "requires a Subject")
	assert.ErrorContains(t, errNoTo, "requires at least one To")
	assert.Equal(t, 0, httpmock.GetTotalCallCount(), "invalid payloads must not reach the network")
}

func TestTransmissionsService_SendRejectsNilRequest(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey), WithTransmissionAPIKey("mt_api_key"))

	_, err := c.Transmissions.Send(nil)

	assert.ErrorContains(t, err, "requires a request")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

// The secret key guard still comes first, before the API key check.
func TestTransmissionsService_GuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t, WithTransmissionAPIKey("mt_api_key"))

	_, err := c.Transmissions.Send(sendRequest())

	var configErr *ConfigError
	assert.ErrorAs(t, err, &configErr)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestTransmissionsService_SendCarriesFullPayload(t *testing.T) {
	c := newTestClient(t, WithSecretKey(testSecretKey), WithTransmissionAPIKey("mt_api_key"))

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/transmissions",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"transmissionId":"tx-2"}}`), nil
		})

	_, err := c.Transmissions.Send(&SendEmailRequest{
		From:              &Address{Email: "s@b.co"},
		Subject:           "Full",
		To:                []Address{{Email: "a@b.co", Name: "A"}},
		CC:                []Address{{Email: "c@b.co"}},
		BodyHTML:          "<p>hi</p>",
		Headers:           []Header{{Name: "X-Trace", Value: "1"}},
		Attachments:       []Attachment{{Filename: "a.png", MimeType: "image/png", Value: "AAAA", ContentID: "cid1"}},
		Metadata:          map[string]string{"k": "v"},
		SubstitutionData:  map[string]string{"name": "A"},
		TemplateID:        "tpl1",
		OptionsAttributes: &OptionsAttributes{ClickTracking: true, Transactional: true},
	})

	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(gotBody), &payload))
	assert.Equal(t, "tpl1", payload["templateId"])
	assert.Equal(t, "<p>hi</p>", payload["bodyHtml"])
	assert.Len(t, payload["attachments"], 1)
	assert.Equal(t, map[string]any{"clickTracking": true, "openTracking": false, "transactional": true},
		payload["optionsAttributes"])
}
