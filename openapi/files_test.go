package openapi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fileJSON = `{"id":"f1","name":"banner.png","extension":"png","type":"image","size":5000000000,
	"width":800,"height":600,"labels":["Promo"],"url":"https://cdn.test/banner.png",
	"thumbnailUrl":"https://cdn.test/banner_thumb.png","createdAt":1760000000000}`

func TestFilesService_List(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/files",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.Query().Encode()
			return httpmock.NewStringResponse(200,
				`{"data":[`+fileJSON+`],"meta":{"page":1,"perPage":20,"total":1}}`), nil
		})

	page, err := c.Files.List(&ListFilesParams{
		Page: 1, PerPage: 20, Search: "banner", Sort: FileSortSizeDesc,
		Type:   []string{FileTypeImage, FileTypeDocument},
		Labels: []string{"Promo", "Hero"},
	})

	require.NoError(t, err)
	assert.Equal(t,
		"labels=Promo%2CHero&page=1&perPage=20&search=banner&sort=-size&type=image%2Cdocument",
		gotQuery)
	require.Len(t, page.Items, 1)
	f := page.Items[0]
	assert.Equal(t, "f1", f.ID)
	assert.Equal(t, "banner.png", f.Name)
	assert.Equal(t, "png", f.Extension)
	assert.Equal(t, FileTypeImage, f.Type)
	assert.Equal(t, int64(5000000000), f.Size)
	assert.Equal(t, 800, f.Width)
	assert.Equal(t, 600, f.Height)
	assert.Equal(t, []string{"Promo"}, f.Labels)
	assert.Equal(t, "https://cdn.test/banner.png", f.URL)
	require.NotNil(t, f.ThumbnailURL)
	assert.Equal(t, "https://cdn.test/banner_thumb.png", *f.ThumbnailURL)
	assert.Equal(t, int64(1760000000000), f.CreatedAt)
	assert.Equal(t, 1, page.Meta.Total)
}

func TestFilesService_ListWithoutParamsSendsNoQuery(t *testing.T) {
	c := withSecret(t)

	var gotQuery string
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/files",
		func(req *http.Request) (*http.Response, error) {
			gotQuery = req.URL.RawQuery
			return httpmock.NewStringResponse(200, `{"data":[],"meta":{"page":1,"perPage":20,"total":0}}`), nil
		})

	_, err := c.Files.List(nil)

	require.NoError(t, err)
	assert.Empty(t, gotQuery)
}

func TestFilesService_GetWithoutThumbnail(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/files/f%2F2",
		httpmock.NewStringResponder(200, `{"data":{"id":"f/2","name":"report.pdf","extension":"pdf",
			"type":"document","size":1024,"width":0,"height":0,"labels":[],
			"url":"https://cdn.test/report.pdf","thumbnailUrl":null,"createdAt":1760000000000}}`))

	f, err := c.Files.Get("f/2")

	require.NoError(t, err)
	assert.Equal(t, FileTypeDocument, f.Type)
	assert.Nil(t, f.ThumbnailURL)
	assert.Equal(t, 0, f.Width)
}

func TestFilesService_UploadSendsMultipart(t *testing.T) {
	c := withSecret(t)

	var gotAuth, gotFile, gotFilename, gotContentType string
	var gotForm map[string][]string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/files",
		func(req *http.Request) (*http.Response, error) {
			gotAuth = req.Header.Get("Authorization")
			gotContentType = req.Header.Get("Content-Type")
			require.NoError(t, req.ParseMultipartForm(1<<20))
			gotForm = req.MultipartForm.Value
			file, header, err := req.FormFile("file")
			require.NoError(t, err)
			content, _ := io.ReadAll(file)
			gotFile, gotFilename = string(content), header.Filename
			return httpmock.NewStringResponse(201, `{"data":`+fileJSON+`}`), nil
		})

	f, err := c.Files.Upload(&UploadFileRequest{
		File: strings.NewReader("PNGDATA"), Filename: "banner.png", Labels: []string{"Promo", "Hero"},
	})

	require.NoError(t, err)
	assert.Equal(t, "Bearer "+testSecretKey, gotAuth)
	assert.True(t, strings.HasPrefix(gotContentType, "multipart/form-data"))
	assert.Equal(t, "banner.png", gotFilename)
	assert.Equal(t, "PNGDATA", gotFile)
	assert.Equal(t, []string{"Promo,Hero"}, gotForm["labels"])
	assert.Equal(t, "https://cdn.test/banner.png", f.URL)
}

func TestFilesService_UploadPathWithoutLabels(t *testing.T) {
	c := withSecret(t)
	path := filepath.Join(t.TempDir(), "notes.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello"), 0o600))

	var gotFile, gotFilename string
	var gotForm map[string][]string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/files",
		func(req *http.Request) (*http.Response, error) {
			require.NoError(t, req.ParseMultipartForm(1<<20))
			gotForm = req.MultipartForm.Value
			file, header, err := req.FormFile("file")
			require.NoError(t, err)
			content, _ := io.ReadAll(file)
			gotFile, gotFilename = string(content), header.Filename
			return httpmock.NewStringResponse(201, `{"data":`+fileJSON+`}`), nil
		})

	_, err := c.Files.UploadPath(path)

	require.NoError(t, err)
	assert.Equal(t, "notes.txt", gotFilename)
	assert.Equal(t, "hello", gotFile)
	assert.NotContains(t, gotForm, "labels")
}

