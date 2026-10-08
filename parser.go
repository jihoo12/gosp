package gosp

import (
	"fmt"
	"strconv"
	"strings"
)

func parse(filename, source string) (page, error) {
	var p page
	for len(source) > 0 {
		start := strings.Index(source, "<%")
		if start < 0 {
			p.fragments = append(p.fragments, fragment{literal, source})
			break
		}
		if start > 0 {
			p.fragments = append(p.fragments, fragment{literal, source[:start]})
		}
		source = source[start:]
		if strings.HasPrefix(source, "<%%") {
			p.fragments = append(p.fragments, fragment{literal, "<%"})
			source = source[3:]
			continue
		}
		if strings.HasPrefix(source, "<%--") {
			end := strings.Index(source[4:], "--%>")
			if end < 0 {
				return p, fmt.Errorf("%s: unclosed GOSP comment", filename)
			}
			source = source[4+end+4:]
			continue
		}
		end := strings.Index(source[2:], "%>")
		if end < 0 {
			return p, fmt.Errorf("%s: unclosed GOSP tag", filename)
		}
		content := strings.TrimSpace(source[2 : 2+end])
		source = source[2+end+2:]
		switch {
		case strings.HasPrefix(content, "="):
			expression := strings.TrimSpace(content[1:])
			if expression == "" {
				return p, fmt.Errorf("%s: empty escaped expression", filename)
			}
			p.fragments = append(p.fragments, fragment{escapedExpression, expression})
			p.hasEscaped = true
		case strings.HasPrefix(content, "-"):
			expression := strings.TrimSpace(content[1:])
			if expression == "" {
				return p, fmt.Errorf("%s: empty raw expression", filename)
			}
			p.fragments = append(p.fragments, fragment{rawExpression, expression})
			p.hasRaw = true
		case strings.HasPrefix(content, "@"):
			directive := strings.Fields(strings.TrimSpace(content[1:]))
			if len(directive) != 2 || directive[0] != "import" {
				return p, fmt.Errorf("%s: expected directive: <%%@ import \"package/path\" %%>", filename)
			}
			path, err := strconv.Unquote(directive[1])
			if err != nil || path == "" || strings.ContainsAny(path, " \n\t\r") {
				return p, fmt.Errorf("%s: invalid import path %q", filename, directive[1])
			}
			p.imports = append(p.imports, path)
		default:
			if content != "" {
				p.fragments = append(p.fragments, fragment{statement, content})
			}
		}
	}
	return p, nil
}
