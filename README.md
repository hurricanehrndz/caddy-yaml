# caddy-yaml

[![codecov](https://codecov.io/gh/hurricanehrndz/caddy-yaml/graph/badge.svg?token=XW8QF76UVU)](https://codecov.io/gh/hurricanehrndz/caddy-yaml)
[![Tests](https://github.com/hurricanehrndz/caddy-yaml/actions/workflows/test.yml/badge.svg)](https://github.com/hurricanehrndz/caddy-yaml/actions)

YAML config adapter for Caddy with templating and YAML 1.2 support.

> **⚠️ INCOMPATIBILITY WARNING**
> This adapter is **incompatible** with all other YAML adapters (e.g.,
> [iamd3vil/caddy_yaml_adapter](https://github.com/iamd3vil/caddy_yaml_adapter)).
> Both register as the `yaml` adapter - only one can be used at a time.

## Install

Install with [xcaddy](https://github.com/caddyserver/xcaddy).

```sh
xcaddy build --with github.com/hurricanehrndz/caddy-yaml
```

## Usage

Specify with the `--adapter` flag for `caddy run`.

```sh
caddy run --config /path/to/yaml/config.yaml --adapter yaml
```

## Features

### Include Files

> **⚠️ WARNING**
> This adapter expects each individual file to be a valid YAML file on its own.

Split your configuration across multiple files using Docker Compose-style includes:

```yaml
include:
  - path: ./common/defaults.yaml
  - path: ./routes/api.yaml
  - path:
      - ./sites/site1.yaml
      - ./sites/site2.yaml
  - path: ./config.d  # Include all .yaml/.yml files from directory

apps:
  http:
    # Main configuration here
```

Includes support:
- Relative paths (resolved from the including file's directory)
- Multiple files per include entry
- Directory includes (processes all `.yaml` and `.yml` files in alphabetical order)
- Circular dependency detection
- Deep merging of configurations (conflicts will cause an error)

**Note:** When including a directory, files are processed in alphabetical order. Subdirectories are processed recursively.

### YAML 1.2 with Anchors & Aliases

Full support for YAML 1.2 anchors (`&`) and aliases (`*`):

```yaml
# anchor declaration
x-file_server: &file_server
  handler: file_server
  hide: [".git"]
  index_names: [index.html]

# reuse alias
apps:
  http:
    servers:
      srv0:
        routes:
          - handle:
              - <<: *file_server
                root: /var/www/blog/public
          # reuse alias again
          - handle:
              - <<: *file_server
                root: /var/www/api/docs
```

### Extension Fields

Top level keys prefixed with x- are discarded. This makes it easier to leverage
YAML anchors and aliases, while avoiding Caddy errors due to unknown fields.
This convention is similar and inspired by the extension feature in  [Docker
Compose](https://docs.docker.com/reference/compose-file/extension/).

```yaml
# anchor declaration
x-file_server: &file_server
  handler: file_server
  hide: [".git"]
  index_names: [index.html]

# reuse alias
...
handle:
  - <<: *file_server
    root: /var/www/blog/public

# reuse alias
...
handle:
  - <<: *file_server
    root: /var/www/api/docs
```

Extension fields can also be used as template variables (see Templating section below).

### Conditional Configurations with Templates

Use Go templates for dynamic configurations:

```yaml
#{if ne $ENVIRONMENT "production"}
logging:
  logs:
    default: { level: DEBUG }
#{end}
```

### Config-time Environment Variables

Without the Caddyfile, Caddy's native configuration limits to runtime
environment variables. There are use cases for knowing the environment
variables at configuration time (e.g., troubleshooting purposes).

```yaml
listen: "#{ $PORT }"
```


## Templating

Anything supported by [Go templates](https://pkg.go.dev/text/template) can be
used, as well as any [Sprig](https://masterminds.github.io/sprig) function.

### Delimiters

Delimiters are `#{` and `}`. e.g. `#{ .title }`. The choice of delimiters
ensures the YAML config file remains a valid YAML file that can be validated by
the schema.

### Values

Extension fields can be reused anywhere else in the YAML config as template variables.
Hyphens in extension field names are automatically converted to underscores for
template compatibility.

```yaml
x-hello: Hello from YAML template
x-api-version: v2  # Available as .api_version in templates
x-nest:
  value: nesting
```

Referencing them without `x-` prefix (hyphens become underscores).

```yaml
...
handle:
  - handler: static_response
    body: "#{ .hello } with #{ .nest.value }"
    version: "#{ .api_version }"  # Hyphen became underscore
```

_If string interpolation is not needed, YAML anchors and aliases can also be
used to achieve this_.

### Environment Variables

Environment variables can be used in a template by prefixing with `$`.

```yaml
listen:
  - "#{ $PORT }"
...
handler: file_server
root: "#{ $APP_ROOT_DIR }/public"
```

Caddy supports runtime environment variables via [`{env.*}` placeholders](https://caddyserver.com/docs/caddyfile/concepts#environment-variables).

## Processing Pipeline

The adapter processes YAML configuration in the following order:

1. **Include Processing** - Load and merge included files (with cycle detection)
2. **Extension Variable Extraction** - Extract `x-` fields for use as template variables
3. **Template Application** - Apply Go templates with environment variables and extension variables
4. **Extension Removal** - Recursively remove all `x-` prefixed keys
5. **JSON Conversion** - Convert final YAML to Caddy JSON format

This pipeline allows:
- Included files to define extension fields used in the main file
- Extension fields to use environment variables in templates
- Templates to reference both extension fields and environment variables
- Clean final output with no extension fields remaining

## Examples

- [Basic YAML configuration with templates](testdata/test.caddy.yaml)
- [Include example](testdata/test.include.yaml)
- [Nested extension fields](testdata/test.nested-extensions.yaml)
- [**Complete example with includes and templates**](EXAMPLE.md)

## Acknowledgments

The include processing and extension field handling features were inspired by
[compose-go](https://github.com/compose-spec/compose-go) (Copyright 2020 The Compose
Specification Authors), licensed under Apache License 2.0. While the implementation
is original, the design patterns and concepts follow the Docker Compose specification.

## License

Apache 2.0 - See [LICENSE](LICENSE) for details

This project includes inspiration from compose-go. See [NOTICE](NOTICE) for attribution details.
