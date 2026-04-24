package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// trinoProvider is the provider implementation.
type trinoProvider struct {
	version string
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &trinoProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &trinoProvider{
			version: version,
		}
	}
}

// trinoProviderModel maps provider schema data to a Go type.
type trinoProviderModel struct {
	Host        types.String `tfsdk:"host"`
	Port        types.Int32  `tfsdk:"port"`
	User        types.String `tfsdk:"user"`
	Password    types.String `tfsdk:"password"`
	Catalog     types.String `tfsdk:"catalog"`
	SchemaName  types.String `tfsdk:"schema_name"`
	HttpScheme  types.String `tfsdk:"http_scheme"`
	Auth        types.String `tfsdk:"auth"`
	Verify      types.String `tfsdk:"verify"`
	PathToPem   types.String `tfsdk:"path_to_pem"`
	FileNamePem types.String `tfsdk:"file_name_pem"`
}

// Metadata returns the provider type name.
func (p *trinoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "trino"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *trinoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{

		Description: "Trino provider for managing Trino resources.",

		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The hostname or IP address of the Trino server.",
				Required:    true,
			},
			"port": schema.Int32Attribute{
				Description: "The port number on which the Trino server is listening.",
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "The username for authenticating with the Trino server.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for authenticating with the Trino server.",
				Required:    true,
				Sensitive:   true,
			},
			"catalog": schema.StringAttribute{
				Description: "The default catalog to use when connecting to the Trino server.",
				Optional:    true,
			},
			"schema_name": schema.StringAttribute{
				Description: "The default schema to use when connecting to the Trino server.",
				Optional:    true,
			},
			"http_scheme": schema.StringAttribute{
				Description: "The HTTP scheme to use when connecting to the Trino server (e.g., http or https).",
				Optional:    true,
			},
			"auth": schema.StringAttribute{
				Description: "The authentication method to use when connecting to the Trino server (e.g., basic, kerberos).",
				Optional:    true,
			},
			"verify": schema.StringAttribute{
				Description: "Whether to verify the server's TLS certificate when connecting to the Trino server (e.g., true or false).",
				Optional:    true,
			},
			"path_to_pem": schema.StringAttribute{
				Description: "The file path to the PEM file containing the TLS certificate for the Trino server.",
				Optional:    true,
			},
			"file_name_pem": schema.StringAttribute{
				Description: "The file name of the PEM file containing the TLS certificate for the Trino server.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a Trino API client for data sources and resources.
func (p *trinoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config trinoProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() || config.Host.IsNull() {

		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Trino host",
			"The provider cannot create a Trino client without a host.",
		)
	}

	if config.User.IsUnknown() || config.User.IsNull() {

		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Missing Trino user",
			"The provider cannot create a Trino client without a user.",
		)
	}

	if config.Password.IsUnknown() || config.Password.IsNull() {

		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Trino password",
			"The provider cannot create a Trino client without a password.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// c := &client.Client{

	// 	Host:     config.Host.ValueString(),
	// 	Port:     int(config.Port.ValueInt32()),
	// 	User:     config.User.ValueString(),
	// 	Password: config.Password.ValueString(),
	// 	Catalog:  config.Catalog.ValueString(),
	// 	Schema:   config.SchemaName.ValueString(),
	// 	// HTTPS:    useHTTPS,
	// 	// Insecure: insecure,
	// }
	// c.BaseURL = fmt.Sprintf("%s://%s:%d", c.Protocol(), c.Host, c.Port)
	// resp.DataSourceData = c
	// resp.ResourceData = c

}

// DataSources defines the data sources implemented in the provider.
func (p *trinoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *trinoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
