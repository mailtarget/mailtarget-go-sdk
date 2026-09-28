package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CampaignsService manages email campaigns.
type CampaignsService struct{ c *Client }

// CampaignSender is the from address of a campaign.
type CampaignSender struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// CampaignRecipients selects a campaign's audience by label.
type CampaignRecipients struct {
	Labels []string `json:"labels,omitempty"`
}

// Campaign is a campaign as returned by the list endpoint.
type Campaign struct {
	ID          string              `json:"id"`
	Subject     string              `json:"subject"`
	Type        string              `json:"type"`
	Stage       string              `json:"stage"`
	EmailID     string              `json:"emailId"`
	Active      bool                `json:"active"`
	MemberCount int                 `json:"memberCount"`
	SentCount   int                 `json:"sentCount"`
	LastUpdate  string              `json:"lastUpdate"`
	Sender      *CampaignSender     `json:"sender"`
	Recipients  *CampaignRecipients `json:"recipients"`
}

// CampaignDetail is the fuller record returned when fetching one campaign.
type CampaignDetail struct {
	ID          string              `json:"id"`
	Subject     string              `json:"subject"`
	Type        string              `json:"type"`
	Stage       string              `json:"stage"`
	EmailID     string              `json:"emailId"`
	EmailType   string              `json:"emailType"`
	HTMLContent string              `json:"htmlContent"`
	Snippet     string              `json:"snippet"`
	StartType   string              `json:"startType"`
	DueDate     string              `json:"dueDate"`
	Active      bool                `json:"active"`
	MemberCount int                 `json:"memberCount"`
	SentCount   int                 `json:"sentCount"`
	LastUpdate  string              `json:"lastUpdate"`
	Sender      *CampaignSender     `json:"sender"`
	Recipients  *CampaignRecipients `json:"recipients"`
}

// CampaignAnalytics aggregates a campaign's delivery and engagement metrics.
type CampaignAnalytics struct {
	CampaignID       string  `json:"campaignId"`
	Subject          string  `json:"subject"`
	SentCount        int     `json:"sentCount"`
	DeliveredCount   int     `json:"deliveredCount"`
	OpenCount        int     `json:"openCount"`
	ClickCount       int     `json:"clickCount"`
	BounceCount      int     `json:"bounceCount"`
	ComplaintCount   int     `json:"complaintCount"`
	UnsubscribeCount int     `json:"unsubscribeCount"`
	OpenRate         float64 `json:"openRate"`
	ClickRate        float64 `json:"clickRate"`
	BounceRate       float64 `json:"bounceRate"`
}

