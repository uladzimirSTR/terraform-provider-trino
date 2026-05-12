package provider

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	rndr "github.com/uladzimirSTR/terraform-provider-trino/internal/sqlrender"
)

func convertToTableColumns(ctx context.Context, cols types.List, prt types.List) []rndr.Column {
	var tfColumns []tableColumnModel
	var tfPartitions []tablePartitionModel

	ColsErr := cols.ElementsAs(ctx, &tfColumns, false)
	PartsErr := prt.ElementsAs(ctx, &tfPartitions, false)

	columns := make([]rndr.Column, 0, len(tfColumns)+len(tfPartitions))

	if ColsErr != nil || PartsErr != nil {
		return columns
	}

	for _, col := range tfColumns {
		columns = append(columns, rndr.Column{
			ColName: col.Name.ValueString(),
			ColType: col.Type.ValueString(),
			Comment: col.Description.ValueString(),
		})
	}

	for _, part := range tfPartitions {
		columns = append(columns, rndr.Column{
			ColName: part.Name.ValueString(),
			ColType: part.Type.ValueString(),
			Comment: "partition key",
		})
	}

	return columns
}

func convertToPartitionKeys(ctx context.Context, list types.List) []string {
	if list.IsNull() {
		return []string{}
	}

	var tfKeys []tablePartitionModel

	err := list.ElementsAs(ctx, &tfKeys, false)

	if err != nil {
		return []string{}
	}

	keys := make([]string, 0, len(tfKeys))

	for _, k := range tfKeys {
		keys = append(keys, k.Name.ValueString())
	}

	return keys
}

func (r *tableResource) createTable(ctx context.Context, data tableResourceModel) error {
	rnd, err := rndr.NewRenderer[rndr.CreateTableData]()

	if err != nil {
		return fmt.Errorf("create SQL renderer: %w", err)
	}

	sql, err := rnd.Render("create_table.sql.tmpl", rndr.CreateTableData{
		IfNotExists: true,
		Table: rndr.Table{
			TableSchema: rndr.TableSchema{
				Catalog:  data.Catalog.ValueString(),
				Name:     data.SchemaName.ValueString(),
				Location: data.Location.ValueString(),
			},
			TableName: data.Name.ValueString(),
			Columns:   convertToTableColumns(ctx, data.Columns, data.PartitionKeys),
			TableProp: rndr.TableProperties{
				Format:        data.Format.ValueString(),
				PartitionedBy: convertToPartitionKeys(ctx, data.PartitionKeys),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("render create table SQL: %w", err)
	}

	return r.client.Exec(ctx, sql)
}

func (r *tableResource) dropTable(ctx context.Context, data tableResourceModel) error {
	rnd, err := rndr.NewRenderer[rndr.DropTableData]()

	if err != nil {
		return fmt.Errorf("create SQL renderer: %w", err)
	}

	sql, err := rnd.Render("drop_table.sql.tmpl", rndr.DropTableData{
		IfExists: true,
		Table: rndr.TableRef{
			TableSchema: rndr.TableSchema{
				Catalog:  data.Catalog.ValueString(),
				Name:     data.SchemaName.ValueString(),
				Location: data.Location.ValueString(),
			},
			TableName: data.Name.ValueString(),
		},
	})

	if err != nil {
		return fmt.Errorf("render drop table SQL: %w", err)
	}

	return r.client.Exec(ctx, sql)
}

func (r *tableResource) TableExists(ctx context.Context, data tableResourceModel) (bool, error) {
	rnd, err := rndr.NewRenderer[rndr.TableRef]()

	if err != nil {
		return false, fmt.Errorf("create SQL renderer: %w", err)
	}

	query, err := rnd.Render("table_exists.sql.tmpl", rndr.TableRef{
		TableSchema: rndr.TableSchema{
			Catalog:  data.Catalog.ValueString(),
			Name:     data.SchemaName.ValueString(),
			Location: data.Location.ValueString(),
		},
		TableName: data.Name.ValueString(),
	})

	if err != nil {
		return false, fmt.Errorf("render table exists SQL: %w", err)
	}

	var found string

	err = r.client.QueryRow(ctx, query).Scan(&found)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("check trino table exists: %w", err)
	}

	return true, nil

}
