package caddyyaml

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// yamlToJSON converts YAML bytes to JSON bytes.
// It recursively removes all entries with "x-" prefix before marshaling to JSON.
func yamlToJSON(b []byte) ([]byte, error) {
	var tmp map[string]any

	if err := yaml.Unmarshal(b, &tmp); err != nil {
		return nil, err
	}

	// Recursively discard all entries with x- prefix
	removeExtensions(tmp)

	return json.Marshal(tmp)
}
