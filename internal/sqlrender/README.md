# SQL Render Package

`sqlrender` is a lightweight SQL template rendering package designed for Trino-based Terraform providers and infrastructure tooling.

The package provides:

* SQL template rendering;
* safe SQL identifier escaping;
* schema and table DDL generation;
* table property rendering;
* S3 location generation;
* reusable template helper functions;
* strongly typed rendering models.

It is designed specifically for infrastructure-as-code workflows where SQL must be generated predictably and consistently.

---

# Features

* Generic SQL renderer using Go templates
* Embedded SQL templates
* Trino-compatible DDL generation
* Safe identifier quoting
* SQL string escaping
* S3 path generation helpers
* Table property rendering
* Strongly typed render models
* Fully testable rendering layer

---

# Project Structure

```text
internal/sqlrender/
├── funcs.go                 # Template helper functions
├── models.go                # Rendering models
├── renderer.go              # Generic renderer
├── renderer_test.go         # Unit tests
└── templates/
    ├── add_columns.sql.tmpl
    ├── create_schema.sql.tmpl
    ├── create_table.sql.tmpl
    ├── drop_columns.sql.tmpl
    ├── drop_schema.sql.tmpl
    ├── drop_table.sql.tmpl
    ├── rename_column.sql.tmpl
    ├── rename_table.sql.tmpl
    ├── schema_exists.sql.tmpl
    ├── set_file_format.sql.tmpl
    ├── set_partitioning.sql.tmpl
    ├── set_table_location.sql.tmpl
    └── table_exists.sql.tmpl
```

---

# Renderer

The renderer uses Go generics and `text/template`.

```go
r, err := NewRenderer[CreateTableData]()
if err != nil {
    panic(err)
}
```

Render SQL:

```go
sql, err := r.Render("create_table.sql.tmpl", data)
if err != nil {
    panic(err)
}
```

---

# Embedded Templates

Templates are embedded directly into the binary:

```go
//go:embed templates/*.sql.tmpl
var templatesFS embed.FS
```

Advantages:

* no runtime dependency on template files;
* portable builds;
* easier Terraform provider distribution;
* deterministic rendering.

---

# Supported Operations

| Operation              | Template                      |
| ---------------------- | ----------------------------- |
| Create schema          | `create_schema.sql.tmpl`      |
| Drop schema            | `drop_schema.sql.tmpl`        |
| Create table           | `create_table.sql.tmpl`       |
| Drop table             | `drop_table.sql.tmpl`         |
| Add columns            | `add_columns.sql.tmpl`        |
| Drop columns           | `drop_columns.sql.tmpl`       |
| Rename column          | `rename_column.sql.tmpl`      |
| Rename table           | `rename_table.sql.tmpl`       |
| Set file format        | `set_file_format.sql.tmpl`    |
| Set partitioning       | `set_partitioning.sql.tmpl`   |
| Set table location     | `set_table_location.sql.tmpl` |
| Schema existence check | `schema_exists.sql.tmpl`      |
| Table existence check  | `table_exists.sql.tmpl`       |

---

# Models

## TableSchema

Represents a Trino schema.

```go
type TableSchema struct {
    Catalog  string
    Name     string
    Location string
}
```

Example:

```go
TableSchema{
    Catalog:  "datalake",
    Name:     "stage",
    Location: "s3://bucket/data",
}
```

---

## TableRef

Reference to an existing table.

```go
type TableRef struct {
    TableSchema TableSchema
    TableName   string
}
```

---

## Column

Represents a table column.

```go
type Column struct {
    ColName string
    ColType string
    Comment string
}
```

---

## TableProperties

Additional Trino table properties.

```go
type TableProperties struct {
    Format        string
    PartitionedBy []string
    Extra         map[string]any
}
```

---

## Table

Represents a full table definition.

```go
type Table struct {
    TableSchema TableSchema
    TableName   string
    Columns     []Column
    TableProp   TableProperties
}
```

