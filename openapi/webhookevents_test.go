package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookEventsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/webhook-events",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{"data":[
				{"id":1,"name":"Delivered","slug":"delivered","detail":"Message delivered"},
				{"id":2,"name":"Bounced","slug":"bounced"}
			]}`), nil
		})

	events, err := c.WebhookEvents.List(&ListWebhookEventsParams{Page: 1, Size: 50, Name: "e"})

	require.NoError(t, err)
	// This endpoint spells the page size "size", not "perPage".
	assert.Equal(t, "name=e&page=1&size=50", gotQuery)
	require.Len(t, events, 2)
	assert.Equal(t, "Message delivered", events[0].Detail)
	assert.Equal(t, "bounced", events[1].Slug)
}

func TestWebhookEventsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/webhook-events/1",
		httpmock.NewStringResponder(200, `{"data":{"id":1,"name":"Delivered","slug":"delivered"}}`))

	event, err := c.WebhookEvents.Get(1)

	require.NoError(t, err)
	assert.Equal(t, "Delivered", event.Name)
}

func TestWebhookEventsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"List": func() error { _, err := c.WebhookEvents.List(nil); return err },
		"Get":  func() error { _, err := c.WebhookEvents.Get(1); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "WebhookEvents.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
