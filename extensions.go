// Package caddyyaml provides a YAML config adapter for Caddy v2.
//
// Extension field processing inspired by compose-go:
// https://github.com/compose-spec/compose-go
// Copyright 2020 The Compose Specification Authors
// Licensed under the Apache License, Version 2.0
package caddyyaml

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var extensionLineRegexp = regexp.MustCompile(`^x\-([a-zA-Z0-9\.\_\-]+)(\s*)\:`)

// removeExtensions removes only top-level x- prefixed keys from the config.
// Nested x- fields are preserved. This follows the Docker Compose convention
// where extension fields are only meaningful at the document root.
func removeExtensions(m map[string]any) {
	for key := range m {
		// Skip only top-level extension fields (x- prefix)
		if strings.HasPrefix(key, "x-") {
			delete(m, key)
		}
	}
}

// extractRawExtensions extracts x- prefixed extension fields from the body.
func extractRawExtensions(body []byte) ([]byte, error) {
	extensions, _ := extractAllMatchingTopLevelSections(body, extensionLineRegexp)
	return extensions, nil
}

// parseExtensionVars parses YAML and extracts x- variables for templates.
// It uses line-based extraction to preserve raw YAML structure (including anchors),
// then applies template processing to the extension fields themselves.
func parseExtensionVars(body []byte, env []string) (map[string]any, error) {
	// Extract raw x- field lines (preserves YAML anchors and structure)
	varsBytes, err := extractRawExtensions(body)
	if err != nil {
		return nil, err
	}

	// Apply templates to x- fields using only env vars
	// This allows x- fields to reference environment variables
	varsBytes, err = applyTemplate(varsBytes, nil, env, nil)
	if err != nil {
		return nil, err
	}

	// Parse the processed x- fields
	var tmp map[string]any
	if err := yaml.Unmarshal(varsBytes, &tmp); err != nil {
		return nil, err
	}

	// Create vars map with x- prefix removed
	// Convert hyphens to underscores for template variable names
	vars := make(map[string]any)
	for xkey, val := range tmp {
		key := xkey[2:] // Remove x- prefix
		// Replace hyphens with underscores for template compatibility
		key = strings.ReplaceAll(key, "-", "_")
		vars[key] = val
	}

	return vars, nil
}