---

# Example: Create Schema

```go
r, _ := NewRenderer[CreateSchemaData]()

sql, err := r.Render("create_schema.sql.tmpl", CreateSchemaData{
    IfNotExists: true,
    TableSchema: TableSchema{
        Catalog:  "datalake",
        Name:     "stage",
        Location: "s3a://test-data/",
    },
})
```

Generated SQL:

```sql
CREATE SCHEMA IF NOT EXISTS "datalake"."stage"
WITH (
  location = 's3a://test-data/stage'
)
```

---

# Example: Create Table

```go
r, _ := NewRenderer[CreateTableData]()

sql, err := r.Render("create_table.sql.tmpl", CreateTableData{
    IfNotExists: true,
    Table: Table{
        TableSchema: TableSchema{
            Catalog:  "datalake",
            Name:     "stage",
            Location: "s3://bucket/data",
        },
        TableName: "users",
        Columns: []Column{
            {
                ColName: "id",
                ColType: "BIGINT",
            },
            {
                ColName: "email",
                ColType: "VARCHAR",
                Comment: "user email",
            },
        },
        TableProp: TableProperties{
            Format: "PARQUET",
            PartitionedBy: []string{
                "created_date",
            },
        },
    },
})
```

Generated SQL:

```sql
CREATE TABLE IF NOT EXISTS "datalake"."stage"."users" (
  "id" BIGINT,
  "email" VARCHAR COMMENT 'user email'
)
WITH (
  format = 'PARQUET',
  partitioned_by = ARRAY['created_date'],
  external_location = 's3://bucket/data/stage/users'
)
```

---

# Helper Functions

## qident

Safely quotes SQL identifiers.

```go
qident(`user`)
```

Result:

```sql
"user"
```

---

## sqlString

Escapes SQL strings.

```go
sqlString("user's email")
```

Result:

```sql
'user''s email'
```

---

## fqSchema

Builds fully-qualified schema names.

```sql
"catalog"."schema"
```

---

## fqTable

Builds fully-qualified table names.

```sql
"catalog"."schema"."table"
```

---

## s3Join

Safely joins S3 paths.

Example:

```go
s3Join("s3://bucket/data", "stage", "users")
```

Result:

```text
s3://bucket/data/stage/users
```

---

## sqlValue

Converts Go values into SQL literals.

Supported types:

* `nil`
* `bool`
* `string`
* `int`
* `int64`
* `float64`
* `[]string`

Examples:

```go
sqlValue("hello")
```

```sql
'hello'
```

```go
sqlValue([]string{"a", "b"})
```

```sql
ARRAY['a', 'b']
```

---

# Table Properties

The renderer automatically generates table properties.

Example input:

```go
TableProperties{
    Format: "PARQUET",
    PartitionedBy: []string{"dt"},
}
```

Generated SQL:

```sql
WITH (
  format = 'PARQUET',
  partitioned_by = ARRAY['dt']
)
```

---

# Automatic External Location

`tableWithProps()` automatically generates:

```sql
external_location
```

Based on:

```text
schema_location/schema_name/table_name
```

Example:

```text
s3://bucket/data/stage/users
```

---

# Tests

Run tests:

```bash
go test ./...
```

Covered areas:

* schema rendering;
* table rendering;
* column rendering;
* SQL escaping;
* property rendering;
* S3 path generation.

---

# Design Goals

The package intentionally separates:

| Layer               | Responsibility               |
| ------------------- | ---------------------------- |
| `sqlrender`         | SQL generation               |
| `client`            | Trino connectivity           |
| Terraform resources | Business logic and lifecycle |

This keeps rendering deterministic and testable.

---

# Future Improvements

Potential extensions:

* column type validation;
* SQL formatter integration;
* template inheritance;
* CTAS support;
* INSERT generation;
* MERGE generation;
* Iceberg optimization helpers;
* view rendering;
* materialized view rendering;
* transactional batch rendering.

---
