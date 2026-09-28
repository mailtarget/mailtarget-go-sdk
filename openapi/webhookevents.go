package openapi

import (
	"net/http"
	"net/url"
	"strconv"
)

// WebhookEventsService reads the platform catalog of webhook event types,
// whose IDs go into the EventIDs of a webhook request.
type WebhookEventsService struct{ c *Client }

// ListWebhookEventsParams holds the optional query parameters of List.
type ListWebhookEventsParams struct {
	Page int
	// Size is the page size. Like SubAccounts, this endpoint spells it "size"
	// rather than "perPage".
	Size int
	// Name filters by event name.
	Name string
}

func (p *ListWebhookEventsParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.Size > 0 {
		values.Set("size", strconv.Itoa(p.Size))
	}
	if p.Name != "" {
		values.Set("name", p.Name)
	}
	return values
}

// List returns the webhook event types available to subscribe to. The
// response carries no pagination metadata.
func (s *WebhookEventsService) List(params *ListWebhookEventsParams) ([]WebhookEvent, error) {
	return objects[WebhookEvent](s.c, request{
		capability: "WebhookEvents.List",
		method:     http.MethodGet,
		path:       "/webhook-events",
		query:      params.values(),
	})
}

// Get returns one webhook event type by ID.
func (s *WebhookEventsService) Get(id int) (*WebhookEvent, error) {
	return object[WebhookEvent](s.c, request{
		capability: "WebhookEvents.Get",
		method:     http.MethodGet,
		path:       "/webhook-events/" + strconv.Itoa(id),
	})
}
