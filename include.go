// Package caddyyaml provides a YAML config adapter for Caddy v2.
//
// Include processing implementation inspired by compose-go:
// https://github.com/compose-spec/compose-go
// Copyright 2020 The Compose Specification Authors
// Licensed under the Apache License, Version 2.0
package caddyyaml

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"

	"gopkg.in/yaml.v3"
)

// includeConfig represents an include directive in the YAML config.
type includeConfig struct {
	Path []string `yaml:"path"`
}

// processIncludes handles include directives by loading and merging referenced files.
// It detects circular dependencies and returns the merged configuration bytes.
func processIncludes(body []byte, baseDir string, included []string) ([]byte, error) {
	var config map[string]any
	if err := yaml.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML for includes: %w", err)
	}

	// Check if there are any includes
	includeValue, hasInclude := config["include"]
	if !hasInclude {
		return body, nil
	}

	// Parse include configurations
	includes, err := loadIncludeConfig(includeValue)
	if err != nil {
		return nil, err
	}

	// Remove include directive from config
	delete(config, "include")

	// Process each include
	for _, inc := range includes {
		for _, path := range inc.Path {
			if err := processIncludeStatements(path, baseDir, included, config); err != nil {
				return nil, err
			}
		}
	}

	// Marshal back to YAML
	return yaml.Marshal(config)
}

// processIncludeStatements loads, processes, and merges a single include file or directory into the config.
// If path is a directory, all .yaml and .yml files in the directory are processed.
func processIncludeStatements(path, baseDir string, included []string, config map[string]any) error {
	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}

	// Check if path is a directory
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	if info.IsDir() {
		return processIncludeDir(path, included, config)
	}

	return processIncludeSingleFile(path, included, config)
}

// processIncludeDir recursively processes all YAML files in a directory and its subdirectories.
func processIncludeDir(dirPath string, included []string, config map[string]any) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Process entries in sorted order for deterministic results
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		err := processIncludeDirEntry(entry, fullPath, included, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// processIncludeDirEntry processes a single directory entry (file or subdirectory).
func processIncludeDirEntry(entry os.DirEntry, fullPath string, included []string, config map[string]any) error {
	if entry.IsDir() {
		// Recursively process subdirectories
		return processIncludeDir(fullPath, included, config)
	}

	// Only process .yaml and .yml files
	ext := filepath.Ext(entry.Name())
	if ext != ".yaml" && ext != ".yml" {
		return nil
	}

	return processIncludeSingleFile(fullPath, included, config)
}

// processIncludeSingleFile loads, processes, and merges a single include file into the config.
func processIncludeSingleFile(path string, included []string, config map[string]any) error {
	// Check for circular includes
	if slices.Contains(included, path) {
		return fmt.Errorf("circular include detected: %s", path)
	}

	// Read included file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read include file %s: %w", path, err)
	}

	// Recursively process includes in the included file
	includedDir := filepath.Dir(path)
	newIncluded := append(included, path)
	content, err = processIncludes(content, includedDir, newIncluded)
	if err != nil {
		return err
	}

	// Parse included content
	var includedConfig map[string]any
	if err := yaml.Unmarshal(content, &includedConfig); err != nil {
		return fmt.Errorf("failed to parse included file %s: %w", path, err)
	}

	// Merge included config into main config
	if err := mergeConfig(config, includedConfig); err != nil {
		return fmt.Errorf("failed to merge include %s: %w", path, err)
	}

	return nil
}

// loadIncludeConfig parses the include configuration from raw YAML.
func loadIncludeConfig(source any) ([]includeConfig, error) {
	if source == nil {
		return nil, nil
	}

	// Handle both string and list formats
	configs, ok := source.([]any)
	if !ok {
		return nil, fmt.Errorf("include must be a list, got %T", source)
	}

	result := make([]includeConfig, 0, len(configs))
	for i, config := range configs {
		// Convert string shorthand to map format
		if v, ok := config.(string); ok {
			result = append(result, includeConfig{
				Path: []string{v},
			})
			continue
		}

		// Parse map format
		inc, err := parseIncludeMap(config, i)
		if err != nil {
			return nil, err
		}

		result = append(result, inc)
	}

	return result, nil
}

// parseIncludeMap parses a map-format include configuration.
func parseIncludeMap(config any, index int) (includeConfig, error) {
	configMap, ok := config.(map[string]any)
	if !ok {
		return includeConfig{}, fmt.Errorf("include[%d] must be a string or map, got %T", index, config)
	}

	var inc includeConfig
	pathValue, exists := configMap["path"]
	if !exists {
		return includeConfig{}, fmt.Errorf("include[%d] missing required 'path' field", index)
	}

	switch p := pathValue.(type) {
	case string:
		inc.Path = []string{p}
	case []any:
		inc.Path = make([]string, len(p))
		for j, pathItem := range p {
			pathStr, ok := pathItem.(string)
			if !ok {
				return includeConfig{}, fmt.Errorf("include[%d].path[%d] must be a string", index, j)
			}
			inc.Path[j] = pathStr
		}
	default:
		return includeConfig{}, fmt.Errorf("include[%d].path must be a string or list of strings", index)
	}

	return inc, nil
}

// mergeConfig deep merges source into target.
// For conflicting keys, source takes precedence.
func mergeConfig(target, source map[string]any) error {
	for key, sourceValue := range source {
		targetValue, exists := target[key]
		if !exists {
			target[key] = sourceValue
			continue
		}

		if err := mergeValue(target, key, targetValue, sourceValue); err != nil {
			return err
		}
	}

	return nil
}

// mergeValue merges a single value into the target map at the specified key.
func mergeValue(target map[string]any, key string, targetValue, sourceValue any) error {
	// If both are maps, merge recursively
	sourceMap, sourceIsMap := sourceValue.(map[string]any)
	targetMap, targetIsMap := targetValue.(map[string]any)

	if sourceIsMap && targetIsMap {
		return mergeConfig(targetMap, sourceMap)
	}

	// If both are arrays, concatenate them
	sourceSlice, sourceIsSlice := sourceValue.([]any)
	targetSlice, targetIsSlice := targetValue.([]any)

	if sourceIsSlice && targetIsSlice {
		target[key] = append(targetSlice, sourceSlice...)
		return nil
	}

	// Check for conflicts (different types or non-map values)
	if !reflect.DeepEqual(sourceValue, targetValue) {
		return fmt.Errorf("conflict at key %q: cannot merge %T with %T", key, targetValue, sourceValue)
	}

	return nil
}
