resource "trino_schema" "example" {
  catalog  = "datalake"
  name     = "example_schema"
  location = "s3a://prod-datalake-gypsy"
}
