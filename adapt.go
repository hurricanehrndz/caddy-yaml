package caddyyaml

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

// adapt processes YAML configuration and converts it to Caddy JSON format.
// Processing pipeline:
// 1. Process includes (if present)
// 2. Extract x- variables for templates
// 3. Apply Go templates
// 4. Remove extension fields recursively
// 5. Convert to JSON
func adapt(body []byte, options map[string]any) ([]byte, []caddyconfig.Warning, error) {
	filename, ok := options["filename"].(string)
	if !ok {
		return nil, nil, errors.New("missing filename option")
	}

	env, ok := options[envOptionName].([]string)
	if !ok {
		env = os.Environ()
	}

	wc := newWarningsCollector(filename)

	// Phase 1: Process includes
	baseDir := filepath.Dir(filename)
	body, err := processIncludes(body, baseDir, []string{filename})
	if err != nil {
		return nil, wc.warnings, err
	}

	// Phase 2: Extract x- variables for templates
	vars, err := parseExtensionVars(body, env)
	if err != nil {
		return nil, wc.warnings, err
	}

	// Phase 3: Apply Go templates
	body, err = applyTemplate(body, vars, env, wc)
	if err != nil {
		return nil, wc.warnings, err
	}

	// Phase 4 & 5: Remove extensions and convert to JSON
	result, err := yamlToJSON(body)
	return result, wc.warnings, err
}
