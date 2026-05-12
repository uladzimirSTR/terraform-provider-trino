resource "trino_table" "example" {
  catalog  = "datalake"
  schema_name = "example_schema"
  name     = "example_table "
  location = "s3a://prod-datalake-gypsy"
  format = "parquet"

  partition_keys = [
    {
      name = "dt"
      type = "string"
    },
    {
      name = "country"
      type = "string"
    }
  ]

  columns = [
    {
      name = "id"
      type = "bigint"
    },
    {
      name = "name"
      type = "string"
    },
    {
      name = "age"
      type = "integer"
    }
  ]

}
