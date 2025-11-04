# Example: Complex Caddy Configuration with Includes and Templates

This example demonstrates the new include feature combined with templates and
extension fields.

## File Structure

```
config/
├── caddy.yaml              # Main configuration
├── common/
│   ├── defaults.yaml       # Common defaults
│   └── tls.yaml           # TLS configuration
└── sites/
    ├── blog.yaml          # Blog site config
    └── api.yaml           # API site config
```

## Main Configuration (caddy.yaml)

```yaml
include:
  - path: ./common/defaults.yaml
  - path: ./common/tls.yaml
  - path:
      - ./sites/blog.yaml
      - ./sites/api.yaml

# Extension fields for templates
x-environment: production
x-log-level: INFO

apps:
  http:
    servers:
      main:
        logs:
          default_logger_name: default

#{if ne $ENVIRONMENT "production"}
logging:
  logs:
    default:
      level: DEBUG
#{else}
logging:
  logs:
    default:
      level: #{ .log_level }
#{end}
```

## Common Defaults (common/defaults.yaml)

```yaml
x-default-port: &port 443
x-timeout-duration: 30s

apps:
  http:
    servers:
      main:
        listen: [":#{ .default_port }"]
        timeouts:
          read: #{ .timeout_duration }
          write: #{ .timeout_duration }
```

## TLS Configuration (common/tls.yaml)

```yaml
x-tls-email: admin@example.com

apps:
  http:
    servers:
      main:
        automatic_https:
          email: #{ .tls_email }
```

## Blog Site (sites/blog.yaml)

```yaml
x-blog-handler: &blog
  handler: file_server
  root: /var/www/blog

apps:
  http:
    servers:
      main:
        routes:
          - match:
              - host: [blog.example.com]
            handle:
              - <<: *blog
```

## API Site (sites/api.yaml)

```yaml
x-api-handler: &api
  handler: reverse_proxy
  upstreams:
    - dial: localhost:8080

apps:
  http:
    servers:
      main:
        routes:
          - match:
              - host: [api.example.com]
            handle:
              - <<: *api
```

## Key Features Demonstrated

1. **Include Files**: Configuration split across multiple files for better organization
2. **Nested Extension Fields**: `x-` fields used within route handlers
3. **Templates with Includes**: Extension fields from included files used in main config
4. **YAML Anchors**: Anchors (`&`) and aliases (`*`) work across the merged config
5. **Environment Variables**: Mixed with extension field templates
6. **Hyphen Conversion**: `x-default-port` becomes `.default_port` in templates

## Processing Order

1. All includes are loaded and merged
2. All `x-` fields from all files become template variables
3. Templates are resolved (environment + extension variables)
4. All `x-` fields are removed from final output
5. Clean JSON configuration is generated

This results in a clean, modular configuration system that avoids repetition
while maintaining type safety and clarity.
