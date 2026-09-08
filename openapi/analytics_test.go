package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsService_Summary(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/analytics/summary",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{"data":{
				"name":"All","sent":100,"delivered":98,"opened":40,"openRate":40.8,
				"period":{"from":"2026-08-01","to":"2026-08-31"}
			}}`), nil
		})

	summary, err := c.Analytics.Summary(&AnalyticsSummaryParams{From: "2026-08-01", To: "2026-08-31"})

	require.NoError(t, err)
	assert.Equal(t, "from=2026-08-01&to=2026-08-31", gotQuery)
	assert.Equal(t, 100, summary.Sent)
	assert.Equal(t, 40.8, summary.OpenRate)
	assert.Equal(t, "2026-08-31", summary.Period.To)
}

func TestAnalyticsService_SummaryPassesOptionalParams(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/analytics/summary",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{"data":{}}`), nil
		})

	_, err := c.Analytics.Summary(&AnalyticsSummaryParams{
		From: "2026-08-01", To: "2026-08-31", Sandbox: true, SandboxSubAccountID: 7,
	})

	require.NoError(t, err)
	assert.Equal(t, "from=2026-08-01&sandbox=true&sandboxSubAccountId=7&to=2026-08-31", gotQuery)
}

// from/to are required by the API, so the SDK rejects the call locally.
func TestAnalyticsService_SummaryRequiresDateRange(t *testing.T) {
	c := withSecret(t)

	_, errNil := c.Analytics.Summary(nil)
	_, errMissingTo := c.Analytics.Summary(&AnalyticsSummaryParams{From: "2026-08-01"})

	assert.ErrorContains(t, errNil, "From and To")
	assert.ErrorContains(t, errMissingTo, "From and To")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAnalyticsService_SummaryBreakdown(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/analytics/summary/breakdown",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200,
				`{"data":[{"name":"sender-a","sent":10},{"name":"sender-b","sent":5}]}`), nil
		})

	summaries, err := c.Analytics.SummaryBreakdown(&AnalyticsSummaryParams{
		From: "2026-08-01", To: "2026-08-31", GroupBy: "sender",
	})

	require.NoError(t, err)
	assert.Equal(t, "from=2026-08-01&groupBy=sender&to=2026-08-31", gotQuery)
	require.Len(t, summaries, 2)
	assert.Equal(t, "sender-a", summaries[0].Name)
	assert.Equal(t, 5, summaries[1].Sent)
}

func TestAnalyticsService_Transmission(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/analytics/transmission/tx-1",
		httpmock.NewStringResponder(200, `{"data":{
			"transmissionId":"tx-1","subject":"Hi","status":"delivered","recipientCount":2,
			"from":{"email":"s@b.co","name":"S"},
			"stats":{"delivered":2,"opened":1}
		}}`))

	detail, err := c.Analytics.Transmission("tx-1")

	require.NoError(t, err)
	assert.Equal(t, "tx-1", detail.TransmissionID)
	assert.Equal(t, "s@b.co", detail.From.Email)
	assert.Equal(t, 1, detail.Stats.Opened)
}

// The events payload nests the list under data.events.
func TestAnalyticsService_TransmissionEvents(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/analytics/transmission/tx-1/events",
		httpmock.NewStringResponder(200, `{"data":{"events":[
			{"type":"delivery","recipient":"a@b.co","timestamp":1756800000},
			{"type":"open","recipient":"a@b.co","userAgent":"Mozilla"}
		]}}`))

	events, err := c.Analytics.TransmissionEvents("tx-1")

	require.NoError(t, err)
	require.Len(t, events, 2)
	assert.Equal(t, "delivery", events[0].Type)
	assert.Equal(t, int64(1756800000), events[0].Timestamp)
	assert.Equal(t, "Mozilla", events[1].UserAgent)
}

func TestAnalyticsService_TransmissionEventsHandlesEmptyPayload(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/analytics/transmission/tx-1/events",
		httpmock.NewStringResponder(200, `{"data":{}}`))

	events, err := c.Analytics.TransmissionEvents("tx-1")

	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestAnalyticsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)
	params := &AnalyticsSummaryParams{From: "2026-08-01", To: "2026-08-31"}

	calls := map[string]func() error{
		"Summary":            func() error { _, err := c.Analytics.Summary(params); return err },
		"SummaryBreakdown":   func() error { _, err := c.Analytics.SummaryBreakdown(params); return err },
		"Transmission":       func() error { _, err := c.Analytics.Transmission("tx-1"); return err },
		"TransmissionEvents": func() error { _, err := c.Analytics.TransmissionEvents("tx-1"); return err },
	}

	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Analytics.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount(), "no guarded call may reach the network")
}