// CampaignRecipient is one recipient of a campaign and its delivery state.
type CampaignRecipient struct {
	ContactID    string   `json:"contactId"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Firstname    string   `json:"firstname"`
	Lastname     string   `json:"lastname"`
	Status       string   `json:"status"`
	Labels       []string `json:"labels"`
	DeliveredAt  string   `json:"deliveredAt"`
	BouncedAt    string   `json:"bouncedAt"`
	BounceReason string   `json:"bounceReason"`
	FirstVisited string   `json:"firstVisited"`
	LastVisited  string   `json:"lastVisited"`
	Peak         string   `json:"peak"`
	VisitCount   int      `json:"visitCount"`
}

// CampaignRequest is the payload of Create and Update. The API uses the same
// shape for both, and only the fields you set are sent.
type CampaignRequest struct {
	Subject     string              `json:"subject,omitempty"`
	Type        string              `json:"type,omitempty"`
	EmailType   string              `json:"emailType,omitempty"`
	HTMLContent string              `json:"htmlContent,omitempty"`
	Snippet     string              `json:"snippet,omitempty"`
	StartType   string              `json:"startType,omitempty"`
	DueDate     string              `json:"dueDate,omitempty"`
	TemplateID  string              `json:"templateId,omitempty"`
	Sender      *CampaignSender     `json:"sender,omitempty"`
	Recipients  *CampaignRecipients `json:"recipients,omitempty"`
}

// ListCampaignsParams holds the optional query parameters of List.
type ListCampaignsParams struct {
	Page    int
	PerPage int
	// Search matches against the campaign subject.
	Search string
}

func (p *ListCampaignsParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	return values
}

// ListCampaignRecipientsParams holds the parameters of Recipients. Status is
// required by the API.
type ListCampaignRecipientsParams struct {
	// Status filters by delivery state: delivered, opened, clicked, bounced or
	// complained.
	Status  string
	Page    int
	PerPage int
}

func (p *ListCampaignRecipientsParams) values() url.Values {
	values := url.Values{}
	if p == nil {
		return values
	}
	values.Set("status", p.Status)
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("perPage", strconv.Itoa(p.PerPage))
	}
	return values
}

func (p *ListCampaignRecipientsParams) validate() error {
	if p == nil || p.Status == "" {
		return errors.New("campaign recipients requires a Status filter " +
			"(delivered, opened, clicked, bounced or complained)")
	}
	return nil
}

// List returns a page of campaigns.
func (s *CampaignsService) List(params *ListCampaignsParams) (*Page[Campaign], error) {
	return paged[Campaign](s.c, request{
		capability: "Campaigns.List",
		method:     http.MethodGet,
		path:       "/campaigns",
		query:      params.values(),
	})
}

// Get returns one campaign by ID.
func (s *CampaignsService) Get(id string) (*CampaignDetail, error) {
	return object[CampaignDetail](s.c, request{
		capability: "Campaigns.Get",
		method:     http.MethodGet,
		path:       "/campaigns/" + url.PathEscape(id),
	})
}

// Create adds a campaign. The sender email must be a verified, permitted
// sender.
func (s *CampaignsService) Create(req *CampaignRequest) (*CampaignDetail, error) {
	return object[CampaignDetail](s.c, request{
		capability: "Campaigns.Create",
		method:     http.MethodPost,
		path:       "/campaigns",
		body:       req,
	})
}

// Update changes an existing campaign.
func (s *CampaignsService) Update(id string, req *CampaignRequest) (*CampaignDetail, error) {
	return object[CampaignDetail](s.c, request{
		capability: "Campaigns.Update",
		method:     http.MethodPut,
		path:       "/campaigns/" + url.PathEscape(id),
		body:       req,
	})
}

// Delete permanently removes a campaign.
func (s *CampaignsService) Delete(id string) error {
	return s.c.do(request{
		capability: "Campaigns.Delete",
		method:     http.MethodDelete,
		path:       "/campaigns/" + url.PathEscape(id),
	})
}

// Analytics returns a campaign's delivery and engagement metrics.
func (s *CampaignsService) Analytics(id string) (*CampaignAnalytics, error) {
	return object[CampaignAnalytics](s.c, request{
		capability: "Campaigns.Analytics",
		method:     http.MethodGet,
		path:       "/campaigns/" + url.PathEscape(id) + "/analytics",
	})
}

// Recipients returns a page of a campaign's recipients filtered by delivery
// status.
func (s *CampaignsService) Recipients(id string, params *ListCampaignRecipientsParams) (*Page[CampaignRecipient], error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	return paged[CampaignRecipient](s.c, request{
		capability: "Campaigns.Recipients",
		method:     http.MethodGet,
		path:       "/campaigns/" + url.PathEscape(id) + "/recipients",
		query:      params.values(),
	})
}

// Send triggers immediate sending of a campaign.
func (s *CampaignsService) Send(id string) error {
	return s.c.do(request{
		capability: "Campaigns.Send",
		method:     http.MethodPost,
		path:       "/campaigns/" + url.PathEscape(id) + "/send",
	})
}

// SendTest sends a campaign preview to a single address.
func (s *CampaignsService) SendTest(id, recipient string) error {
	return s.c.do(request{
		capability: "Campaigns.SendTest",
		method:     http.MethodPost,
		path:       "/campaigns/" + url.PathEscape(id) + "/send-test",
		body:       map[string]string{"recipient": recipient},
	})
}

// SetSchedule schedules a campaign to send at dueDate and returns the updated
// campaign.
func (s *CampaignsService) SetSchedule(id, dueDate string) (*CampaignDetail, error) {
	if strings.TrimSpace(dueDate) == "" {
		return nil, errors.New("campaign set-schedule requires a dueDate")
	}
	return object[CampaignDetail](s.c, request{
		capability: "Campaigns.SetSchedule",
		method:     http.MethodPost,
		path:       "/campaigns/" + url.PathEscape(id) + "/set-schedule",
		body:       map[string]string{"dueDate": dueDate},
	})
}

// CancelSchedule removes a campaign's schedule so it does not send.
func (s *CampaignsService) CancelSchedule(id string) error {
	return s.c.do(request{
		capability: "Campaigns.CancelSchedule",
		method:     http.MethodPost,
		path:       "/campaigns/" + url.PathEscape(id) + "/cancel-schedule",
	})
}
