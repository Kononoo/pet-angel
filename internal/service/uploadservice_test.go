package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	uploadv1 "pet-angel/api/upload/v1"
	"pet-angel/internal/conf"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func TestUploadFile_Image(t *testing.T) {
	tmp := t.TempDir()
	svc := NewUploadService(&conf.Storage{LocalRoot: tmp, PublicPrefix: "/static/"}, nil)
	srv := khttp.NewServer()
	uploadv1.RegisterUploadServiceHTTPServer(srv, svc)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	// build a small png header bytes (just to satisfy MIME)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	_ = mw.WriteField("type", "image")
	fw, _ := mw.CreateFormFile("file", "test.png")
	// minimal PNG signature bytes
	fw.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/upload/file", buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("/static/image/")) {
		t.Fatalf("unexpected response: %s", string(body))
	}
	// ensure file written
	// find any file under tmp/image/*
	matches, _ := filepath.Glob(filepath.Join(tmp, "image", "*", "*", "*", "*"))
	if len(matches) == 0 {
		// it may be deeper path YYYY/MM/DD
		var found bool
		filepath.Walk(tmp, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(path) == ".png" {
				found = true
			}
			return nil
		})
		if !found {
			t.Fatal("uploaded file not found")
		}
	}
}
