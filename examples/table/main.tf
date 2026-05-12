resource "trino_schema" "example" {
  catalog  = "datalake"
  name     = "example_schema"
  location = "s3a://prod-datalake-gypsy"
}

resource "trino_table" "example" {
  catalog  = trino_schema.example.catalog
  schema_name = trino_schema.example.name
  name     = "example_table"
  location = trino_schema.example.location
  format = "parquet"

  partition_keys = [
    {
      name = "dt"
      type = "varchar"
    },
    {
      name = "country"
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
