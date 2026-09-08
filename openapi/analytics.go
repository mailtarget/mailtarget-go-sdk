package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// AnalyticsService exposes the analytics endpoints.
type AnalyticsService struct{ c *Client }

// AnalyticsPeriod is the date range an analytics summary covers.
type AnalyticsPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// AnalyticsSummary aggregates delivery and engagement metrics.
type AnalyticsSummary struct {
	Name            string           `json:"name"`
	Period          *AnalyticsPeriod `json:"period"`
	Sent            int              `json:"sent"`
	Delivered       int              `json:"delivered"`
	Opened          int              `json:"opened"`
	Clicked         int              `json:"clicked"`
	Bounced         int              `json:"bounced"`
	Unsubscribed    int              `json:"unsubscribed"`
	DeliveryRate    float64          `json:"deliveryRate"`
	OpenRate        float64          `json:"openRate"`
	ClickRate       float64          `json:"clickRate"`
	BounceRate      float64          `json:"bounceRate"`
	UnsubscribeRate float64          `json:"unsubscribeRate"`
}

// TransmissionFrom is the sender of a transmission.
type TransmissionFrom struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TransmissionStats counts the events recorded for a transmission.
type TransmissionStats struct {
	Delivered    int `json:"delivered"`
	Opened       int `json:"opened"`
	Clicked      int `json:"clicked"`
	Bounced      int `json:"bounced"`
	Unsubscribed int `json:"unsubscribed"`
}

// TransmissionDetail describes a single transactional transmission.
type TransmissionDetail struct {
	TransmissionID string             `json:"transmissionId"`
	Subject        string             `json:"subject"`
	Status         string             `json:"status"`
	CreatedAt      int64              `json:"createdAt"`
	RecipientCount int                `json:"recipientCount"`
	From           *TransmissionFrom  `json:"from"`
	Stats          *TransmissionStats `json:"stats"`
}

// TransmissionEvent is one delivery or engagement event.
type TransmissionEvent struct {
	Type         string `json:"type"`
	Recipient    string `json:"recipient"`
	Subject      string `json:"subject"`
	FriendlyFrom string `json:"friendlyFrom"`
	Reason       string `json:"reason"`
	URL          string `json:"url"`
	UserAgent    string `json:"userAgent"`
	Timestamp    int64  `json:"timestamp"`
}

// transmissionEventList matches the events payload, which nests the list one
// level deeper than usual: {"data": {"events": [...]}}.
type transmissionEventList struct {
	Events []TransmissionEvent `json:"events"`
}

// AnalyticsSummaryParams holds the query parameters of Summary and
// SummaryBreakdown. From and To are required by the API.
type AnalyticsSummaryParams struct {
	// From is the inclusive start date, formatted YYYY-MM-DD.
	From string
	// To is the inclusive end date, formatted YYYY-MM-DD.
	To string
	// GroupBy selects the breakdown dimension, e.g. "sender" or "domain". Only
	// SummaryBreakdown uses it.
	GroupBy string
	// Sandbox includes sandbox traffic in the result.
	Sandbox bool
	// SandboxSubAccountID filters sandbox traffic to one sub-account.
	SandboxSubAccountID int
}

func (p *AnalyticsSummaryParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	values.Set("from", p.From)
	values.Set("to", p.To)
	if p.GroupBy != "" {
		values.Set("groupBy", p.GroupBy)
	}
	if p.Sandbox {
		values.Set("sandbox", "true")
	}
	if p.SandboxSubAccountID > 0 {
		values.Set("sandboxSubAccountId", strconv.Itoa(p.SandboxSubAccountID))
	}
	return values
}

func (p *AnalyticsSummaryParams) validate() error {
	if p == nil || p.From == "" || p.To == "" {
		return errors.New("analytics summary requires From and To dates in YYYY-MM-DD format")
	}
	return nil
}

// Summary returns aggregated analytics for a date range.
func (s *AnalyticsService) Summary(params *AnalyticsSummaryParams) (*AnalyticsSummary, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	return object[AnalyticsSummary](s.c, request{
		capability: "Analytics.Summary",
		method:     http.MethodGet,
		path:       "/analytics/summary",
		query:      params.values(),
	})
}

// SummaryBreakdown returns analytics for a date range grouped by the dimension
// in params.GroupBy.
func (s *AnalyticsService) SummaryBreakdown(params *AnalyticsSummaryParams) ([]AnalyticsSummary, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	return objects[AnalyticsSummary](s.c, request{
		capability: "Analytics.SummaryBreakdown",
		method:     http.MethodGet,
		path:       "/analytics/summary/breakdown",
		query:      params.values(),
	})
}

// Transmission returns detail and stats for one transactional transmission.
func (s *AnalyticsService) Transmission(id string) (*TransmissionDetail, error) {
	return object[TransmissionDetail](s.c, request{
		capability: "Analytics.Transmission",
		method:     http.MethodGet,
		path:       "/analytics/transmission/" + url.PathEscape(id),
	})
}

// TransmissionEvents returns every delivery and engagement event recorded for
// a transmission.
func (s *AnalyticsService) TransmissionEvents(id string) ([]TransmissionEvent, error) {
	out, err := object[transmissionEventList](s.c, request{
		capability: "Analytics.TransmissionEvents",
		method:     http.MethodGet,
		path:       "/analytics/transmission/" + url.PathEscape(id) + "/events",
	})
	if err != nil {
		return nil, err
	}
	return out.Events, nil
}
