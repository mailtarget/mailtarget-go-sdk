package layang

import (
	"strconv"
	"time"

	"github.com/mailtarget/mailtarget-go-sdk/openapi"
)

// MailtargetClient is the entry point of the SDK.
//
// It bundles two capabilities that use two different credentials:
//
//   - Sending email through the Transmission API, exposed as Layang. This only
//     needs the Mailtarget API key, which every account can generate from the
//     dashboard, and is the preferred way to send email.
//   - The Mailtarget Open API resources (e.g. Contacts, Analytics), which
//     additionally need an Open API secret key. That key is optional: it is
//     issued on request, so a client built without it still sends email
//     normally, and only Open API calls are rejected — locally, before any
//     request is made, with a *openapi.ConfigError.
//
// Open API resources are embedded rather than nested under their own field, so
// they read the same way as in the Java and Python SDKs:
//
//	client.Contacts.List(nil)
//	client.Analytics.Summary(params)
type MailtargetClient struct {
	// Layang is the preferred email sending path. It is the same type returned
	// by NewLayang and never routes through the Open API.
	Layang *Layang

	// *openapi.Client is embedded so its resources (Contacts, Analytics, ...)
	// are promoted directly onto MailtargetClient. Adding a resource to the
	// openapi package makes it available here automatically, with no change
	// needed in this file.
	*openapi.Client
}

// ClientOption customises a MailtargetClient at construction time.
type ClientOption func(*clientConfig)

type clientConfig struct {
	openAPI []openapi.Option
}

// WithOpenAPISecretKey configures the optional Open API secret key. Without it
// every Open API call fails with an *openapi.ConfigError before reaching the
// network.
func WithOpenAPISecretKey(secretKey string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.openAPI = append(cfg.openAPI, openapi.WithSecretKey(secretKey))
	}
}

// WithOpenAPIBaseURL overrides the Open API base URL. Intended for tests and
// staging environments.
func WithOpenAPIBaseURL(baseURL string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.openAPI = append(cfg.openAPI, openapi.WithBaseURL(baseURL))
	}
}

// WithOpenAPITimeout overrides the per-request Open API timeout.
func WithOpenAPITimeout(timeout time.Duration) ClientOption {
	return func(cfg *clientConfig) {
		cfg.openAPI = append(cfg.openAPI, openapi.WithTimeout(timeout))
	}
}

// NewMailtargetClient builds a client from the required Mailtarget API key.
// Pass WithOpenAPISecretKey to also enable the Open API resources:
//
//	client := layang.NewMailtargetClient(apiKey)
//	client := layang.NewMailtargetClient(apiKey, layang.WithOpenAPISecretKey(secretKey))
func NewMailtargetClient(apiKey string, opts ...ClientOption) *MailtargetClient {
	cfg := &clientConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	layang := NewLayang(apiKey)

	// Both capabilities share one HTTP client so they share a connection pool.
	// That is safe because auth headers are always set per request, never as a
	// client default, so neither credential can reach the other's host.
	options := append([]openapi.Option{openapi.WithHTTPClient(layang.resty.GetClient())}, cfg.openAPI...)

	return &MailtargetClient{
		Layang: layang,
		Client: openapi.New(options...),
	}
}

// Send delivers an email through the Transmission API. It is a convenience
// wrapper around Layang.Send and never uses the Open API.
func (c *MailtargetClient) Send(message *Message) (*SuccessResponse, *ErrorResponse, error) {
	return c.Layang.Send(message)
}

// SetOpenAPISecretKey configures the Open API secret key after construction,
// for callers that resolve it lazily. Equivalent to c.SetSecretKey, named to
// match the Java and Python SDKs.
func (c *MailtargetClient) SetOpenAPISecretKey(secretKey string) {
	c.Client.SetSecretKey(secretKey)
}

// HasOpenAPIAccess reports whether an Open API secret key is configured, so
// callers can skip Open API features instead of handling an error. Equivalent
// to c.HasAccess, named to match the Java and Python SDKs.
func (c *MailtargetClient) HasOpenAPIAccess() bool {
	return c.Client.HasAccess()
}

// String reports the client without its credentials: Go's fmt package prints
// unexported fields, so %v would otherwise leak them.
func (c *MailtargetClient) String() string {
	return "layang.MailtargetClient{openAPIAccess:" +
		strconv.FormatBool(c.HasOpenAPIAccess()) + "}"
}
