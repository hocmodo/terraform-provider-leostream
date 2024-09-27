// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
func NewCenterDataSource() datasource.DataSource {
	return &centerDataSource{}
}

// centersDataSource is the data source implementation.
type centerDataSource struct {
	client *leostream.Client
}

// Metadata returns the data source type name.
func (d *centerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_center_ds"
}

// Schema defines the schema for the data source.
func (d *centerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"images": schema.ListNestedAttribute{
				Description: "List of available AMI's",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"name": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"center_definition": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Name of the center.",
						Optional: true,
						Computed: true,
					},
					"allow_rogue": schema.Int64Attribute{
						Description: "Assign rogue users to desktops from this center.",
						Optional: true,
						Computed: true,
					},
					"allow_rogue_policy_id": schema.Int64Attribute{
						Description: "Policy for rogue users.",
						Optional: true,
						Computed: true,
					},
					"continuous_autotag": schema.Int64Attribute{
						Description: "Apply auto-tags every time center is scanned.",
						Optional: true,
						Computed: true,
					},
					"init_unavailable": schema.Int64Attribute{
						Description: "Initialize desktops as unavailable.",
						Optional: true,
						Computed: true,
					},
					"new_as_deletable": schema.Int64Attribute{
						Description: "New desktops are deletable.",
						Optional: true,
						Computed: true,
					},
					"notes": schema.StringAttribute{
						Description: "Notes for the center.",
						Optional: true,
						Computed: true,
					},
					"offer_vms": schema.Int64Attribute{
						Description: "Offer VMs to users from this center.",
						Optional: true,
						Computed: true,
					},
					"poll_interval": schema.Int64Attribute{
						Description: "Interval in minutes to poll the center, 0 is don't poll.",
						Optional: true,
						Computed: true,
					},
					"proxy_address": schema.StringAttribute{
						Description: "Proxy address for the center.",
						Optional: true,
						Computed: true,
					},
					"type": schema.StringAttribute{
						Description: "Type of the center. Currently only 'amazon' is supported.",
						Optional: true,
						Computed: true,
					},
					"vc_auth_method": schema.StringAttribute{
						Description: "Authorization method: For Amazon centers, either Access Key or any attached IAM role: 'access_key' or 'attached_role'.",
						Optional: true,
						Computed: true,
					},
					"vc_datacenter": schema.StringAttribute{
						Description: "AWS region or a predefined value _custom if custom region is used.",
						Optional: true,
						Computed: true,
					},
					"vc_name": schema.StringAttribute{
						Description: "The Access Key ID for a user with permission to access EC2.",
						Optional: true,
						Computed: true,
					},
					"vc_password": schema.StringAttribute{
						Description: "The Secret Access Key for the user.",
						Optional: true,
						Computed: true,
					},
					"wait_inst_status": schema.Int64Attribute{
						Description: "Wait for instance status to be running before assigning desktops.",
						Optional: true,
						Computed: true,
					},
					"wait_sys_status": schema.Int64Attribute{
						Description: "Wait for system status to be valid before assigning desktops.",
						Optional: true,
						Computed: true,
					},
				},
			},
			"center_info": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"os": schema.StringAttribute{
						Description: "Operating System",
						Optional: true,
						Computed: true,
					},
					"os_version": schema.StringAttribute{
						Description: "Operating System Version",
						Optional: true,
						Computed: true,
					},
					"aws_sizes": schema.ListAttribute{
						Description: "List of available AWS sizes",
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
					},
					"aws_sub_nets": schema.ListAttribute{
						Description: "List of available AWS subnets",
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
					},
					// "aws_sec_groups": schema.SetNestedAttribute{
					// 	Optional: true,
					// 	Computed: true,
					// 	NestedObject: schema.NestedAttributeObject{
					// 		Attributes: map[string]schema.Attribute{
					// 			"gid": schema.StringAttribute{
					// 				Optional: true,
					// 				Computed: true,
					// 			},
					// 			"gname": schema.StringAttribute{
					// 				Optional: true,
					// 				Computed: true,
					// 			},
					// 			"gdesc": schema.StringAttribute{
					// 				Optional: true,
					// 				Computed: true,
					// 			},
					// 			"vpcid": schema.StringAttribute{
					// 				Optional: true,
					// 				Computed: true,
					// 			},
					// 		},
					// 	},
					// },
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *centerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state centerDSModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	center, err := d.client.GetCenter(state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Leostream Center",
			err.Error(),
		)
		return
	}

	// Convert center.ID from int64 to string
	state.ID = types.Int64Value(center.ID)

	// Create a new center_definition state model
	var stateCenterDefinitionDataSourceModel centerDefinitionDataSourceModel

	stateCenterDefinitionDataSourceModel.Name = types.StringValue(center.Center_definition.Name)
	stateCenterDefinitionDataSourceModel.Allow_rogue = types.Int64Value(center.Center_definition.Allow_rogue)
	stateCenterDefinitionDataSourceModel.Allow_rogue_policy_id = types.Int64Value(center.Center_definition.Allow_rogue_policy_id)
	stateCenterDefinitionDataSourceModel.Continuous_autotag = types.Int64Value(center.Center_definition.Continuous_autotag)
	stateCenterDefinitionDataSourceModel.Init_unavailable = types.Int64Value(center.Center_definition.Init_unavailable)
	stateCenterDefinitionDataSourceModel.New_as_deletable = types.Int64Value(center.Center_definition.New_as_deletable)
	stateCenterDefinitionDataSourceModel.Notes = types.StringValue(center.Center_definition.Notes)
	stateCenterDefinitionDataSourceModel.Offer_vms = types.Int64Value(center.Center_definition.Offer_vms)
	stateCenterDefinitionDataSourceModel.Poll_interval = types.Int64Value(center.Center_definition.Poll_interval)
	stateCenterDefinitionDataSourceModel.Proxy_address = types.StringValue(center.Center_definition.Proxy_address)
	stateCenterDefinitionDataSourceModel.Type = types.StringValue(center.Center_definition.Type)
	stateCenterDefinitionDataSourceModel.Vc_auth_method = types.StringValue(center.Center_definition.Vc_auth_method)
	stateCenterDefinitionDataSourceModel.Vc_datacenter = types.StringValue(center.Center_definition.Vc_datacenter)
	stateCenterDefinitionDataSourceModel.Vc_name = types.StringValue(center.Center_definition.Vc_name)
	stateCenterDefinitionDataSourceModel.Vc_password = types.StringValue(center.Center_definition.Vc_password)
	stateCenterDefinitionDataSourceModel.Wait_inst_status = types.Int64Value(center.Center_definition.Wait_inst_status)
	stateCenterDefinitionDataSourceModel.Wait_sys_status = types.Int64Value(center.Center_definition.Wait_sys_status)
	//stateCenterDefinitionDataSourceModel.aws_sizes, *diags = types.ListValueFrom(ctx, types.StringType, center.Center_definition.Aws_sizes)

	// Map response body to model
	state.Center_definition, _ = types.ObjectValueFrom(ctx, centerDefinitionDataSourceModel{}.attrTypes(), &stateCenterDefinitionDataSourceModel)

	// Create a new center_info state model
	var stateCenterInfoDataSourceModel centerInfoDataSourceModel

	//stateCenterInfoDataSourceModel.aws_sec_groups, *diags = types.SetValueFrom(ctx, awsSecGroupsModel{}.attrTypes(), center.Center_info.Aws_sec_groups)
	stateCenterInfoDataSourceModel.Aws_sizes, _ = types.ListValueFrom(ctx, types.StringType, center.Center_info.Aws_sizes)
	stateCenterInfoDataSourceModel.Aws_sub_nets, _ = types.ListValueFrom(ctx, types.StringType, center.Center_info.Aws_sub_nets)
	stateCenterInfoDataSourceModel.Os = types.StringValue(center.Center_info.Os)
	stateCenterInfoDataSourceModel.Os_version = types.StringValue(center.Center_info.Os_version)

	// Map response body to model
	state.Center_info, _ = types.ObjectValueFrom(ctx, centerInfoDataSourceModel{}.attrTypes(), &stateCenterInfoDataSourceModel)

	var stateImages []imagesModel
	// Loop through images and convert to state model
	for _, image := range center.Images {
		var stateImage imagesModel
		//Convert image.ID from int64 to string
		stateImage.ID = types.Int64Value(image.ID)
		stateImage.Name = types.StringValue(image.Name)
		stateImages = append(stateImages, stateImage)
	}

	state.Images, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: imagesModel{}.attrTypes()}, stateImages)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *centerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*leostream.Client)
}

