package provider

import (
	"context"
	"database/sql"
	"fmt"

	rndr "github.com/uladzimirSTR/terraform-provider-trino/internal/sqlrender"
)

func (r *schemaResource) createSchema(ctx context.Context, data schemaResourceModel) error {
	rnd, err := rndr.NewRenderer[rndr.CreateSchemaData]()
	if err != nil {
		return fmt.Errorf("create SQL renderer: %w", err)
	}

	sql, err := rnd.Render("create_schema.sql.tmpl", rndr.CreateSchemaData{
		IfNotExists: true,
		TableSchema: rndr.TableSchema{
			Catalog:  data.Catalog.ValueString(),
			Name:     data.Name.ValueString(),
			Location: data.Location.ValueString(),
		},
	})
	if err != nil {
		return fmt.Errorf("render create schema SQL: %w", err)
	}

	return r.client.Exec(ctx, sql)
}

func (r *schemaResource) dropSchema(ctx context.Context, data schemaResourceModel) error {
	rnd, err := rndr.NewRenderer[rndr.DropSchemaData]()
	if err != nil {
		return fmt.Errorf("create SQL renderer: %w", err)
	}

	sql, err := rnd.Render("drop_schema.sql.tmpl", rndr.DropSchemaData{
		IfExists: true,
		TableSchema: rndr.TableSchema{
			Catalog:  data.Catalog.ValueString(),
			Name:     data.Name.ValueString(),
			Location: data.Location.ValueString(),
		},
	})
	if err != nil {
		return fmt.Errorf("render drop schema SQL: %w", err)
	}

	return r.client.Exec(ctx, sql)
}

func (r *schemaResource) SchemaExists(ctx context.Context, data schemaResourceModel) (bool, error) {
	rnd, err := rndr.NewRenderer[rndr.TableSchema]()

	if err != nil {
		return false, fmt.Errorf("create SQL renderer: %w", err)
	}

	query, err := rnd.Render("schema_exists.sql.tmpl", rndr.TableSchema{
		Catalog:  data.Catalog.ValueString(),
		Name:     data.Name.ValueString(),
		Location: data.Location.ValueString(),
	})

	if err != nil {
		return false, fmt.Errorf("render schema exists SQL: %w", err)
	}

	var found string

	err = r.client.QueryRow(ctx, query).Scan(&found)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("check trino schema exists: %w", err)
	}

	return true, nil
}
