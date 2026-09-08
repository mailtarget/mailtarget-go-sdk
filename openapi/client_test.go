package openapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew_Defaults(t *testing.T) {
	c := New()

	assert.Equal(t, DefaultBaseURL, c.baseURL)
	assert.Equal(t, APIVersion, c.apiVersion)
	assert.Equal(t, DefaultTimeout, c.timeout)
	assert.False(t, c.HasAccess())
	assert.NotNil(t, c.Contacts)
	assert.NotNil(t, c.Analytics)
}

func TestNew_Options(t *testing.T) {
	c := New(
		WithSecretKey(testSecretKey),
		WithBaseURL("https://staging.test/"),
		WithTimeout(5*time.Second),
	)

	assert.True(t, c.HasAccess())
	assert.Equal(t, 5*time.Second, c.timeout)
	// A trailing slash on the base URL must not double up in the path.
	assert.Equal(t, "https://staging.test/v1/contacts", c.url("/contacts"))
}

func TestClient_SetSecretKey(t *testing.T) {
	c := New()
	assert.False(t, c.HasAccess())

	c.SetSecretKey(testSecretKey)

	assert.True(t, c.HasAccess())
}

func TestClient_HasAccessIgnoresBlankKey(t *testing.T) {
	assert.False(t, New(WithSecretKey("   ")).HasAccess())
}

func TestClient_StringDoesNotLeakSecret(t *testing.T) {
	dump := New(WithSecretKey("oa_super_secret_key")).String()

	assert.NotContains(t, dump, "oa_super_secret_key")
	assert.Contains(t, dump, "access:true")
}
