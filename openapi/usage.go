package openapi

import "net/http"

// UsageService reports the account's email sending usage.
type UsageService struct{ c *Client }

// Usage is the sending usage of the current billing period.
//
// On plans without a fixed monthly quota (for example FREE or
// PAY_AS_YOU_GROW) EmailsQuota, EmailsRemaining and PercentUsed are 0, which
// does not mean sending is blocked.
type Usage struct {
	Plan string `json:"plan"`
	// PeriodEnd is a Unix millisecond timestamp, or 0 when not reported.
	PeriodEnd       int64   `json:"periodEnd"`
	EmailsSent      int64   `json:"emailsSent"`
	EmailsQuota     int64   `json:"emailsQuota"`
	EmailsRemaining int64   `json:"emailsRemaining"`
	PercentUsed     float64 `json:"percentUsed"`
}

// Get returns how many emails the account has sent in the current billing
// period against its quota.
func (s *UsageService) Get() (*Usage, error) {
	return object[Usage](s.c, request{
		capability: "Usage.Get",
		method:     http.MethodGet,
		path:       "/usage",
	})
}
