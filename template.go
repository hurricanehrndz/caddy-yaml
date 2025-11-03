package caddyyaml

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

const (
	openingDelim = "#{"
	closingDelim = "}"
)

// applyTemplate processes the YAML body as a Go template with sprig functions.
// It prepends environment variables as template variables and executes the template with the provided values.
// Returns the processed template output or an error if template parsing or execution fails.
func applyTemplate(body []byte, values map[string]any, env []string, wc *warningsCollector) ([]byte, error) {
	tplBody := envVarsTemplate(env, wc) + string(body)

	tpl, err := template.New("yaml").
		Funcs(sprig.TxtFuncMap()).
		Delims(openingDelim, closingDelim).
		Parse(tplBody)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, values); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// envVarsTemplate generates template variable declarations from environment variables.
// It filters out environment variables with invalid identifiers and adds warnings for them.
// Returns a string containing template variable assignments for valid environment variables.
func envVarsTemplate(env []string, wc *warningsCollector) string {
	var builder strings.Builder
	line := func(key, val string) string {
		return tplWrap(fmt.Sprintf(`$%s := %q`, key, val))
	}
	for _, env := range env {
		key, val, _ := strings.Cut(env, "=")
		if !token.IsIdentifier(key) {
			if wc != nil {
				wc.Add(-1, "", fmt.Sprintf("environment variable %q cannot be used in template", key))
			}
			continue
		}
		fmt.Fprintln(&builder, line(key, val))
	}
	return builder.String()
}

// tplWrap wraps a string with template delimiters.
func tplWrap(s string) string {
	return fmt.Sprintf("%s %s %s", openingDelim, s, closingDelim)
}
