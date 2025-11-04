# caddy-yaml

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

Top-level keys prefixed with `x-` are discarded, following the [Docker Compose
convention](https://docs.docker.com/reference/compose-file/extension/).
This makes it easier to leverage YAML anchors and aliases, while avoiding Caddy
errors due to unknown fields.

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

Extension fields can be reused anywhere else in
the YAML config.

```yaml
x-hello: Hello from YAML template
x-nest:
  value: nesting
```

Referencing them without `x-` prefix.

```yaml
...
handle:
  - handler: static_response
    body: "#{ .hello } with #{ .nest.value }"
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

## Example Config

Check the [test YAML configuration file](testdata/test.caddy.yaml).

## License

Apache 2
