package openapi

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsageService_Get(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/usage",
		httpmock.NewStringResponder(200, `{"data":{
			"plan":"PRO","periodEnd":1790000000000,"emailsSent":2500,
			"emailsQuota":10000,"emailsRemaining":7500,"percentUsed":25.5
		}}`))

	usage, err := c.Usage.Get()

	require.NoError(t, err)
	assert.Equal(t, "PRO", usage.Plan)
	assert.Equal(t, int64(1790000000000), usage.PeriodEnd)
	assert.Equal(t, int64(2500), usage.EmailsSent)
	assert.Equal(t, int64(7500), usage.EmailsRemaining)
	assert.Equal(t, 25.5, usage.PercentUsed)
}

func TestUsageService_GuardedWithoutSecretKey(t *testing.T) {
	c := newTestClient(t)

	_, err := c.Usage.Get()

	var configErr *ConfigError
	assert.ErrorAs(t, err, &configErr)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}
