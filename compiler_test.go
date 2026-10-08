package gosp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	dir := t.TempDir()
	for path, content := range map[string]string{
		"index.gosp":       `<%@ import "time" %>Hello <%= r.URL.Query().Get("name") %>!<%-- hidden --%><% for i:=0;i<2;i++ { %><%- i %><% } %>`,
		"admin/index.gosp": `Admin`,
		"admin/users.gosp": `Users`,
	} {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// The test's import directive needs to be used by the page to compile.
	if err := os.WriteFile(filepath.Join(dir, "index.gosp"), []byte(`<%@ import "time" %><% _ = time.Now() %>Hello <%= r.URL.Query().Get("name") %>!<%-- hidden --%><% for i:=0;i<2;i++ { %><%- i %><% } %>`), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := Generate(dir, "main")
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"func RenderIndex(", "func RenderAdminIndex(", "func RenderAdminUsers(",
		`mux.HandleFunc("/{$}", RenderIndex)`,
		`mux.HandleFunc("/admin/{$}", RenderAdminIndex)`,
		`mux.HandleFunc("/admin/users", RenderAdminUsers)`,
		`html.EscapeString(fmt.Sprint(r.URL.Query().Get("name")))`,
		`_gospOutput.WriteString(fmt.Sprint(i))`,
		`"time"`,
	} {
		if !strings.Contains(string(result), expected) {
			t.Errorf("generated source missing %q:\n%s", expected, result)
		}
	}
	if strings.Contains(string(result), "hidden") {
		t.Error("GOSP comments should be omitted from Go output")
	}
}

func TestGenerateErrors(t *testing.T) {
	for _, tc := range []struct{ name, text, want string }{
		{"unterminated-tag", "hi <% x", "unclosed GOSP tag"},
		{"unterminated-comment", "hi <%-- x", "unclosed GOSP comment"},
		{"empty-expression", "<%= %>", "empty escaped expression"},
		{"empty-raw", "<%- %>", "empty raw expression"},
		{"unknown-directive", `<%@ page foo %>`, "expected directive"},
		{"bad-import", `<%@ import nope %>`, "invalid import path"},
		{"invalid-Go", `<% if true { %>`, "invalid generated Go code"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "index.gosp"), []byte(tc.text), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := Generate(dir, "main")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("expected %q, got %v", tc.want, err)
			}
		})
	}
}

func TestGenerateCollision(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"foo-bar.gosp", "foo_bar.gosp"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("Hello"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Generate(dir, "main"); err == nil || !strings.Contains(err.Error(), "both generate") {
		t.Fatalf("expected a function-name collision, got %v", err)
	}
}

func TestGenerateNoFiles(t *testing.T) {
	if _, err := Generate(t.TempDir(), "main"); err == nil {
		t.Fatal("expected no files error")
	}
}

func TestGenerateInvalidPackage(t *testing.T) {
	if _, err := Generate(t.TempDir(), "for"); err == nil {
		t.Fatal("expected invalid package name")
	}
}

func TestEscapedDelimiter(t *testing.T) {
	p, err := parse("test.gosp", `before <%% after`)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.fragments) != 3 || p.fragments[1].value != "<%" {
		t.Fatalf("escaped delimiter not parsed: %#v", p.fragments)
	}
}
