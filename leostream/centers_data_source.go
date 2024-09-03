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
	_ datasource.DataSource              = &centersDataSource{}
	_ datasource.DataSourceWithConfigure = &centersDataSource{}
)

// NewCentersDataSource is a helper function to simplify the provider implementation.
func NewCentersDataSource() datasource.DataSource {
	return &centersDataSource{}
}

// centersDataSource is the data source implementation.
type centersDataSource struct {
	client *leostream.Client
}

// Metadata returns the data source type name.
func (d *centersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_centers"
}

// Schema defines the schema for the data source.
func (d *centersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"centers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"os": schema.StringAttribute{
							Computed: true,
						},
						"flavor": schema.StringAttribute{
							Computed: true,
						},
						"online": schema.Int64Attribute{
							Computed: true,
						},
						"status": schema.Int64Attribute{
							Computed: true,
						},
						"status_label": schema.StringAttribute{
							Computed: true,
						},
						"center_type": schema.StringAttribute{
							Computed: true,
						},
						"type_label": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *centersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state centersDataSourceModel

	centers, err := d.client.GetCenters()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Leostream Centers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, center := range centers {
		centerState := centersModel{
			ID:           types.Int64Value(int64(center.ID)),
			Name:         types.StringValue(center.Name),
			Os:           types.StringValue(center.Os),
			Flavor:       types.StringValue(center.Flavor),
			Online:       types.Int64Value(int64(center.Online)),
			Status:       types.Int64Value(int64(center.Status)),
			Status_label: types.StringValue(center.Status_label),
			Center_type:  types.StringValue(center.Center_type),
			Type_label:   types.StringValue(center.Type_label),
		}

		state.Centers = append(state.Centers, centerState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *centersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*leostream.Client)
}

// centersDataSourceModel maps the data source schema data.
type centersDataSourceModel struct {
	Centers []centersModel `tfsdk:"centers"`
}

// centersModel maps centers schema data.
type centersModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Os           types.String `tfsdk:"os"`
	Flavor       types.String `tfsdk:"flavor"`
	Online       types.Int64  `tfsdk:"online"`
	Status       types.Int64  `tfsdk:"status"`
	Status_label types.String `tfsdk:"status_label"`
	Center_type  types.String `tfsdk:"center_type"`
	Type_label   types.String `tfsdk:"type_label"`
}
