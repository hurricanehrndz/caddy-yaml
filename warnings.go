package caddyyaml

import "github.com/caddyserver/caddy/v2/caddyconfig"

// warningsCollector collects warnings during YAML processing.
type warningsCollector struct {
	filename string
	warnings []caddyconfig.Warning
}

// Add adds a warning to the collector with file, line, directive, and message information.
func (w *warningsCollector) Add(line int, directive string, message string) {
	w.warnings = append(w.warnings, caddyconfig.Warning{
		File:      w.filename,
		Line:      line,
		Directive: directive,
		Message:   message,
	})
}

// newWarningsCollector creates a new warnings collector for the specified filename.
func newWarningsCollector(filename string) *warningsCollector {
	return &warningsCollector{filename, nil}
}
