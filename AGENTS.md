# Agent Guidelines for caddy-yaml

## Build/Test Commands
- Build: `go build -v .`
- Run all tests: `go test -v .`
- Run single test: `go test -v . -run TestApply`
- Get dependencies: `go get -v ./...`

## Code Style Guidelines

**Package & Imports:**
- Package: `package caddyyaml` (single word, lowercase)
- Group imports: stdlib first, then third-party (github.com/caddyserver, github.com/ghodss, github.com/Masterminds)
- Use blank lines to separate import groups

**Naming & Types:**
- Use camelCase for private functions/vars: `applyTemplate`, `varsFromBody`, `warningsCollector`
- Use PascalCase for exported types: `Adapter`
- Receiver names: single letter or short abbreviation (e.g., `a Adapter`, `w *warningsCollector`)

**Error Handling:**
- Return errors explicitly: `return nil, wc.warnings, err` or `return nil, err`
- Use `fmt.Errorf()` for formatted error messages with context
- Check errors immediately after function calls

**Testing:**
- Table-driven tests with subtests using `t.Run(tt.name, func(t *testing.T) {...})`
- Use `testdata/` directory for test fixtures
- Helper functions like `jsonToObj()` should panic on errors for test clarity

**Style:**
- Align the happy path to the left
- Ensure code adheres to following coding principles: KISS, SRP, DRY and prioritize clarity over brevity or clever 
- Use modern features of Go when possible
