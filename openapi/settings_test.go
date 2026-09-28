package openapi

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Both settings endpoints answer with the object directly, not wrapped in the
// usual {"data": ...} envelope.
func TestSettingsService_Company(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/settings/mtarget/company",
		httpmock.NewStringResponder(200, `{
			"companyId":"cmp1","name":"Sunrise Bank","packet":"enterprise",
			"paid":true,"autoLinkUtm":false,"expiredDate":1790000000
		}`))

	company, err := c.Settings.Company()

	require.NoError(t, err)
	assert.Equal(t, "cmp1", company.CompanyID)
	assert.Equal(t, "Sunrise Bank", company.Name)
	assert.True(t, company.Paid)
	assert.Equal(t, int64(1790000000), company.ExpiredDate)
}

func TestSettingsService_Profile(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/settings/mtarget/profile",
		httpmock.NewStringResponder(200, `{
			"userId":"usr1","email":"ryan@mtarget.co","fullname":"Ryan",
			"tours":["welcome"],"validate":true
		}`))

	profile, err := c.Settings.Profile()

	require.NoError(t, err)
	assert.Equal(t, "usr1", profile.UserID)
	assert.Equal(t, "ryan@mtarget.co", profile.Email)
	assert.Equal(t, []string{"welcome"}, profile.Tours)
	assert.True(t, profile.Validate)
}

// A `data` wrapper must not be mistaken for the payload on these endpoints.
func TestSettingsService_CompanyIgnoresEnvelopeShape(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/settings/mtarget/company",
		httpmock.NewStringResponder(200, `{"name":"Direct Object"}`))

	company, err := c.Settings.Company()

	require.NoError(t, err)
	assert.Equal(t, "Direct Object", company.Name)
}

func TestSettingsService_AllMethodsGuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	calls := map[string]func() error{
		"Company": func() error { _, err := c.Settings.Company(); return err },
		"Profile": func() error { _, err := c.Settings.Profile(); return err },
	}
	for name, call := range calls {
		var configErr *ConfigError
		assert.ErrorAs(t, call(), &configErr, "Settings.%s must be guarded", name)
	}
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
