package openapi

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// File types, as found in File.Type and accepted by ListFilesParams.Type.
const (
	FileTypeImage    = "image"
	FileTypeDocument = "document"
	FileTypeVideo    = "video"
	FileTypeAudio    = "audio"
	FileTypeOther    = "other"
)

// Sort orders accepted by ListFilesParams.Sort. A leading "-" sorts
// descending; the API default is FileSortCreatedAtDesc.
const (
	FileSortCreatedAt     = "createdAt"
	FileSortCreatedAtDesc = "-createdAt"
	FileSortName          = "name"
	FileSortNameDesc      = "-name"
	FileSortSize          = "size"
	FileSortSizeDesc      = "-size"
)

// FilesService manages the file manager: uploaded images, documents and
// media served from a public CDN URL, and the labels used to organise them.
//
// File labels are separate from contact labels (LabelsService).
type FilesService struct{ c *Client }

// File is one file of the file manager.
type File struct {
	ID string `json:"id"`
	// Name is the file name, extension included.
	Name      string `json:"name"`
	Extension string `json:"extension"`
	// Type is one of the FileType* constants.
	Type string `json:"type"`
	// Size is in bytes.
	Size int64 `json:"size"`
	// Width and Height are the image dimensions in pixels; 0 for other
	// types.
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Labels []string `json:"labels"`
	// URL is a public CDN URL. It does not change when the file is renamed
	// and stops working once the file is deleted.
	URL string `json:"url"`
	// ThumbnailURL is a 200x280 preview for images; nil for other types.
	ThumbnailURL *string `json:"thumbnailUrl"`
	// CreatedAt is a Unix millisecond timestamp.
	CreatedAt int64 `json:"createdAt"`
}

// FileLabel is a file manager label with the number of files carrying it.
type FileLabel struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// FileTypeCounts is the number of files per type.
type FileTypeCounts struct {
	Image    int `json:"image"`
	Document int `json:"document"`
	Video    int `json:"video"`
	Audio    int `json:"audio"`
	Other    int `json:"other"`
}

// FileStats summarises the file manager's content and storage, in bytes.
type FileStats struct {
	FileCount        int            `json:"fileCount"`
	ByType           FileTypeCounts `json:"byType"`
	StorageUsedBytes int64          `json:"storageUsedBytes"`
	// StorageLimitBytes and StorageRemainingBytes are nil when the plan has
	// no storage limit.
	StorageLimitBytes     *int64 `json:"storageLimitBytes"`
	StorageRemainingBytes *int64 `json:"storageRemainingBytes"`
}

// ListFilesParams holds the optional query parameters of List.
type ListFilesParams struct {
	Page    int
	PerPage int
	// Sort is one of the FileSort* constants; the API default is
	// FileSortCreatedAtDesc.
	Sort string
	// Search matches whole words of the file name, not prefixes.
	Search string
	// Type keeps only files of any of these types (FileType* constants).
	Type []string
	// Labels keeps only files carrying any of these labels.
	Labels []string
}

func (p *ListFilesParams) values() url.Values {
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
	if p.Sort != "" {
		values.Set("sort", p.Sort)
	}
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	if len(p.Type) > 0 {
		values.Set("type", strings.Join(p.Type, ","))
	}
	if len(p.Labels) > 0 {
		values.Set("labels", strings.Join(p.Labels, ","))
	}
	return values
}

// UploadFileRequest is the payload of Upload. File and Filename are required.
type UploadFileRequest struct {
	// File is the content. It is streamed, so it can come from disk or
	// memory.
	File io.Reader
	// Filename decides the stored name and extension. Allowed extensions:
	// jpg, jpeg, jfif, png, gif, webp, tiff, avif, pdf, doc, docx, xls, xlsx,
	// ppt, pptx, csv, txt, md, mp3, mp4, mov, avi. The name without extension
	// may only contain letters, digits, spaces, "-" and "_".
	Filename string
	// Labels are attached to the uploaded file.
	Labels []string
}

// UpdateFileRequest is the payload of Update. A nil field keeps its current
// value.
type UpdateFileRequest struct {
	// Name is the new file name. The extension cannot change; a name without
	// one keeps the current extension.
	Name *string `json:"name,omitempty"`
	// Labels replaces the whole label set when non-nil; an empty, non-nil
	// slice clears it.
	Labels *[]string `json:"labels,omitempty"`
}

// List returns a page of files, newest first unless Sort says otherwise.
func (s *FilesService) List(params *ListFilesParams) (*Page[File], error) {
	return paged[File](s.c, request{
		capability: "Files.List",
		method:     http.MethodGet,
		path:       "/files",
		query:      params.values(),
	})
}

