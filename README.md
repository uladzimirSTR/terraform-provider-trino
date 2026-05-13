# Terraform Provider Trino

Terraform provider for managing Trino schemas and tables.

The provider generates Trino-compatible SQL and executes it through the official Go Trino client.

Currently supported resources:

- `trino_schema`
- `trino_table`

---

# Features

- Manage Trino schemas
- Manage Trino tables
- SQL rendering through embedded templates
- TLS support with custom PEM certificates
- Strongly typed Terraform resources
- Generic SQL renderer
- Unit tests for provider, renderer, and client
- Local provider installation support

---

# Project Structure

```text
.
├── main.go
├── go.mod
├── Makefile
├── examples/
│   ├── provider/
│   ├── schema/
│   └── table/
└── internal/
    ├── client/
    ├── provider/
    └── sqlrender/
````

---

# Architecture

The project is split into three layers.

## Provider Layer

Location:

```text
internal/provider
```

Responsibilities:

* Terraform provider schema
* Terraform resources
* resource lifecycle handling
* provider configuration
* Trino client initialization

---

## Client Layer

Location:

```text
internal/client
```

Responsibilities:

* Trino connection management
* TLS configuration
* DSN generation
* SQL execution
* configuration validation

Uses:

```go
database/sql
github.com/trinodb/trino-go-client/trino
```

---

## SQL Render Layer

Location:

```text
internal/sqlrender
```

Responsibilities:

* SQL template rendering
* SQL escaping
* identifier quoting
* Trino DDL generation
* S3 path generation
* table property rendering

Templates are embedded into the provider binary.

---

# Requirements

Required tools:

* Go
* Terraform
* Trino server

Recommended versions:

```text
Go:        1.25+
Terraform: 1.x
```

---

# Main Dependencies

## Terraform Plugin Framework

```text
github.com/hashicorp/terraform-plugin-framework
```

Used for:

* provider implementation
* resource lifecycle
* Terraform schemas
* diagnostics
* state management

---

## Trino Go Client

```text
github.com/trinodb/trino-go-client/trino
```

Used for:

* Trino connectivity
* SQL execution
* TLS client integration

---

## database/sql

Standard Go SQL abstraction layer.

Used for:

* query execution
* connection pooling
* context support

---

# Build Provider

Build binary:

```bash
make build
```

Equivalent:

```bash
go build -o terraform-provider-trino
```

Generated binary:

```text
terraform-provider-trino
```

---

# Run Tests

```bash
make test
```

Equivalent:

```bash
go test ./...
```

Tests cover:

* provider metadata
* provider schema
* resource schemas
* SQL rendering
* client configuration validation
* DSN generation
* TLS configuration

---

# Local Installation

Install provider locally:

```bash
make install-local
```

Provider binary will be copied into:

```text
~/.terraform.d/plugins/registry.terraform.io/<namespace>/trino/0.1.0/linux_amd64/
```

---

# Important Namespace Note

Current `main.go`:

```go
Address: "registry.terraform.io/ulstr/trino"
```

Current examples:

```hcl
source = "ulstr/trino"
```

---

# Provider Configuration

Example:

```hcl
terraform {
  required_providers {
    trino = {
      source  = "ulstr/trino"
      version = "0.1.0"
    }
  }
}

