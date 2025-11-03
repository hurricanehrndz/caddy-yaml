package caddyyaml

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// yamlToJSON converts YAML bytes to JSON bytes.
// It removes all entries with "x-" prefix before marshaling to JSON.
func yamlToJSON(b []byte) ([]byte, error) {
	var tmp map[string]interface{}

	if err := yaml.Unmarshal(b, &tmp); err != nil {
		return nil, err
	}

	// discard all entries with x- prefix
	for key := range tmp {
		if strings.HasPrefix(key, "x-") {
			// this is safe to do
			// https://stackoverflow.com/a/23230406/524060
			delete(tmp, key)
		}
	}

	return json.Marshal(tmp)
}

// varsFromBody extracts template variables from x- prefixed YAML fields.
// It validates that variable names don't contain hyphens or dots (except the x- prefix).
// Returns a map of variable names to values for use in template execution.
func varsFromBody(b []byte, env []string) (map[string]interface{}, error) {
	varsBytes, err := extractVariables(b)
	if err != nil {
		return nil, err
	}

	varsBytes, err = applyTemplate(varsBytes, nil, env, nil)
	if err != nil {
		return nil, err
	}

	var tmp map[string]interface{}
	if err := yaml.Unmarshal(varsBytes, &tmp); err != nil {
		return nil, err
	}

	vars := make(map[string]interface{})
	for xkey, val := range tmp {
		key := xkey[2:] // discard x- prefix

		if strings.Contains(key, "-") {
			return nil, fmt.Errorf("template: apart from 'x-' prefix, '-' is not permitted in extension field name for %s", xkey)
		}

		if strings.Contains(key, ".") {
			return nil, fmt.Errorf("template: '.' is not permitted in extension field name for %s", xkey)
		}

		vars[key] = val
	}

	return vars, nil
}
