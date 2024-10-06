// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
	"os"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &leostreamProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &leostreamProvider{}
}

// leostreamProvider is the provider implementation.
type leostreamProvider struct{}

// Metadata returns the provider type name.
func (p *leostreamProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "leostream"
}

// Schema defines the provider-level schema for configuration data.
func (p *leostreamProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The Leostream provider allows you to interact with the Leostream REST API to manage Leostream resources. The username created must have API access.`,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Leostream REST API.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Leostream REST API.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Leostream REST API.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *leostreamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring Leostream client")

	// Retrieve provider data from configuration
	var config leostreamProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown leostream API Host",
			"The provider cannot create the leostream API client as there is an unknown configuration value for the leostream API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the leostream_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown leostream API Username",
			"The provider cannot create the leostream API client as there is an unknown configuration value for the leostream API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the leostream_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown leostream API Password",
			"The provider cannot create the leostream API client as there is an unknown configuration value for the leostream API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the leostream_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("leostream_HOST")
	username := os.Getenv("leostream_USERNAME")
	password := os.Getenv("leostream_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Leostream API Host",
			"The provider cannot create the Leostream API client as there is a missing or empty value for the leostream API host. "+
				"Set the host value in the configuration or use the Leostream environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Leostream API Username",
			"The provider cannot create the leostream API client as there is a missing or empty value for the Leostream API username. "+
				"Set the username value in the configuration or use the Leostream environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Leostream API Password",
			"The provider cannot create the Leostream API client as there is a missing or empty value for the Leostream API password. "+
				"Set the password value in the configuration or use the Leostream environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "leostream_host", host)
	ctx = tflog.SetField(ctx, "leostream_username", username)
	ctx = tflog.SetField(ctx, "leostream_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "leostream_password")

	tflog.Debug(ctx, "Creating Leostream client")

	// Create a new Leostream client using the configuration values
	client, err := leostream.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Leostream API Client",
			"An unexpected error occurred when creating the Leostream API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Leostream Client Error: "+err.Error(),
		)
		return
	}

	// Make the Leostream client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Leostream client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *leostreamProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCentersDataSource,
		NewGatewaysDataSource,
		NewCenterDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *leostreamProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGatewayResource,
		NewAwsPoolResource,
		NewBasicPoolResource,
		NewCenterResource,
	}
}

// leostreamProviderModel maps provider schema data to a Go type.
type leostreamProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}