// centersDataSourceModel maps the data source schema data.
type centerDefinitionDataSourceModel struct {
	Name                  types.String `tfsdk:"name"`
	Allow_rogue           types.Int64  `tfsdk:"allow_rogue"`
	Allow_rogue_policy_id types.Int64  `tfsdk:"allow_rogue_policy_id"`
	Continuous_autotag    types.Int64  `tfsdk:"continuous_autotag"`
	Init_unavailable      types.Int64  `tfsdk:"init_unavailable"`
	New_as_deletable      types.Int64  `tfsdk:"new_as_deletable"`
	Notes                 types.String `tfsdk:"notes"`
	Offer_vms             types.Int64  `tfsdk:"offer_vms"`
	Poll_interval         types.Int64  `tfsdk:"poll_interval"`
	Proxy_address         types.String `tfsdk:"proxy_address"`
	Type                  types.String `tfsdk:"type"`
	Vc_auth_method        types.String `tfsdk:"vc_auth_method"`
	Vc_datacenter         types.String `tfsdk:"vc_datacenter"`
	Vc_name               types.String `tfsdk:"vc_name"`
	Vc_password           types.String `tfsdk:"vc_password"`
	Wait_inst_status      types.Int64  `tfsdk:"wait_inst_status"`
	Wait_sys_status       types.Int64  `tfsdk:"wait_sys_status"`
}

