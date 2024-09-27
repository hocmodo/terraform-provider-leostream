// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &gatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewaysDataSource{}
)

// NewGatewaysDataSource is a helper function to simplify the provider implementation.
func NewGatewaysDataSource() datasource.DataSource {
	return &gatewaysDataSource{}
}

// gatewaysDataSource is the data source implementation.
type gatewaysDataSource struct {
	client *leostream.Client
}

// Metadata returns the data source type name.
func (d *gatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateways"
}

// Schema defines the schema for the data source.
func (d *gatewaysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"gateways": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Unique identifier for the gateway.",
							Computed: true,
						},
						"name": schema.StringAttribute{
							Description: "Display name of the gateway.",
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *gatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state gatewaysDataSourceModel

	gateways, err := d.client.GetGateways()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Leostream Gateways",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, gateway := range gateways {
		gatewayState := gatewaysModel{
			ID:   types.Int64Value(int64(gateway.ID)),
			Name: types.StringValue(gateway.Name),
		}

		state.Gateways = append(state.Gateways, gatewayState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *gatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*leostream.Client)
}

// gatewaysDataSourceModel maps the data source schema data.
type gatewaysDataSourceModel struct {
	Gateways []gatewaysModel `tfsdk:"gateways"`
}

// gatewaysModel maps gateways schema data.
type gatewaysModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
