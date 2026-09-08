package openapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
)

// request describes a single Open API call. Resource methods fill it in and
// hand it to the helpers below; they never touch HTTP themselves.
type request struct {
	// capability names the caller for guard errors, e.g. "Contacts.List".
	capability string
	method     string
	// path is appended to the versioned base URL, e.g. "/contacts/123".
	path  string
	query url.Values
	body  any

	// result receives the `data` member of the response envelope.
	result any
	// meta receives the `meta` block of a paginated response.
	meta *Meta
	// rawResult receives the whole response body instead of `data`, for the
	// few endpoints that answer with a bare object such as {"count": N}.
	rawResult any
}

// do is the single path every Open API call takes. It enforces the secret key
// guard, signs and bounds the request, and turns the response into either the
// caller's target or a typed error.
//
// Because resource methods cannot reach the network any other way, a new
// resource cannot accidentally skip the guard.
func (c *Client) do(req request) error {
	if !c.HasAccess() {
		return &ConfigError{Capability: req.capability}
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	r := c.resty.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+c.secretKey).
		SetHeader("Accept", "application/json")

	if len(req.query) > 0 {
		r.SetQueryParamsFromValues(req.query)
	}
	if req.body != nil {
		r.SetHeader("Content-Type", "application/json").SetBody(req.body)
	}

	resp, err := r.Execute(strings.ToUpper(req.method), c.url(req.path))
	if err != nil {
		return err
	}

	body := resp.Body()

	if resp.IsError() {
		apiErr := &Error{StatusCode: resp.StatusCode()}
		var env errorEnvelope
		if len(body) > 0 && json.Unmarshal(body, &env) == nil {
			apiErr.Code = env.Error
			apiErr.Message = env.Message
		}
		if apiErr.Message == "" {
			apiErr.Message = resp.Status()
		}
		return apiErr
	}

	if len(body) == 0 {
		return nil
	}
	if req.rawResult != nil {
		return json.Unmarshal(body, req.rawResult)
	}
	if req.result == nil && req.meta == nil {
		return nil
	}

	var env dataEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return err
	}
	if req.result != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, req.result); err != nil {
			return err
		}
	}
	if req.meta != nil && env.Meta != nil {
		*req.meta = *env.Meta
	}
	return nil
}

// The helpers below cover the four response shapes the Open API uses, so a
// resource method is a single call instead of a decode-and-assemble block.
// They are package level functions because Go methods cannot take type
// parameters.

// object decodes `data` into a single T.
func object[T any](c *Client, req request) (*T, error) {
	var out T
	req.result = &out
	if err := c.do(req); err != nil {
		return nil, err
	}
	return &out, nil
}

// objects decodes `data` into a slice of T, for list endpoints without meta.
func objects[T any](c *Client, req request) ([]T, error) {
	var out []T
	req.result = &out
	if err := c.do(req); err != nil {
		return nil, err
	}
	return out, nil
}

// paged decodes `data` plus `meta` into a Page of T.
func paged[T any](c *Client, req request) (*Page[T], error) {
	var items []T
	var meta Meta
	req.result = &items
	req.meta = &meta
	if err := c.do(req); err != nil {
		return nil, err
	}
	return &Page[T]{Items: items, Meta: meta}, nil
}

// bare decodes the whole response body into T, for the endpoints that answer
// without the `data` envelope.
func bare[T any](c *Client, req request) (*T, error) {
	var out T
	req.rawResult = &out
	if err := c.do(req); err != nil {
		return nil, err
	}
	return &out, nil
}
