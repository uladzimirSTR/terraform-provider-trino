package sqlrender

import (
	"strings"
	"testing"
)

func TestRenderCreateSchema(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatal(err)
	}

	sql, err := r.Render("create_schema.sql.tmpl", CreateSchemaData{
		IfNotExists: true,
		TableSchema: TableSchema{
			Catalog:  "datalake",
			Name:     "stage",
			Location: "s3a://test-data/",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []string{
		`CREATE SCHEMA IF NOT EXISTS "datalake"."stage"`,
		`WITH (`,
		`location = 's3a://test-data/stage'`,
		`)`,
	}

	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
}

func TestRenderCreateTable(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatal(err)
	}

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
				{ColName: "id", ColType: "BIGINT"},
				{ColName: "email", ColType: "VARCHAR", Comment: "user's email"},
			},
			TableProp: TableProperties{
				Format:        "PARQUET",
				PartitionedBy: []string{"created_date"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []string{
		`CREATE TABLE IF NOT EXISTS "datalake"."stage"."users"`,
		`"id" BIGINT,`,
		`"email" VARCHAR COMMENT 'user''s email'`,
		`format = 'PARQUET'`,
		`partitioned_by = ARRAY['created_date']`,
		`external_location = 's3://bucket/data/stage/users'`,
	}

	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
}

func TestRenderAddColumns(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatal(err)
	}

	sql, err := r.Render("add_columns.sql.tmpl", AddColumnsData{
		Table: TableRef{
			TableSchema: TableSchema{
				Catalog: "datalake",
				Name:    "stage",
			},
			TableName: "users",
		},
		Columns: []Column{
			{ColName: "age", ColType: "INTEGER"},
			{ColName: "country", ColType: "VARCHAR"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []string{
		`ALTER TABLE "datalake"."stage"."users"`,
		`"age" INTEGER,`,
		`"country" VARCHAR`,
	}

	for _, check := range checks {
		if !strings.Contains(sql, check) {
			t.Fatalf("expected SQL to contain %q\nSQL:\n%s", check, sql)
		}
	}
}
