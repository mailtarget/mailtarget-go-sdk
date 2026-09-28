package openapi

import (
	"errors"
	"net/http"
	"strings"
)

// TransmissionsService sends a transactional email through the Open API.
//
// Prefer layang.Layang.Send: it needs only the Mailtarget API key, which every
// account has, while this endpoint requires the Open API secret key as well.
// This exists so the SDK covers the whole Open API surface.
type TransmissionsService struct{ c *Client }

// Address is an email address. It mirrors layang.Address, which the
// transmission API uses; the two are separate because the Open API documents
// its own schema for this endpoint.
type Address struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// Header is a custom email header.
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Attachment is a file attached to an email. Value holds its base64 content.
type Attachment struct {
	Filename    string `json:"filename,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
	Value       string `json:"value,omitempty"`
	ContentID   string `json:"contentId,omitempty"`
	Disposition string `json:"disposition,omitempty"`
}

// OptionsAttributes toggles per-message tracking behaviour.
type OptionsAttributes struct {
	ClickTracking bool `json:"clickTracking"`
	OpenTracking  bool `json:"openTracking"`
	Transactional bool `json:"transactional"`
}

// SendEmailRequest is the payload of Send. From, Subject and To are required.
//
// APIKey is filled in by the SDK from the client's Mailtarget API key, so
// callers normally leave it empty. The Open API forwards it upstream for
// authentication and rate limiting rather than reading it from the header.
type SendEmailRequest struct {
	APIKey            string             `json:"apiKey"`
	From              *Address           `json:"from"`
	Subject           string             `json:"subject"`
	To                []Address          `json:"to"`
	CC                []Address          `json:"cc,omitempty"`
	BCC               []Address          `json:"bcc,omitempty"`
	ReplyTo           []Address          `json:"replyTo,omitempty"`
	BodyText          string             `json:"bodyText,omitempty"`
	BodyHTML          string             `json:"bodyHtml,omitempty"`
	Headers           []Header           `json:"headers,omitempty"`
	Attachments       []Attachment       `json:"attachments,omitempty"`
	Metadata          map[string]string  `json:"metadata,omitempty"`
	SubstitutionData  map[string]string  `json:"substitutionData,omitempty"`
	TemplateID        string             `json:"templateId,omitempty"`
	AllRCPTTo         []string           `json:"allRCPTto,omitempty"`
	OptionsAttributes *OptionsAttributes `json:"optionsAttributes,omitempty"`
}

// validate enforces the fields the Open API marks as required, so a payload
// that cannot succeed is rejected here instead of costing a round trip.
func (r *SendEmailRequest) validate() error {
	if r.From == nil || strings.TrimSpace(r.From.Email) == "" {
		return errors.New("transmissions send requires From with an email address")
	}
	if strings.TrimSpace(r.Subject) == "" {
		return errors.New("transmissions send requires a Subject")
	}
	if len(r.To) == 0 {
		return errors.New("transmissions send requires at least one To address")
	}
	return nil
}

// SendEmailResult identifies the accepted transmission.
type SendEmailResult struct {
	TransmissionID string `json:"transmissionId"`
}

// Send delivers one transactional email through the Open API.
func (s *TransmissionsService) Send(req *SendEmailRequest) (*SendEmailResult, error) {
	if req == nil {
		return nil, errors.New("transmissions send requires a request")
	}
	if err := req.validate(); err != nil {
		return nil, err
	}

	// Copy so filling in the API key never mutates the caller's value.
	payload := *req
	if payload.APIKey == "" {
		payload.APIKey = s.c.transmissionAPIKey
	}
	if payload.APIKey == "" {
		return nil, errors.New("transmissions send requires a Mailtarget API key — " +
			"build the client with layang.NewMailtargetClient(apiKey, ...), or set " +
			"openapi.WithTransmissionAPIKey(apiKey), or fill SendEmailRequest.APIKey")
	}

	return object[SendEmailResult](s.c, request{
		capability: "Transmissions.Send",
		method:     http.MethodPost,
		path:       "/transmissions",
		body:       &payload,
	})
}
