package openapi

import "net/http"

// SettingsService reads the authenticated account's own settings.
type SettingsService struct{ c *Client }

// CompanyDetail is the company profile and configuration.
type CompanyDetail struct {
	CompanyID       string `json:"companyId"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Website         string `json:"website"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Country         string `json:"country"`
	Zipcode         string `json:"zipcode"`
	Industry        string `json:"industry"`
	Language        string `json:"language"`
	Timezone        string `json:"timezone"`
	Image           string `json:"image"`
	Logo            string `json:"logo"`
	Packet          string `json:"packet"`
	Role            string `json:"role"`
	Paid            bool   `json:"paid"`
	Default         bool   `json:"default"`
	AutoLinkUTM     bool   `json:"autoLinkUtm"`
	ExpiredDate     int64  `json:"expiredDate"`
	GracePeriodDate int64  `json:"gracePeriodDate"`
}

// UserProfile is the authenticated user's profile.
type UserProfile struct {
	UserID        string   `json:"userId"`
	Email         string   `json:"email"`
	Fullname      string   `json:"fullname"`
	Firstname     string   `json:"firstname"`
	Lastname      string   `json:"lastname"`
	Phone         string   `json:"phone"`
	Image         string   `json:"image"`
	Language      string   `json:"language"`
	PolicyVersion string   `json:"policyVersion"`
	Tours         []string `json:"tours"`
	Hints         any      `json:"hints"`
	Validate      bool     `json:"validate"`
	Time          int64    `json:"time"`
}

// Company returns the company profile.
//
// This endpoint answers with the object directly instead of the usual
// {"data": ...} envelope, so it decodes the whole body.
func (s *SettingsService) Company() (*CompanyDetail, error) {
	return bare[CompanyDetail](s.c, request{
		capability: "Settings.Company",
		method:     http.MethodGet,
		path:       "/settings/mtarget/company",
	})
}

// Profile returns the authenticated user's profile.
//
// Like Company, this endpoint is not enveloped.
func (s *SettingsService) Profile() (*UserProfile, error) {
	return bare[UserProfile](s.c, request{
		capability: "Settings.Profile",
		method:     http.MethodGet,
		path:       "/settings/mtarget/profile",
	})
}
