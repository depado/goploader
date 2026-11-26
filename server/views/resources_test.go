package views

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
)

func TestSetFileHeaders(t *testing.T) {
	cases := []struct {
		name, ctype, disp string
	}{
		{"xss.html", "application/octet-stream", "attachment"},
		{"PAGE.HTML", "application/octet-stream", "attachment"},
		{"pic.svg", "application/octet-stream", "attachment"},
		{"app.js", "application/octet-stream", "attachment"},
		{"payload", "application/octet-stream", "attachment"},
		{"note.txt", "text/plain; charset=utf-8", "filename"},
		{"pic.png", "image/png", "filename"},
		{"doc.pdf", "application/pdf", "filename"},
	}
	defer func() { conf.C.AlwaysDownload = false }()
	for _, tc := range cases {
		for _, dl := range []bool{false, true} {
			conf.C.AlwaysDownload = dl
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			setFileHeaders(c, tc.name)

			wantType, wantDisp := tc.ctype, tc.disp
			if dl {
				wantType, wantDisp = "application/octet-stream", "attachment"
			}
			if got := w.Header().Get("Content-Type"); got != wantType {
				t.Errorf("%s (always_download=%v): Content-Type = %q, want %q", tc.name, dl, got, wantType)
			}
			if got := w.Header().Get("Content-Disposition"); !strings.HasPrefix(got, wantDisp) {
				t.Errorf("%s (always_download=%v): Content-Disposition = %q, want prefix %q", tc.name, dl, got, wantDisp)
			}
			if got := w.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("%s (always_download=%v): X-Content-Type-Options = %q, want nosniff", tc.name, dl, got)
			}
		}
	}
}
