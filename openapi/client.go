// Package openapi is the client for the Mailtarget Open API — contacts,
// analytics and the other account level resources.
//
// It is a separate capability from sending email: the Open API authenticates
// with an Open API secret key, which is issued on request, while sending only
// needs the Mailtarget API key every account can generate from the dashboard.
// A client without a secret key is legal to build; every call is then refused
// locally, before any request is made, with a *ConfigError.
//
// Most callers reach this package through the parent SDK, which embeds this
// Client so its resources are available directly:
//
//	client := layang.NewMailtargetClient(apiKey, layang.WithOpenAPISecretKey(secretKey))
//	page, err := client.Contacts.List(nil)
//
// It also works standalone when sending email is not needed:
//
//	oa := openapi.New(openapi.WithSecretKey(secretKey))
//	page, err := oa.Contacts.List(nil)
package openapi

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultBaseURL is the production Open API host.
	DefaultBaseURL = "https://api.mailtarget.co"

	// APIVersion is the Open API version this client targets.
	APIVersion = "v1"

	// DefaultTimeout bounds every Open API request.
	DefaultTimeout = 30 * time.Second
)

// Client is the Mailtarget Open API client. Resources hang off it as fields,
// so adding a capability means adding one file and one line of wiring here.
type Client struct {
	Contacts  *ContactsService
	Analytics *AnalyticsService

	baseURL    string
	apiVersion string
	secretKey  string
	timeout    time.Duration
	httpClient *http.Client
	resty      *resty.Client
}

// Option customises a Client at construction time.
type Option func(*Client)

// WithSecretKey sets the Open API secret key. Without it every call fails with
// a *ConfigError before reaching the network.
func WithSecretKey(secretKey string) Option {
	return func(c *Client) { c.secretKey = secretKey }
}

// WithBaseURL overrides the Open API host. Intended for tests and staging.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

// WithTimeout overrides the per-request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) { c.timeout = timeout }
}

// WithHTTPClient reuses an existing HTTP client, so the Open API and the
// transmission API can share one connection pool.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// New builds an Open API client. The secret key is optional, so a client can be
// created before it is known and filled in later with SetSecretKey.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		apiVersion: APIVersion,
		timeout:    DefaultTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient != nil {
		c.resty = resty.NewWithClient(c.httpClient)
	} else {
		c.resty = resty.New()
	}

	c.Contacts = &ContactsService{c: c}
	c.Analytics = &AnalyticsService{c: c}
	return c
}

// SetSecretKey sets the Open API secret key after construction, for callers
// that resolve it lazily.
func (c *Client) SetSecretKey(secretKey string) {
	c.secretKey = secretKey
}

// HasAccess reports whether a secret key is configured, so callers can skip
// Open API features instead of handling a *ConfigError.
func (c *Client) HasAccess() bool {
	return strings.TrimSpace(c.secretKey) != ""
}

// String deliberately redacts the secret key: Go's fmt package prints
// unexported fields, so %v on this client would otherwise leak it.
func (c *Client) String() string {
	return "openapi.Client{baseURL:" + c.baseURL +
		", access:" + strconv.FormatBool(c.HasAccess()) + "}"
}

func (c *Client) url(path string) string {
	return strings.TrimSuffix(c.baseURL, "/") + "/" + c.apiVersion + path
}
