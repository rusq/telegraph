// Package telegraph provides some functions to interface with APIs of telegra.ph.
package telegraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	BaseURL = "https://telegra.ph"

	uploadEP = BaseURL + "/upload"
)

const (
	ulKeyName     = "name"
	ulKeyFilename = "filename"
)

// File is the relative path to the file on the telegra.ph service.
type File struct {
	Src string `json:"src"`
}

// UploadResult is the result of file upload request.
type UploadResult []File

// Upload uploads the file to telegra.ph service.
func Upload(ctx context.Context, r io.Reader) (UploadResult, error) {
	return upload(ctx, uploadEP, r)
}

// upload is the actual uploader with mockable endpoint.
func upload(ctx context.Context, ep string, r io.Reader) (UploadResult, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("blob", "filename")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, r); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: %s, message: %q", resp.Status, string(data))
	}

	var result UploadResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