func TestFilesService_UploadValidatesLocally(t *testing.T) {
	c := withSecret(t)

	_, err := c.Files.Upload(nil)
	assert.Error(t, err)
	_, err = c.Files.Upload(&UploadFileRequest{File: strings.NewReader("x"), Filename: " "})
	assert.Error(t, err)
	_, err = c.Files.Upload(&UploadFileRequest{File: strings.NewReader("x"), Filename: "a.png", Labels: []string{"a,b"}})
	assert.Error(t, err)
	_, err = c.Files.UploadPath(filepath.Join(t.TempDir(), "missing.png"))
	assert.Error(t, err)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestFilesService_Update(t *testing.T) {
	cases := map[string]struct {
		req  *UpdateFileRequest
		want string
	}{
		"name only":    {&UpdateFileRequest{Name: ptr("hero")}, `{"name":"hero"}`},
		"clear labels": {&UpdateFileRequest{Labels: &[]string{}}, `{"labels":[]}`},
		"both":         {&UpdateFileRequest{Name: ptr("hero"), Labels: &[]string{"Promo"}}, `{"name":"hero","labels":["Promo"]}`},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := withSecret(t)
			var gotBody string
			httpmock.RegisterResponder("PATCH", testBaseURL+"/v1/files/f1",
				func(req *http.Request) (*http.Response, error) {
					b, _ := io.ReadAll(req.Body)
					gotBody = string(b)
					return httpmock.NewStringResponse(200, `{"data":`+fileJSON+`}`), nil
				})

			f, err := c.Files.Update("f1", tc.req)

			require.NoError(t, err)
			assert.JSONEq(t, tc.want, gotBody)
			assert.Equal(t, "f1", f.ID)
		})
	}
}

func TestFilesService_UpdateRequiresAField(t *testing.T) {
	c := withSecret(t)

	_, err := c.Files.Update("f1", nil)
	assert.Error(t, err)
	_, err = c.Files.Update("f1", &UpdateFileRequest{})
	assert.Error(t, err)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestFilesService_Stats(t *testing.T) {
	t.Run("limited plan", func(t *testing.T) {
		c := withSecret(t)
		httpmock.RegisterResponder("GET", testBaseURL+"/v1/files/stats",
			httpmock.NewStringResponder(200, `{"data":{"fileCount":7,
				"byType":{"image":3,"document":2,"video":1,"audio":1,"other":0},
				"storageUsedBytes":3000000000,"storageLimitBytes":10000000000,"storageRemainingBytes":7000000000}}`))

		stats, err := c.Files.Stats()

		require.NoError(t, err)
		assert.Equal(t, 7, stats.FileCount)
		assert.Equal(t, FileTypeCounts{Image: 3, Document: 2, Video: 1, Audio: 1}, stats.ByType)
		assert.Equal(t, int64(3000000000), stats.StorageUsedBytes)
		require.NotNil(t, stats.StorageLimitBytes)
		assert.Equal(t, int64(10000000000), *stats.StorageLimitBytes)
		require.NotNil(t, stats.StorageRemainingBytes)
		assert.Equal(t, int64(7000000000), *stats.StorageRemainingBytes)
	})
	t.Run("unlimited plan", func(t *testing.T) {
		c := withSecret(t)
		httpmock.RegisterResponder("GET", testBaseURL+"/v1/files/stats",
			httpmock.NewStringResponder(200, `{"data":{"fileCount":0,"byType":{},
				"storageUsedBytes":0,"storageLimitBytes":null,"storageRemainingBytes":null}}`))

		stats, err := c.Files.Stats()

		require.NoError(t, err)
		assert.Nil(t, stats.StorageLimitBytes)
		assert.Nil(t, stats.StorageRemainingBytes)
	})
}

func TestFilesService_Labels(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("GET", testBaseURL+"/v1/files/labels",
		httpmock.NewStringResponder(200, `{"data":[{"name":"Promo","count":3},{"name":"Hero","count":0}]}`))

	labels, err := c.Files.Labels()

	require.NoError(t, err)
	assert.Equal(t, []FileLabel{{Name: "Promo", Count: 3}, {Name: "Hero", Count: 0}}, labels)
}

func TestFilesService_CreateLabel(t *testing.T) {
	c := withSecret(t)
	var gotBody map[string]string
	httpmock.RegisterResponder("POST", testBaseURL+"/v1/files/labels",
		func(req *http.Request) (*http.Response, error) {
			require.NoError(t, json.NewDecoder(req.Body).Decode(&gotBody))
			return httpmock.NewStringResponse(201, `{"data":{"name":"My Banner","count":0}}`), nil
		})

	label, err := c.Files.CreateLabel("my banner")

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"name": "my banner"}, gotBody)
	assert.Equal(t, "My Banner", label.Name)

	_, err = c.Files.CreateLabel(" ")
	assert.Error(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestFilesService_RenameLabelEscapesName(t *testing.T) {
	c := withSecret(t)
	var gotPath string
	var gotBody map[string]string
	httpmock.RegisterResponder("PATCH", `=~^`+testBaseURL+`/v1/files/labels/`,
		func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.EscapedPath()
			require.NoError(t, json.NewDecoder(req.Body).Decode(&gotBody))
			return httpmock.NewStringResponse(200, `{"data":{"name":"Hero","count":3}}`), nil
		})

	label, err := c.Files.RenameLabel("my/banner", "hero")

	require.NoError(t, err)
	assert.Equal(t, "/v1/files/labels/my%2Fbanner", gotPath)
	assert.Equal(t, map[string]string{"name": "hero"}, gotBody)
	assert.Equal(t, 3, label.Count)

	_, err = c.Files.RenameLabel("Hero", "")
	assert.Error(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestFilesService_DeleteLabel(t *testing.T) {
	c := withSecret(t)
	httpmock.RegisterResponder("DELETE", testBaseURL+"/v1/files/labels/My%20Banner",
		httpmock.NewStringResponder(200, `{"data":true}`))

	assert.NoError(t, c.Files.DeleteLabel("My Banner"))
}

func ptr[T any](v T) *T { return &v }