provider "trino" {
  host        = "localhost"
  port        = 8080
  user        = "admin"
  password    = "admin"
  http_scheme = "https"

  # path_to_pem   = "/path/to/certs"
  # file_name_pem = "ca.pem"
}
```

---

# Provider Arguments

| Argument        |      Required | Description       |
| --------------- | ------------: | ----------------- |
| `host`          |           yes | Trino host        |
| `port`          |           yes | Trino port        |
| `user`          |           yes | Trino username    |
| `password`      |           yes | Trino password    |
| `catalog`       |            no | Default catalog   |
| `schema_name`   |            no | Default schema    |
| `http_scheme`   | should be set | `http` or `https` |
| `path_to_pem`   |            no | PEM directory     |
| `file_name_pem` |            no | PEM filename      |

---

# TLS Support

Custom CA certificates are supported.

Example:

```hcl
provider "trino" {
  host        = "localhost"
  port        = 8443
  user        = "admin"
  password    = "admin"
  http_scheme = "https"

  path_to_pem   = "/etc/trino/certs"
  file_name_pem = "ca.pem"
}
```

Both fields must be provided together:

```hcl
path_to_pem
file_name_pem
```

---

# Resource: trino_schema

Example:

```hcl
resource "trino_schema" "example" {
  catalog  = "datalake"
  name     = "example_schema"
  location = "s3a://prod-datalake"
}
```

---

# Resource: trino_table

Example:

```hcl
resource "trino_schema" "example" {
  catalog  = "datalake"
  name     = "example_schema"
  location = "s3a://prod-datalake"
}

resource "trino_table" "example" {
  catalog     = trino_schema.example.catalog
  schema_name = trino_schema.example.name
  name        = "example_table"
  location    = trino_schema.example.location
  format      = "parquet"

  partition_keys = [
    {
      name = "dt"
      type = "varchar"
    }
  ]

  columns = [
    {
      name = "id"
      type = "bigint"
    },
    {
      name = "name"
      type = "varchar"
    }
  ]
}
```

---

# Terraform Dependency Management

The table resource references:

```hcl
trino_schema.example.catalog
trino_schema.example.name
trino_schema.example.location
```

Because of this Terraform automatically creates:

```text
schema -> table
```

without explicit `depends_on`.

---

# Generated SQL

Example schema SQL:

```sql
CREATE SCHEMA IF NOT EXISTS "datalake"."example_schema"
WITH (
  location = 's3a://prod-datalake/example_schema'
)
```

Example table SQL:

```sql
CREATE TABLE IF NOT EXISTS "datalake"."example_schema"."example_table" (
  "id" bigint,
  "name" varchar,
  "dt" varchar COMMENT 'partition key'
)
WITH (
  external_location = 's3a://prod-datalake/example_schema/example_table',
  format = 'parquet',
  partitioned_by = ARRAY['dt']
)
```

---

# Examples

## Initialize

```bash
make example-init
```

or:

```bash
cd examples/schema
terraform init
```

---

## Plan

```bash
make example-plan
```

---

## Apply

```bash
make example-apply
```

---

## Destroy

```bash
make example-destroy
```

---

# Debug Mode

Run provider in debug mode:

```bash
terraform-provider-trino -debug
```

---

# How It Works

Terraform starts the provider as a plugin process.

The provider:

1. reads Terraform configuration;
2. creates Trino client;
3. renders SQL through templates;
4. executes SQL in Trino;
5. stores Terraform state.

---

# SQL Rendering

SQL generation is separated from provider logic.

Renderer location:

```text
internal/sqlrender
```

Templates:

```text
internal/sqlrender/templates
```

Example template:

```sql
CREATE TABLE {{ if .IfNotExists }}IF NOT EXISTS {{ end }}{{ fqTableFromTable .Table }}
```

This architecture keeps:

* provider logic clean;
* SQL reusable;
* rendering testable.

Good separation of concerns. Already much better than “concatenate SQL strings and hope production survives”.

---

# Current Limitations

Current update behavior:

```text
drop old resource
create new resource
```

This is simple but destructive.

Not implemented yet:

* import support
* data sources
* safe ALTER TABLE diffing
* partial updates
* view support
* materialized views

---

# Development Workflow

Typical workflow:

```bash
go test ./...
make build
make install-local

cd examples/schema

terraform init
terraform plan
terraform apply
```

---

# Clean Build Files

```bash
make clean
```

---

# Future Improvements

Planned improvements:

* import support
* ALTER TABLE diffing
* column updates
* partition updates
* view resources
* materialized views
* data sources
* Docker test environment
* CI/CD pipeline
* Terraform Registry publishing
* acceptance tests

---
