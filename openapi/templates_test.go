package openapi

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplatesService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/templates",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200, `{
				"data":[{"id":"tpl1","name":"Welcome","emailType":"regular","labels":["onboarding"]}],
				"meta":{"page":1,"perPage":20,"total":1}
			}`), nil
		})

	page, err := c.Templates.List(&ListTemplatesParams{
		Page: 1, PerPage: 20, Search: "Wel", Type: "custom", Sort: "updatedAt",
		Category: "promo", EmailType: "regular",
	})

	require.NoError(t, err)
	assert.Equal(t,
		"category=promo&emailType=regular&page=1&perPage=20&search=Wel&sort=updatedAt&type=custom", gotQuery)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "Welcome", page.Items[0].Name)
	assert.Equal(t, []string{"onboarding"}, page.Items[0].Labels)
}

func TestTemplatesService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/templates/tpl1",
		httpmock.NewStringResponder(200,
			`{"data":{"id":"tpl1","content":"<p>hi</p>","css":"p{}","body":"{}"}}`))

	tpl, err := c.Templates.Get("tpl1")

	require.NoError(t, err)
	assert.Equal(t, "<p>hi</p>", tpl.Content)
	assert.Equal(t, "p{}", tpl.CSS)
}

func TestCampaignsService_SetSchedule(t *testing.T) {
	c := withSecret(t)

	var gotBody string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/campaigns/cmp1/set-schedule",
		func(req *http.Request) (*http.Response, error) {
			gotBody = readBody(t, req)
			return httpmock.NewStringResponse(200,
				`{"data":{"id":"cmp1","dueDate":"2026-10-01 09:00","startType":"scheduled"}}`), nil
		})

	campaign, err := c.Campaigns.SetSchedule("cmp1", "2026-10-01 09:00")

	require.NoError(t, err)
	assert.JSONEq(t, `{"dueDate":"2026-10-01 09:00"}`, gotBody)
	assert.Equal(t, "2026-10-01 09:00", campaign.DueDate)
}

func TestCampaignsService_SetScheduleRequiresDueDate(t *testing.T) {
	c := withSecret(t)

	_, err := c.Campaigns.SetSchedule("cmp1", "  ")

	assert.ErrorContains(t, err, "requires a dueDate")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCampaignsService_CancelSchedule(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/campaigns/cmp1/cancel-schedule",
		httpmock.NewStringResponder(200, `{"message":"schedule cancelled"}`))

	assert.NoError(t, c.Campaigns.CancelSchedule("cmp1"))
}

func TestAPIKeysService_CreateRequiresPermissionIDs(t *testing.T) {
	c := withSecret(t)

	_, errNil := c.APIKeys.Create(nil)
	_, errEmpty := c.APIKeys.Create(&CreateAPIKeyRequest{Name: "ci"})

	assert.ErrorContains(t, errNil, "at least one PermissionID")
	assert.ErrorContains(t, errEmpty, "at least one PermissionID")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestTemplatesAndScheduleGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"Templates.List":           func() error { _, err := c.Templates.List(nil); return err },
		"Templates.Get":            func() error { _, err := c.Templates.Get("tpl1"); return err },
		"Campaigns.SetSchedule":    func() error { _, err := c.Campaigns.SetSchedule("cmp1", "2026-10-01"); return err },
		"Campaigns.CancelSchedule": func() error { return c.Campaigns.CancelSchedule("cmp1") },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
