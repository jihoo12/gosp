# GOSP

GOSP is a small **JSP-inspired template compiler for Go**. It turns `.gosp` pages containing HTML and Go scriptlets into ordinary Go HTTP handlers. No interpreter or extra runtime dependencies are needed.

> **Important:** `.gosp` pages contain executable Go code, so compile only templates you trust. Unlike Java JSP, GOSP compiles pages **at build time**, not automatically on each request.

## Quick start

Using the included demo:

```sh
go generate ./examples/hello
go run ./examples/hello
```

Open <http://localhost:8080/?name=Go> or <http://localhost:8080/about>.

To use GOSP in your own app:

1. Add `index.gosp` and other `.gosp` pages next to your `main.go`.
2. Generate Go handlers from the **repository root** or from your module as appropriate:

   ```sh
   go run ./cmd/gosp -src ./templates -out ./pages_gen.go -pkg main
   ```

3. Register generated handlers on a Go HTTP server:

   ```go
   mux := http.NewServeMux()
   RegisterPages(mux)
   http.ListenAndServe(":8080", mux)
   ```

The generator recursively reads `.gosp` files and generates `RenderIndex(w, r)`, `RenderAbout(w, r)`, etc. The parameters are `http.ResponseWriter` and `*http.Request`. For other Go modules, install GOSP or add it as a dependency, then run its `cmd/gosp` command.

## Template syntax

```html
<%@ import "strings" %>
<% name := strings.TrimSpace(r.URL.Query().Get("name")) %>
<h1>Hello, <%= name %>!</h1>
<ul>
  <% for i := 1; i <= 3; i++ { %>
    <li>Item <%= i %></li>
  <% } %>
</ul>
<%-- server-only comment --%>
<%- "<strong>Trusted HTML only</strong>" %>
```

| Syntax | Behavior |
| --- | --- |
| Plain HTML | Copied verbatim |
| `<% statements %>` | Runs ordinary Go statements; braces can span multiple tags |
| `<%= expression %>` | Outputs the expression with **HTML escaping** |
| `<%- expression %>` | Outputs the expression **without escaping** (trusted HTML only) |
| `<%-- comment --%>` | Removes a server-side comment |
| `<%@ import "pkg/path" %>` | Adds a Go import to the generated file |
| `<%%` | Emits the literal characters `<%` |

Expression output is escaped using `html.EscapeString`. **This is not context-aware escaping** for JavaScript, CSS, or URLs; never interpolate untrusted values into those contexts. For raw output (`<%- ... %>`), the page author is responsible for safety. There is no user-supplied template execution sandbox.

## URL mapping

| File | HTTP route | Generated function |
| --- | --- | --- |
| `index.gosp` | `/` | `RenderIndex` |
| `about.gosp` | `/about` | `RenderAbout` |
| `admin/index.gosp` | `/admin/` | `RenderAdminIndex` |
| `admin/users.gosp` | `/admin/users` | `RenderAdminUsers` |

`RegisterPages(mux)` registers these as exact paths (directory index pages do not swallow nested routes). HTTP methods are not restricted, so forms can use GET or POST. Page paths should contain ASCII letters, digits, `-`, and `_`. Colliding routes or function names cause generation to fail.

## Development

Run `go test ./...` after `go generate ./examples/hello`. GOSP uses only Go's standard library. The committed `pages_gen.go` in the demo is generated code; regenerate it when editing its templates.
