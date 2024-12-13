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
	_ datasource.DataSource              = &policiesDataSource{}
	_ datasource.DataSourceWithConfigure = &policiesDataSource{}
)

// NewPoliciesDataSource is a helper function to simplify the provider implementation.
func NewPoliciesDataSource() datasource.DataSource {
	return &policiesDataSource{}
}

// policiesDataSource is the data source implementation.
type policiesDataSource struct {
	client *leostream.Client
}

// Metadata returns the data source type name.
func (d *policiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policies"
}

// Schema defines the schema for the data source.
func (d *policiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The policies data source allows you to retrieve a list of policies from Leostream.`,
		Attributes: map[string]schema.Attribute{
			"policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Unique identifier for the policy.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Display name of the policy.",
							Computed:    true,
						},
						"notes": schema.StringAttribute{
							Description: "Notes for this policy.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *policiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state policiesDataSourceModel

	policies, err := d.client.GetPolicies()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Leostream Policies",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, policy := range policies {
		policyState := policiesModel{
			ID:    types.Int64Value(int64(policy.ID)),
			Name:  types.StringValue(policy.Name),
			Notes: types.StringValue(policy.Notes),
		}

		state.Policies = append(state.Policies, policyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *policiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*leostream.Client)
}

// policiesDataSourceModel maps the data source schema data.
type policiesDataSourceModel struct {
	Policies []policiesModel `tfsdk:"policies"`
}

// policiesModel maps policies schema data.
type policiesModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Notes types.String `tfsdk:"notes"`
}
