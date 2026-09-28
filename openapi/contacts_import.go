package openapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ContactField describes one contact field, built-in or custom. Use it to
// learn which field names a CSV import or a contact filter can refer to.
type ContactField struct {
	FieldName   string `json:"fieldName"`
	DisplayName string `json:"displayName"`
	// DisplayType is one of EMAIL, TEXT, NUMBER, DATE or CHOICE.
	DisplayType string `json:"displayType"`
	// FieldGroup is DEFAULT_MAIN, DEFAULT_SUB or CUSTOM.
	FieldGroup string `json:"fieldGroup"`
	Category   string `json:"category"`
	// Values lists the allowed options of a CHOICE field (case sensitive).
	Values []any `json:"values"`
}

// ContactImport is one CSV import job.
type ContactImport struct {
	ID         string   `json:"id"`
	Filename   string   `json:"filename"`
	ImportMode string   `json:"importMode"`
	Labels     []string `json:"labels"`
	// State is READY, RUNNING, FINISH or INTERUPTED (the API's spelling).
	// PREPARING appears only transiently.
	State string `json:"state"`
	// Message carries the failure reason; it is set only when INTERUPTED.
	Message string `json:"message"`
	// Count is the number of contacts detected in the file.
	Count int `json:"count"`
	// Progress is indicative only; it may read 0 even after FINISH.
	Progress       int    `json:"progress"`
	NewContact     int    `json:"newContact"`
	UpdatedContact int    `json:"updatedContact"`
	PendingContact int    `json:"pendingContact"`
	BounceContact  int    `json:"bounceContact"`
	UnsubContact   int    `json:"unsubContact"`
	Duplicate      int    `json:"duplicate"`
	Invalid        int    `json:"invalid"`
	CreatedAt      string `json:"createdAt"`
	RunningAt      string `json:"runningAt"`
	ProgressAt     string `json:"progressAt"`
}

// ContactImportQueue is the import activity currently in progress.
type ContactImportQueue struct {
	StillImporting bool            `json:"stillImporting"`
	Imports        []ContactImport `json:"imports"`
}

// ContactImportCancel reports the outcome of cancelling an import. State is
// DELETED when the job is gone; cancelling a RUNNING job is best effort and
// may still end INTERUPTED or FINISH.
type ContactImportCancel struct {
	ID        string `json:"id"`
	Cancelled bool   `json:"cancelled"`
	State     string `json:"state"`
}

// ImportIssue is one problem the pre-flight check found in the CSV.
type ImportIssue struct {
	Line    int    `json:"line"`
	Column  string `json:"column"`
	Value   string `json:"value"`
	Reason  string `json:"reason"`
	Outcome string `json:"outcome"`
}

// ImportReport is the pre-flight validation of an uploaded CSV.
type ImportReport struct {
	// TotalRows counts every data record, comments included.
	TotalRows int `json:"totalRows"`
	// EstimatedValid predicts how many rows will be imported. It is an
	// estimate: the import's own per-row email check is more permissive.
	EstimatedValid int            `json:"estimatedValid"`
	IgnoredColumns []string       `json:"ignoredColumns"`
	IssueCounts    map[string]int `json:"issueCounts"`
	// Issues holds the first findings in file order; IssuesTruncated says
	// whether there were more.
	Issues          []ImportIssue `json:"issues"`
	IssuesTruncated bool          `json:"issuesTruncated"`
}

// ImportContactsResult is the outcome of starting an import: the queued job
// and the validation report of the file.
type ImportContactsResult struct {
	Import     ContactImport
	Validation *ImportReport
}

// contactImportCreateEnvelope matches the 202 body, which carries a
// `validation` block next to `data`.
type contactImportCreateEnvelope struct {
	Data       ContactImport `json:"data"`
	Validation *ImportReport `json:"validation"`
}

// Import modes accepted by ImportContactsRequest.ImportMode.
const (
	ImportModeReplaceAll   = "REPLACE_ALL"
	ImportModeSkipExisting = "SKIP_EXISTING"
	ImportModeFillEmpty    = "FILL_EMPTY"
)

