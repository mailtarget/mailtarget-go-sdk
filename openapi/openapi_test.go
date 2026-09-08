package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	testSecretKey = "OPEN_API_SECRET"
	testBaseURL   = "https://openapi.test"
)

// newTestClient returns a client wired to a mocked HTTP transport. It uses only
// the public options, so the tests exercise the same construction path callers
// take.
func newTestClient(t *testing.T, opts ...Option) *Client {
	t.Helper()

	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	opts = append([]Option{
		WithHTTPClient(httpClient),
		WithBaseURL(testBaseURL),
	}, opts...)
	return New(opts...)
}

// withSecret returns a client that is allowed to reach the Open API.
func withSecret(t *testing.T) *Client {
	t.Helper()
	return newTestClient(t, WithSecretKey(testSecretKey))
}

// readBody drains a mocked request body, for asserting what was sent.
func readBody(t *testing.T, req *http.Request) string {
	t.Helper()
	if req.Body == nil {
		return ""
	}
	buf := make([]byte, req.ContentLength)
	if _, err := req.Body.Read(buf); err != nil && err.Error() != "EOF" {
		t.Fatalf("read body: %v", err)
	}
	return string(buf)
}
