package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCampaignsService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/campaigns",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{
					"id":"cmp1","subject":"Promo","stage":"sent","memberCount":100,"sentCount":98,
					"sender":{"email":"no-reply@mtarget.co","name":"Mailtarget"},
					"recipients":{"labels":["vip"]}
				}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	page, err := c.Campaigns.List(&ListCampaignsParams{Page: 1, PerPage: 20, Search: "Promo"})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=20&search=Promo", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "Promo", page.Items[0].Subject)
	assert.Equal(t, "no-reply@mtarget.co", page.Items[0].Sender.Email)
	assert.Equal(t, []string{"vip"}, page.Items[0].Recipients.Labels)
}

func TestCampaignsService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/campaigns/cmp1",
		httpmock.NewStringResponder(200, `{"data":{
			"id":"cmp1","subject":"Promo","htmlContent":"<p>hi</p>",
			"emailType":"regular","startType":"manual","active":true
		}}`))

	campaign, err := c.Campaigns.Get("cmp1")

	require.NoError(t, err)
	assert.Equal(t, "cmp1", campaign.ID)
	assert.Equal(t, "<p>hi</p>", campaign.HTMLContent)
	assert.True(t, campaign.Active)
}

func TestCampaignsService_Create(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/campaigns",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(201, `{"data":{"id":"cmp9","subject":"New"}}`), nil
		})

	campaign, err := c.Campaigns.Create(&CampaignRequest{
		Subject:    "New",
		Sender:     &CampaignSender{Email: "no-reply@mtarget.co"},
		Recipients: &CampaignRecipients{Labels: []string{"vip"}},
	})

	require.NoError(t, err)
	assert.Equal(t, "cmp9", campaign.ID)
	assert.JSONEq(t,
		`{"subject":"New","sender":{"email":"no-reply@mtarget.co"},"recipients":{"labels":["vip"]}}`,
		gotBody)
}

func TestCampaignsService_Update(t *testing.T) {
	c := withSecret(t)

	var gotMethod, gotBody string
	httpmock.RegisterResponder("PUT", testBaseURL+"/v1/campaigns/cmp1",
		func(req *http.Request) (*http.Response, error) {
			gotMethod, gotBody = req.Method, readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":{"id":"cmp1","subject":"Updated"}}`), nil
		})

	campaign, err := c.Campaigns.Update("cmp1", &CampaignRequest{Subject: "Updated"})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.JSONEq(t, `{"subject":"Updated"}`, gotBody)
	assert.Equal(t, "Updated", campaign.Subject)
}

func TestCampaignsService_Delete(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/campaigns/cmp1",
		httpmock.NewStringResponder(200, `{"message":"campaign deleted"}`))

	assert.NoError(t, c.Campaigns.Delete("cmp1"))
}

func TestCampaignsService_Analytics(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/campaigns/cmp1/analytics",
		httpmock.NewStringResponder(200, `{"data":{
			"campaignId":"cmp1","subject":"Promo","sentCount":100,"deliveredCount":98,
			"openCount":40,"openRate":40.8,"complaintCount":1
		}}`))

	stats, err := c.Campaigns.Analytics("cmp1")

	require.NoError(t, err)
	assert.Equal(t, "cmp1", stats.CampaignID)
	assert.Equal(t, 98, stats.DeliveredCount)
	assert.Equal(t, 40.8, stats.OpenRate)
	assert.Equal(t, 1, stats.ComplaintCount)
}

func TestCampaignsService_Recipients(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/campaigns/cmp1/recipients",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"contactId":"ct1","email":"a@b.co","status":"opened","visitCount":3,"labels":["vip"]}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	page, err := c.Campaigns.Recipients("cmp1", &ListCampaignRecipientsParams{
		Status: "opened", Page: 1, PerPage: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, "page=1&perPage=20&status=opened", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "opened", page.Items[0].Status)
	assert.Equal(t, 3, page.Items[0].VisitCount)
}

// status is required by the API, so the SDK refuses the call locally.
func TestCampaignsService_RecipientsRequiresStatus(t *testing.T) {
	c := withSecret(t)

	_, errNil := c.Campaigns.Recipients("cmp1", nil)
	_, errEmpty := c.Campaigns.Recipients("cmp1", &ListCampaignRecipientsParams{Page: 1})

	assert.ErrorContains(t, errNil, "Status filter")
	assert.ErrorContains(t, errEmpty, "Status filter")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCampaignsService_Send(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/campaigns/cmp1/send",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Campaigns.Send("cmp1"))
}

func TestCampaignsService_SendTest(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/campaigns/cmp1/send-test",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200, `{"data":true}`), nil
		})

	require.NoError(t, c.Campaigns.SendTest("cmp1", "qa@mtarget.co"))
	assert.JSONEq(t, `{"recipient":"qa@mtarget.co"}`, gotBody)
}

func TestCampaignsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)
	recipients := &ListCampaignRecipientsParams{Status: "delivered"}

	calls := map[string]func() error{
		"List":       func() error { _, err := c.Campaigns.List(nil); return err },
		"Get":        func() error { _, err := c.Campaigns.Get("1"); return err },
		"Create":     func() error { _, err := c.Campaigns.Create(&CampaignRequest{}); return err },
		"Update":     func() error { _, err := c.Campaigns.Update("1", &CampaignRequest{}); return err },
		"Delete":     func() error { return c.Campaigns.Delete("1") },
		"Analytics":  func() error { _, err := c.Campaigns.Analytics("1"); return err },
		"Recipients": func() error { _, err := c.Campaigns.Recipients("1", recipients); return err },
		"Send":       func() error { return c.Campaigns.Send("1") },
		"SendTest":   func() error { return c.Campaigns.SendTest("1", "a@b.co") },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Campaigns.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
