package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uladzimirSTR/terraform-provider-trino/internal/client"
)

var (
	_ resource.Resource              = &tableResource{}
	_ resource.ResourceWithConfigure = &tableResource{}
	// _ resource.ResourceWithImportState = &schemaResource{}
)

func NewTableResource() resource.Resource {
	return &tableResource{}
}

type tableResource struct {
	client *client.Client
}

type tableColumnModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
}

type tableResourceModel struct {
	Catalog       types.String `tfsdk:"catalog"`
	SchemaName    types.String `tfsdk:"schema_name"`
	Name          types.String `tfsdk:"name"`
	Location      types.String `tfsdk:"location"`
	ID            types.String `tfsdk:"id"`
	Format        types.String `tfsdk:"format"`
	Columns       types.List   `tfsdk:"columns"`
	PartitionKeys types.List   `tfsdk:"partition_keys"`
}

func (r *tableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (r *tableResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rschema.Schema{
		Description: "Manages a Trino table.",

		Attributes: map[string]rschema.Attribute{
			"catalog": rschema.StringAttribute{
				Description: "Trino catalog name.",
				Required:    true,
			},
			"schema_name": rschema.StringAttribute{
				Description: "Trino schema name.",
				Required:    true,
			},
			"name": rschema.StringAttribute{
				Description: "Trino table name.",
				Required:    true,
			},
			"location": rschema.StringAttribute{
				Description: "Trino table location.",
				Optional:    true,
			},
			"id": rschema.StringAttribute{
				Description: "Table identifier in format catalog.schema.table.",
				Computed:    true,
			},
			"format": rschema.StringAttribute{
				Description: "Table data format (e.g. PARQUET).",
				Optional:    true,
			},
			"columns": rschema.ListNestedAttribute{
				Description: "List of table columns.",
				Required:    true,

				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"name":        rschema.StringAttribute{Required: true},
						"type":        rschema.StringAttribute{Required: true},
						"description": rschema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}

}

func (r *tableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tableResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.createTable(ctx, plan); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Trino table",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s.%s.%s", plan.Catalog.ValueString(), plan.SchemaName.ValueString(), plan.Name.ValueString()),
	)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *tableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tableResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	exists, err := r.TableExists(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to check if Trino table exists",
			err.Error(),
		)
		return
	}

	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(
		fmt.Sprintf("%s.%s.%s", state.Catalog.ValueString(), state.SchemaName.ValueString(), state.Name.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tableResourceModel
	var state tableResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.dropTable(ctx, state); err != nil {
		resp.Diagnostics.AddError(
			"Unable to drop old Trino table during update",
			err.Error(),
		)
		return
	}

	if err := r.createTable(ctx, plan); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create new Trino table during update",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s.%s.%s", plan.Catalog.ValueString(), plan.SchemaName.ValueString(), plan.Name.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tableResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.dropTable(ctx, state); err != nil {
		resp.Diagnostics.AddError(
			"Unable to drop Trino table",
			err.Error(),
		)
		return
	}
}

func (r *tableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	trinoClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = trinoClient
}
