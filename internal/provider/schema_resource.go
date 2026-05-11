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
	_ resource.Resource              = &schemaResource{}
	_ resource.ResourceWithConfigure = &schemaResource{}
	// _ resource.ResourceWithImportState = &schemaResource{}
)

func NewSchemaResource() resource.Resource {
	return &schemaResource{}
}

type schemaResource struct {
	client *client.Client
}

type schemaResourceModel struct {
	Catalog  types.String `tfsdk:"catalog"`
	Name     types.String `tfsdk:"name"`
	Location types.String `tfsdk:"location"`
	ID       types.String `tfsdk:"id"`
}

func (r *schemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *schemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rschema.Schema{
		Description: "Manages a Trino schema.",

		Attributes: map[string]rschema.Attribute{
			"catalog": rschema.StringAttribute{
				Description: "Trino catalog name.",
				Required:    true,
			},
			"name": rschema.StringAttribute{
				Description: "Trino schema name.",
				Required:    true,
			},
			"location": rschema.StringAttribute{
				Description: "Trino schema location.",
				Optional:    true,
			},
			"id": rschema.StringAttribute{
				Description: "Schema identifier in format catalog.name.",
				Computed:    true,
			},
		},
	}
}

func (r *schemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemaResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.createSchema(ctx, plan); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Trino schema",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s.%s", plan.Catalog.ValueString(), plan.Name.ValueString()),
	)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemaResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	exists, err := r.readSchema(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Trino schema",
			err.Error(),
		)
		return
	}

	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(
		fmt.Sprintf("%s.%s", state.Catalog.ValueString(), state.Name.ValueString()),
	)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemaResourceModel
	var state schemaResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.dropSchema(ctx, state); err != nil {
		resp.Diagnostics.AddError(
			"Unable to drop old Trino schema during update",
			err.Error(),
		)
		return
	}

	if err := r.createSchema(ctx, plan); err != nil {
		resp.Diagnostics.AddError(
			"Unable to create new Trino schema during update",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s.%s", plan.Catalog.ValueString(), plan.Name.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemaResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.dropSchema(ctx, state); err != nil {
		resp.Diagnostics.AddError(
			"Unable to drop Trino schema",
			err.Error(),
		)
		return
	}
}

func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
