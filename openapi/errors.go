package openapi

import "fmt"

// Error reports a non-2xx response from the Open API.
type Error struct {
	// StatusCode is the HTTP status returned by the Open API.
	StatusCode int
	// Code is the machine readable `error` field of the response body.
	Code string
	// Message is the human readable `message` field of the response body.
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("mailtarget open api: http %d: %s: %s", e.StatusCode, e.Code, e.Message)
}

// ConfigError reports a credential the client needs but was never given. It is
// returned before any request is made, so a caller missing the Open API secret
// key never produces network traffic.
type ConfigError struct {
	// Capability is the call that was refused, e.g. "Contacts.List".
	Capability string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("open API secret key is required to call %s — pass "+
		"layang.WithOpenAPISecretKey(secretKey) when building the client, or call "+
		"SetSecretKey(secretKey) on the Open API client", e.Capability)
}
