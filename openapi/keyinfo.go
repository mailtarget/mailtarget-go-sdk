package openapi

import (
	"errors"
	"net/http"
	"strings"
)

// KeyInfoService describes the Open API key in use and checks Transmission
// API keys against it.
type KeyInfoService struct{ c *Client }

// KeyInfo describes the Open API key used for the request.
type KeyInfo struct {
	APIKeyID       string   `json:"apiKeyId"`
	CompanyID      int      `json:"companyId"`
	EmailAPIID     int      `json:"emailApiId"`
	Scopes         []string `json:"scopes"`
	UnifyAccountID string   `json:"unifyAccountId"`
}

// transmissionKeyCheck matches the {"data": {"belongsToCompany": bool}} body
// of CheckTransmissionKey.
type transmissionKeyCheck struct {
	BelongsToCompany bool `json:"belongsToCompany"`
}

// Get returns the company, account and scopes of the Open API secret key the
// client is configured with. It needs no scope.
func (s *KeyInfoService) Get() (*KeyInfo, error) {
	return object[KeyInfo](s.c, request{
		capability: "KeyInfo.Get",
		method:     http.MethodGet,
		path:       "/open-api/key-info",
	})
}

// CheckTransmissionKey reports whether a Transmission API key is valid and
// belongs to the same Mailtarget company as the Open API key. No email is
// sent. It returns false both for an invalid or inactive key and for a valid
// key of another company; the two are not told apart. Requires the
// transmissions.write scope.
func (s *KeyInfoService) CheckTransmissionKey(apiKey string) (bool, error) {
	if strings.TrimSpace(apiKey) == "" {
		return false, errors.New("key info check transmission key requires an API key")
	}
	out, err := object[transmissionKeyCheck](s.c, request{
		capability: "KeyInfo.CheckTransmissionKey",
		method:     http.MethodPost,
		path:       "/open-api/key-info/transmission-key",
		body:       map[string]string{"apiKey": apiKey},
	})
	if err != nil {
		return false, err
	}
	return out.BelongsToCompany, nil
}
