package telegraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_upload(t *testing.T) {

	const body = "test file body"

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if files := len(r.MultipartForm.File); files != 1 {
			http.Error(rw, fmt.Sprintf("unexpected number of files: %d", files), http.StatusBadRequest)
			return
		}

		var ur UploadResult

		for _, headers := range r.MultipartForm.File {
			for _, hdr := range headers {
				ur = append(ur, File{Src: "/file/" + hdr.Filename})
				rc, err := hdr.Open()
				if err != nil {
					http.Error(rw, fmt.Sprintf("error opening file %q : %s", hdr.Filename, err), http.StatusBadRequest)
				}
				data, err := ioutil.ReadAll(rc)
				rc.Close()
				if err != nil {
					http.Error(rw, fmt.Sprintf("failed to read the body: %s", err), http.StatusBadRequest)
					return
				}
				if !strings.EqualFold(body, string(data)) {
					http.Error(rw, fmt.Sprintf("body doesn't match: want=%q, got=%q", body, string(data)), http.StatusBadRequest)
					return
				}
			}
		}
		if len(ur) != 1 {
			http.Error(rw, fmt.Sprintf("expected one file, but got: %d", len(ur)), http.StatusBadRequest)
			return
		}

		rw.Header().Add("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(&ur); err != nil {
			t.Fatalf("server error sending the response: %s", err)
		}
	}))
	defer srv.Close()

	res, err := upload(context.Background(), srv.URL, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("unexpected result length: %d", len(res))
	}
}
