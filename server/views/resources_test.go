package views

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/models"
)

func TestServedFileContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	orig := conf.C
	t.Cleanup(func() {
		conf.C = orig
		if database.DB != nil {
			database.DB.Close() //nolint:errcheck
			database.DB = nil
		}
	})
	conf.C = conf.NewDefault()
	conf.C.UploadDir = dir
	conf.C.DB = path.Join(dir, "test.db")
	conf.C.LogLevel = "disabled"
	conf.C.AlwaysDownload = false

	if err := database.Initialize(); err != nil {
		t.Fatal(err)
	}
	if err := database.DB.Init(&models.Resource{}); err != nil {
		t.Fatal(err)
	}

	htmlBody := []byte("<!DOCTYPE html><html><body><script>alert(1)</script></body></html>")
	textBody := []byte("hello plain text\n")
	svgBody := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`)
	jsBody := []byte("alert(document.domain)\n")
	xmlBody := []byte("<?xml version=\"1.0\"?><html><script>alert(1)</script></html>")
	pdfBody := []byte("%PDF-1.4\n% fake\n")
	key := "0123456789abcdef"
	txtType := mime.TypeByExtension(".txt")

	type tc struct {
		name       string
		id         string
		filename   string
		body       []byte
		encrypted  bool
		head       bool
		alwaysDown bool
		wantType   string
		attachment bool
	}
	cases := []tc{
		{name: "html", id: "htmlplain", filename: "xss.html", body: htmlBody, wantType: "application/octet-stream", attachment: true},
		{name: "html uppercase", id: "htmlupper", filename: "PAGE.HTML", body: htmlBody, wantType: "application/octet-stream", attachment: true},
		{name: "html no extension", id: "htmlnoext", filename: "payload", body: htmlBody, wantType: "application/octet-stream", attachment: true},
		{name: "html disguised as png", id: "htmlpng", filename: "xss.png", body: htmlBody, wantType: "image/png"},
		{name: "plain text", id: "plaintext", filename: "note.txt", body: textBody, wantType: txtType},
		{name: "svg", id: "svgfile", filename: "pic.svg", body: svgBody, wantType: "application/octet-stream", attachment: true},
		{name: "javascript", id: "jsfile", filename: "app.js", body: jsBody, wantType: "application/octet-stream", attachment: true},
		{name: "xml", id: "xmlfile", filename: "data.xml", body: xmlBody, wantType: "application/octet-stream", attachment: true},
		{name: "pdf", id: "pdffile", filename: "doc.pdf", body: pdfBody, wantType: "application/pdf"},
		{name: "html head", id: "htmlhead", filename: "xss.html", body: htmlBody, head: true, wantType: "application/octet-stream", attachment: true},
		{name: "encrypted html", id: "enchtml", filename: "xss.html", body: htmlBody, encrypted: true, wantType: "application/octet-stream", attachment: true},
		{name: "encrypted text", id: "enctxt", filename: "note.txt", body: textBody, encrypted: true, wantType: txtType},
		{name: "encrypted html disguised as png", id: "encpng", filename: "xss.png", body: htmlBody, encrypted: true, wantType: "image/png"},
		{name: "encrypted html head", id: "enchead", filename: "xss.html", body: htmlBody, encrypted: true, head: true, wantType: "application/octet-stream", attachment: true},
		{name: "always download", id: "alwaysdl", filename: "note.txt", body: textBody, alwaysDown: true, wantType: "application/octet-stream", attachment: true},
	}

	for _, tc := range cases {
		if tc.encrypted {
			writeEncrypted(t, tc.id, key, tc.body)
		} else {
			writePlain(t, tc.id, tc.body)
		}
		saveResource(t, tc.id, tc.filename, len(tc.body))
	}

	r := gin.New()
	r.GET("/p/:uniuri", View)
	r.HEAD("/p/:uniuri", Head)
	r.GET("/e/:uniuri/:key", ViewC)
	r.HEAD("/e/:uniuri/:key", HeadC)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conf.C.AlwaysDownload = tc.alwaysDown
			method := http.MethodGet
			if tc.head {
				method = http.MethodHead
			}
			url := srv.URL + "/p/" + tc.id
			if tc.encrypted {
				url = srv.URL + "/e/" + tc.id + "/" + key
			}
			req, err := http.NewRequest(method, url, nil)
			if err != nil {
				t.Fatal(err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close() //nolint:errcheck
			got, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, body = %q", res.StatusCode, got)
			}
			if ct := res.Header.Get("Content-Type"); ct != tc.wantType {
				t.Fatalf("Content-Type = %q, want %q", ct, tc.wantType)
			}
			if nosniff := res.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
				t.Fatalf("X-Content-Type-Options = %q, want nosniff", nosniff)
			}
			disp := res.Header.Get("Content-Disposition")
			if tc.attachment {
				if disp != "attachment; filename=\""+tc.filename+"\"" {
					t.Fatalf("Content-Disposition = %q", disp)
				}
			} else if disp != "filename=\""+tc.filename+"\"" {
				t.Fatalf("Content-Disposition = %q", disp)
			}
			if strings.Contains(strings.ToLower(res.Header.Get("Content-Type")), "text/html") ||
				strings.Contains(strings.ToLower(res.Header.Get("Content-Type")), "svg") {
				t.Fatalf("executable content type %q", res.Header.Get("Content-Type"))
			}
			if !tc.head && !bytes.Equal(got, tc.body) {
				t.Fatalf("body = %q", got)
			}
		})
	}
}

func writePlain(t *testing.T, id string, body []byte) {
	t.Helper()
	if err := os.WriteFile(path.Join(conf.C.UploadDir, id), body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeEncrypted(t *testing.T, id, key string, body []byte) {
	t.Helper()
	f, err := os.Create(path.Join(conf.C.UploadDir, id))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	var iv [aes.BlockSize]byte
	sw := &cipher.StreamWriter{S: cipher.NewCTR(block, iv[:]), W: f}
	if _, err := sw.Write(body); err != nil {
		t.Fatal(err)
	}
}

func saveResource(t *testing.T, id, name string, size int) {
	t.Helper()
	del := time.Now().Add(time.Hour)
	r := models.Resource{
		Key:          id,
		Name:         name,
		DeleteAt:     del,
		UnixDeleteAt: del.Unix(),
		Size:         int64(size),
	}
	if err := r.Save(); err != nil {
		t.Fatal(err)
	}
}