// attrTypes - return attribute types for this model
func (o centerDefinitionDataSourceModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"allow_rogue":           types.Int64Type,
		"allow_rogue_policy_id": types.Int64Type,
		"continuous_autotag":    types.Int64Type,
		"init_unavailable":      types.Int64Type,
		"new_as_deletable":      types.Int64Type,
		"notes":                 types.StringType,
		"offer_vms":             types.Int64Type,
		"poll_interval":         types.Int64Type,
		"proxy_address":         types.StringType,
		"type":                  types.StringType,
		"vc_auth_method":        types.StringType,
		"vc_datacenter":         types.StringType,
		"vc_name":               types.StringType,
		"vc_password":           types.StringType,
		"wait_inst_status":      types.Int64Type,
		"wait_sys_status":       types.Int64Type,
	}
}

// centersDataSourceModel maps the data source schema data.
type centerInfoDataSourceModel struct {
	//	Aws_sec_groups types.Set    `tfsdk:"aws_sec_groups"`
	Aws_sizes    types.List   `tfsdk:"aws_sizes"`
	Aws_sub_nets types.List   `tfsdk:"aws_sub_nets"`
	Os           types.String `tfsdk:"os"`
	Os_version   types.String `tfsdk:"os_version"`
}

// attrTypes - return attribute types for this model
func (o centerInfoDataSourceModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		//		"aws_sec_groups": types.SetType{ElemType: types.ObjectType{AttrTypes: awsSecGroupsModel{}.attrTypes()}},
		"aws_sizes":    types.ListType{ElemType: types.StringType},
		"aws_sub_nets": types.ListType{ElemType: types.StringType},
		"os":           types.StringType,
		"os_version":   types.StringType,
	}
}

// attributesModel maps pool definition attribute schema data
type imagesModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (o imagesModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.Int64Type,
		"name": types.StringType,
	}
}

// centersModel maps centers schema data.
type centerDSModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Center_definition types.Object `tfsdk:"center_definition"`
	Center_info       types.Object `tfsdk:"center_info"`
	Images            types.List   `tfsdk:"images"`
}

// attrTypes - return attribute types for this model
// func (o centerDSModel) attrTypes() map[string]attr.Type {
// 	return map[string]attr.Type{
// 		"id":                types.Int64Type,
// 		"center_definition": types.ObjectType{AttrTypes: centerDefinitionDataSourceModel{}.attrTypes()},
// 		"center_info":       types.ObjectType{AttrTypes: centerInfoDataSourceModel{}.attrTypes()},
// 		"images":            types.ListType{ElemType: types.ObjectType{AttrTypes: imagesModel{}.attrTypes()}},
// 	}
// }