// Get returns one file.
func (s *FilesService) Get(id string) (*File, error) {
	return object[File](s.c, request{
		capability: "Files.Get",
		method:     http.MethodGet,
		path:       "/files/" + url.PathEscape(id),
	})
}

// Upload stores one file in the file manager and returns it; File.URL is its
// public CDN URL. The content must match the extension, and csv, txt and md
// files must not contain HTML or script. Uploads count against the plan's
// storage quota (422 when it is used up; 413 when the file is too large).
func (s *FilesService) Upload(req *UploadFileRequest) (*File, error) {
	if req == nil || req.File == nil || strings.TrimSpace(req.Filename) == "" {
		return nil, errors.New("file upload requires a File and a Filename")
	}
	form := map[string]string{}
	if len(req.Labels) > 0 {
		for _, label := range req.Labels {
			if strings.Contains(label, ",") {
				return nil, errors.New("label " + label + " contains a comma, which the upload cannot address")
			}
		}
		form["labels"] = strings.Join(req.Labels, ",")
	}
	return object[File](s.c, request{
		capability: "Files.Upload",
		method:     http.MethodPost,
		path:       "/files",
		files:      []multipartFile{{field: "file", filename: req.Filename, reader: req.File}},
		formData:   form,
	})
}

// UploadPath uploads the file at path, named after its base name.
func (s *FilesService) UploadPath(path string, labels ...string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return s.Upload(&UploadFileRequest{File: f, Filename: filepath.Base(path), Labels: labels})
}

// Update renames a file and/or replaces its labels, and returns the updated
// file. Renaming does not change its URL. A name already used by another
// file answers 409.
func (s *FilesService) Update(id string, req *UpdateFileRequest) (*File, error) {
	if req == nil || (req.Name == nil && req.Labels == nil) {
		return nil, errors.New("file update requires a Name or Labels")
	}
	return object[File](s.c, request{
		capability: "Files.Update",
		method:     http.MethodPatch,
		path:       "/files/" + url.PathEscape(id),
		body:       req,
	})
}

// Delete permanently removes a file. The stored file is removed too, so its
// URL stops working everywhere it is used (emails, templates, landing pages).
//
// The API answering `data: false` is reported as an *Error with Code
// ErrCodeOperationFailed.
func (s *FilesService) Delete(id string) error {
	return confirm(s.c, request{
		capability: "Files.Delete",
		method:     http.MethodDelete,
		path:       "/files/" + url.PathEscape(id),
	})
}

// Stats returns the number of files per type and the storage used, limited
// and remaining, in bytes.
func (s *FilesService) Stats() (*FileStats, error) {
	return object[FileStats](s.c, request{
		capability: "Files.Stats",
		method:     http.MethodGet,
		path:       "/files/stats",
	})
}

// Labels returns the file manager's labels with the number of files carrying
// each. It is not paginated, and can take up to 5 minutes to reflect changes
// made outside this API.
func (s *FilesService) Labels() ([]FileLabel, error) {
	return objects[FileLabel](s.c, request{
		capability: "Files.Labels",
		method:     http.MethodGet,
		path:       "/files/labels",
	})
}

// CreateLabel creates a file label. Labels are stored in title case ("my
// banner" becomes "My Banner") and compared case-insensitively; an existing
// label answers 409.
func (s *FilesService) CreateLabel(name string) (*FileLabel, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("file label create requires a name")
	}
	return object[FileLabel](s.c, request{
		capability: "Files.CreateLabel",
		method:     http.MethodPost,
		path:       "/files/labels",
		body:       map[string]string{"name": name},
	})
}

// RenameLabel renames a file label on the label list and on every file
// carrying it. name is matched case-insensitively; newName is stored in title
// case. A newName already used by another label answers 409.
func (s *FilesService) RenameLabel(name, newName string) (*FileLabel, error) {
	if strings.TrimSpace(newName) == "" {
		return nil, errors.New("file label rename requires a new name")
	}
	return object[FileLabel](s.c, request{
		capability: "Files.RenameLabel",
		method:     http.MethodPatch,
		path:       "/files/labels/" + url.PathEscape(name),
		body:       map[string]string{"name": newName},
	})
}

// DeleteLabel removes a file label from the label list and from every file
// carrying it. The files themselves are kept. name is matched
// case-insensitively.
//
// The API answering `data: false` is reported as an *Error with Code
// ErrCodeOperationFailed.
func (s *FilesService) DeleteLabel(name string) error {
	return confirm(s.c, request{
		capability: "Files.DeleteLabel",
		method:     http.MethodDelete,
		path:       "/files/labels/" + url.PathEscape(name),
	})
}
