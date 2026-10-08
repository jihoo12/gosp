package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderEscapesUntrustedHTML(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?name=%3Cscript%3Ealert(1)%3C%2Fscript%3E", nil)
	response := httptest.NewRecorder()
	RenderIndex(response, request)
	body := response.Body.String()
	if strings.Contains(body, "<script>") {
		t.Fatal("expression emitted unescaped HTML")
	}
	if !strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("expected HTML-escaped name; got %s", body)
	}
	if strings.Count(body, "<li>") != 3 {
		t.Fatalf("expected 3 loop items; got %s", body)
	}
	if ct := response.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("unexpected content type: %q", ct)
	}
}

func TestDefaultNameAndRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPages(mux)
	for _, tc := range []struct {
		path, contains string
		status         int
	}{
		{"/", "Hello, World!", http.StatusOK},
		{"/about", "<strong>Go-powered</strong>", http.StatusOK},
		{"/missing", "404 page not found", http.StatusNotFound},
	} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.contains) {
			t.Errorf("%s: status=%d body=%s", tc.path, response.Code, response.Body.String())
		}
	}
}
