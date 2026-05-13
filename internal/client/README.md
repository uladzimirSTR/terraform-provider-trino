# Terraform Provider Trino Client

Lightweight Go client for executing Trino queries inside a Terraform provider.

The client is designed to:

- validate provider configuration;
- build Trino DSN strings;
- support HTTP and HTTPS connections;
- support custom TLS CA certificates via PEM files;
- execute SQL queries through `database/sql`;
- provide safe resource cleanup;
- integrate with `trinodb/trino-go-client`.

---

# Features

- Simple Trino client abstraction
- TLS support with custom CA
- Config validation
- DSN builder
- Query execution helpers
- Context support
- Unit tested

---

# Project Structure

```text
internal/client/
├── client.go         # Main client implementation
├── config.go         # Config structure, TLS client, DSN builder
├── validate.go       # Config validation
└── client_test.go    # Unit tests
````

---

# Installation

```bash
go get github.com/trinodb/trino-go-client/trino
```

---

# Configuration

```go
cfg := client.Config{
    Host:       "trino.example.com",
    Port:       8443,
    User:       "terraform",
    Password:   "secret",
    HTTPScheme: "https",

    Catalog:    "iceberg",
    Schema:     "analytics",

    PathToPEM:   "/etc/certs",
    FileNamePEM: "ca.pem",

    QueryTimeout: 60 * time.Second,
}
```

---

# Config Fields

| Field          | Description                      |
| -------------- | -------------------------------- |
| `Host`         | Trino host                       |
| `Port`         | Trino port                       |
| `User`         | Trino username                   |
| `Password`     | Optional password authentication |
| `HTTPScheme`   | `http` or `https`                |
| `Catalog`      | Default Trino catalog            |
| `Schema`       | Default Trino schema             |
| `PathToPEM`    | Directory containing CA PEM file |
| `FileNamePEM`  | CA PEM filename                  |
| `QueryTimeout` | Query timeout duration           |

---

# TLS Support

The client supports custom CA certificates.

Both fields must be provided together:

```go
PathToPEM
FileNamePEM
```

Example:

```go
PathToPEM: "/etc/trino/certs",
FileNamePEM: "ca.pem",
```

The client:

1. Reads the PEM file
2. Creates a custom CA pool
3. Builds a custom TLS-enabled HTTP client
4. Registers it in the Trino driver

---

# Creating Client

```go
ctx := context.Background()

cfg := client.Config{
    Host:       "localhost",
    Port:       8080,
    User:       "admin",
    HTTPScheme: "http",
}

c, err := client.NewClient(cfg)
if err != nil {
    panic(err)
}

defer c.Close()
```

---

# Executing Queries

```go
err := c.Exec(ctx, `
CREATE SCHEMA IF NOT EXISTS analytics
`)
if err != nil {
    panic(err)
}
```

---

# QueryRow Example

```go
row := c.QueryRow(ctx, `
SELECT current_user
`)

var currentUser string

if err := row.Scan(&currentUser); err != nil {
    panic(err)
}
```

---

# DSN Example

Generated DSN format:

```text
https://user:password@trino.example.com:8443?catalog=iceberg&schema=analytics&source=terraform-provider-trino
```

If custom TLS client is used:

```text
https://user:password@trino.example.com:8443?catalog=iceberg&schema=analytics&custom_client=trino-provider-custom-tls&source=terraform-provider-trino
```

---

# Validation Rules

The client validates:

* host is required;
* port is required;
* user is required;
* scheme must be `http` or `https`;
* password authentication requires HTTPS;
* PEM path and filename must be provided together.

---

# Error Examples

## Invalid Scheme

```text
http scheme must be http or https
```

## Password With HTTP

```text
password authentication requires https
```

## Invalid PEM

```text
failed to append certificates from pem file
```

---

# Tests

Run tests:

```bash
go test ./...
```

Covered areas:

* config validation
* DSN generation
* TLS client creation
* PEM validation
* query validation
* client creation
* connection cleanup

---

# Example Usage Inside Terraform Provider

```go
cfg := client.Config{
    Host:       data.Host.ValueString(),
    Port:       int(data.Port.ValueInt64()),
    User:       data.User.ValueString(),
    Password:   data.Password.ValueString(),
    HTTPScheme: data.HTTPscheme.ValueString(),

    Catalog: data.Catalog.ValueString(),
    Schema:  data.Schema.ValueString(),
}

trinoClient, err := client.NewClient(cfg)
if err != nil {
    resp.Diagnostics.AddError(
        "Unable to create Trino client",
        err.Error(),
    )

    return
}

defer trinoClient.Close()
```

---

# Design Notes

This client intentionally stays minimal.

Responsibilities:

* transport setup;
* authentication;
* TLS handling;
* query execution;
* integration with Terraform provider resources.

Business logic, SQL rendering, and Terraform resource behavior should remain outside the client layer.

---

# Future Improvements

Possible future extensions:

* connection ping support;
* retry logic;
* query cancellation;
* structured logging;
* metrics;
* prepared statement helpers;
* transaction support;
* query result iterators;
* configurable HTTP transport.

---
