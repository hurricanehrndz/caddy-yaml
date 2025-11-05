# Agent Guidelines for caddy-yaml

## Build/Test Commands
- Build: `go build -v .`
- Run all tests: `go test -v -race -covermode=atomic ./...` or `make test`
- Run single test: `go test -v . -run TestApply` (or `-run TestApply/simple` for subtest)
- Lint: `make lint` (runs pre-commit, golangci-lint, and govulncheck)
- Format: `make format` or `go tool -modfile=tools.mod golangci-lint fmt`
- Get dependencies: `go get -v ./...`

## Code Style Guidelines

**Package & Imports:**
- Package: `package caddyyaml` (single word, lowercase)
- Group imports per golangci-lint gci settings: stdlib first, then default (third-party), then caddyserver prefix
- Use blank lines to separate import groups (gci, gofmt, gofumpt, goimports enabled)

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
- Follow KISS, SRP, DRY principles - prioritize clarity over brevity or cleverness
- Use modern Go 1.23 features when appropriate
- Code must pass all enabled golangci-lint checks (see .golangci.yml for full list)

**Comments:**
- Write clear package documentation (see yaml.go for example)
- Document exported functions with their behavior and return values
- Use inline comments for complex logic or non-obvious behavior
