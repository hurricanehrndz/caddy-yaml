package caddyyaml

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name             string
		yamlFile         string
		jsonFile         string
		env              []string
		expectedWarnings []string
	}{
		{
			name:     "simple",
			yamlFile: "test.caddy.yaml",
			jsonFile: "test.caddy.json",
			env:      []string{"ENVIRONMENT=something"},
		},
		{
			name:     "production environment",
			yamlFile: "test.caddy.yaml",
			jsonFile: "test.caddy.prod.json",
			env:      []string{"ENVIRONMENT=production"},
		},
		{
			name:     "bad environment variables",
			yamlFile: "test.caddy.yaml",
			jsonFile: "test.caddy.json",
			env: []string{
				"ENVIRONMENT=bad",
				"INVALID%=invalid_name",
				"VALUE_WITH_BACKSLASH=foo\\bar",
			},
			expectedWarnings: []string{
				"./testdata/test.caddy.yaml:-1: environment variable \"INVALID%\" cannot be used in template",
			},
		},
		{
			name:     "YAML aliases",
			yamlFile: "test.aliases.caddy.yaml",
			jsonFile: "test.aliases.caddy.json",
			env:      []string{"ENVIRONMENT=production"},
		},
		{
			name:     "nested extension fields",
			yamlFile: "test.nested-extensions.yaml",
			jsonFile: "test.nested-extensions.json",
			env:      []string{"ENVIRONMENT=test"},
		},
		{
			name:     "include support",
			yamlFile: "test.include.yaml",
			jsonFile: "test.include.json",
			env:      []string{"ENVIRONMENT=test"},
		},
		{
			name:     "include directory support",
			yamlFile: "test.include-dir.yaml",
			jsonFile: "test.include-dir.json",
			env:      []string{"ENVIRONMENT=test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := os.ReadFile("./testdata/" + tt.yamlFile)
			if err != nil {
				t.Fatal(err)
			}
			adaptedBytes, warnings, err := Adapter{}.Adapt(b, map[string]any{
				"filename":    "./testdata/" + tt.yamlFile,
				envOptionName: tt.env,
			})
			if err != nil {
				t.Fatal(err)
			}

			for i, w := range warnings {
				if len(tt.expectedWarnings) < i+1 {
					t.Fatalf("unexpected warning %q", w)
				}
				eW := tt.expectedWarnings[i]
				if eW != w.String() {
					t.Fatalf("expected warning %q, got %q", eW, w)
				}
			}
			if len(tt.expectedWarnings) > len(warnings) {
				t.Fatalf("expected additional warnings: %v", tt.expectedWarnings[len(warnings):])
			}

			jsonBytes, err := os.ReadFile("./testdata/" + tt.jsonFile)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(jsonToObj(adaptedBytes), jsonToObj(jsonBytes)) {
				t.Log(string(adaptedBytes))
				t.Fatal("adapter config does not match expected config")
			}
		})
	}
}

func jsonToObj(b []byte) (obj map[string]any) {
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return obj
}