// ImportContactsRequest is the payload of Import. File, Filename,
// FieldMapping and Labels are required.
type ImportContactsRequest struct {
	// File is the CSV content. It is streamed, so it can come from disk or
	// memory.
	File     io.Reader
	Filename string
	// FieldMapping maps each CSV header to a contact field name, e.g.
	// {"Email Address": "email"}. See Fields for the valid names.
	FieldMapping map[string]string
	// Labels are existing label names, matched case-insensitively.
	Labels []string
	// ImportMode defaults to ImportModeReplaceAll server side.
	ImportMode string
	// Strict rejects the whole file when any row has a validation issue.
	Strict bool
}

// ListContactImportsParams holds the optional query parameters of Imports.
type ListContactImportsParams struct {
	Page    int
	PerPage int
	// State filters by job state, e.g. RUNNING or FINISH.
	State string
	Sort  string
}

func (p *ListContactImportsParams) values() url.Values {
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
	if p.State != "" {
		values.Set("state", p.State)
	}
	if p.Sort != "" {
		values.Set("sort", p.Sort)
	}
	return values
}

// Fields returns every contact field, built-in and custom.
func (s *ContactsService) Fields() ([]ContactField, error) {
	return objects[ContactField](s.c, request{
		capability: "Contacts.Fields",
		method:     http.MethodGet,
		path:       "/contacts/fields",
	})
}

// Import uploads a CSV file of contacts. The import runs asynchronously: the
// returned job starts queued, so poll GetImport for its progress.
func (s *ContactsService) Import(req *ImportContactsRequest) (*ImportContactsResult, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	mapping, err := json.Marshal(req.FieldMapping)
	if err != nil {
		return nil, err
	}

	form := map[string]string{
		"fieldMapping": string(mapping),
		"labels":       strings.Join(req.Labels, ","),
	}
	if req.ImportMode != "" {
		form["importMode"] = req.ImportMode
	}
	if req.Strict {
		form["strict"] = "true"
	}

	out, err := bare[contactImportCreateEnvelope](s.c, request{
		capability: "Contacts.Import",
		method:     http.MethodPost,
		path:       "/contacts/import",
		files:      []multipartFile{{field: "file", filename: req.Filename, reader: req.File}},
		formData:   form,
	})
	if err != nil {
		return nil, err
	}
	return &ImportContactsResult{Import: out.Data, Validation: out.Validation}, nil
}

func (r *ImportContactsRequest) validate() error {
	if r == nil || r.File == nil || strings.TrimSpace(r.Filename) == "" {
		return errors.New("contacts import requires a File and a Filename")
	}
	if len(r.FieldMapping) == 0 {
		return errors.New("contacts import requires a FieldMapping from CSV headers to contact fields")
	}
	if len(r.Labels) == 0 {
		return errors.New("contacts import requires at least one label")
	}
	for _, label := range r.Labels {
		if strings.Contains(label, ",") {
			return errors.New("label " + label + " contains a comma, which the import cannot address")
		}
	}
	return nil
}

// Imports returns a page of past imports.
func (s *ContactsService) Imports(params *ListContactImportsParams) (*Page[ContactImport], error) {
	return paged[ContactImport](s.c, request{
		capability: "Contacts.Imports",
		method:     http.MethodGet,
		path:       "/contacts/import",
		query:      params.values(),
	})
}

// CurrentImport returns the import activity currently in progress.
func (s *ContactsService) CurrentImport() (*ContactImportQueue, error) {
	return object[ContactImportQueue](s.c, request{
		capability: "Contacts.CurrentImport",
		method:     http.MethodGet,
		path:       "/contacts/import/current",
	})
}

// GetImport returns the status of one import job.
func (s *ContactsService) GetImport(id string) (*ContactImport, error) {
	return object[ContactImport](s.c, request{
		capability: "Contacts.GetImport",
		method:     http.MethodGet,
		path:       "/contacts/import/" + url.PathEscape(id),
	})
}

// CancelImport cancels an import job.
func (s *ContactsService) CancelImport(id string) (*ContactImportCancel, error) {
	return object[ContactImportCancel](s.c, request{
		capability: "Contacts.CancelImport",
		method:     http.MethodDelete,
		path:       "/contacts/import/" + url.PathEscape(id),
	})
}
