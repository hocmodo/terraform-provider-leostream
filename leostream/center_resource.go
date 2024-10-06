// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &centerResource{}
	_ resource.ResourceWithConfigure   = &centerResource{}
	_ resource.ResourceWithImportState = &centerResource{}
)

// NewCenterResource is a helper function to simplify the provider implementation.
func NewCenterResource() resource.Resource {
	return &centerResource{}
}

// centerResource is the resource implementation.
type centerResource struct {
	client *leostream.Client
}

// Metadata returns the resource type name.
func (r *centerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_center"
}

// Schema defines the schema for the resource.
func (r *centerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The center resource allows you to create, read, update, and delete centers in Leostream.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the center.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"center_definition": schema.SingleNestedAttribute{
				Description: "Center definition",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					centerDefinitionModel{}.attrTypes(), centerDefinitionModel{}.defaultObject()),
				),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Name of the center.",
						Required:    true,
						Computed:    false,
					},
					"allow_rogue": schema.Int64Attribute{
						Description: "Assign rogue users to desktops from this center.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"allow_rogue_policy_id": schema.Int64Attribute{
						Description: "Policy for rogue users.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"continuous_autotag": schema.Int64Attribute{
						Description: "Apply auto-tags every time center is scanned.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"init_unavailable": schema.Int64Attribute{
						Description: "Initialize desktops as unavailable.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"new_as_deletable": schema.Int64Attribute{
						Description: "New desktops are deletable.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"notes": schema.StringAttribute{
						Description: "Notes for the center.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"offer_vms": schema.Int64Attribute{
						Description: "Offer VMs to users from this center.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"poll_interval": schema.Int64Attribute{
						Description: "Interval in minutes to poll the center, 0 is don't poll.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"proxy_address": schema.StringAttribute{
						Description: "Proxy address for the center.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"type": schema.StringAttribute{
						Description: "Type of the center. Currently only 'amazon' is supported.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("amazon"),
					},
					"vc_auth_method": schema.StringAttribute{
						Description: "Authorization method: For Amazon centers, either Access Key or any attached IAM role: 'access_key' or 'attached_role'.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"vc_datacenter": schema.StringAttribute{
						Description: "AWS region or a predefined value _custom if custom region is used.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"vc_name": schema.StringAttribute{
						Description: "The Access Key ID for a user with permission to access EC2.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"vc_password": schema.StringAttribute{
						Description: "The Secret Access Key for the user.",
						Optional:    true,
						Computed:    true,
						Sensitive:   true,
						Default:     stringdefault.StaticString("**********"),
					},
					"wait_inst_status": schema.Int64Attribute{
						Description: "Wait for instance status to be running before assigning desktops.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"wait_sys_status": schema.Int64Attribute{
						Description: "Wait for system status to be valid before assigning desktops.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *centerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*leostream.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *leostream.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create a new resource.
func (r *centerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan

	var plan centerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// empty state as it's a create operation
	var state centerResourceModel

	CrStored := r.CreateNested(ctx, &plan, &state, &resp.Diagnostics)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Map response body to schema and populate Computed attribute values
	// convert int64 to string
	plan.ID = types.StringValue(strconv.FormatInt(CrStored.Stored_data.ID, 10))

	// set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *centerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve values from state
	var state centerResourceModel
	tflog.Info(ctx, "Performing state get on center resource")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Performing Read on center resource")

	// // use common model for state
	var newState centerResourceModel
	// use common Read function
	newState.Read(ctx, *r.client, &resp.Diagnostics, "resource", state.ID.ValueString())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	//set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *centerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from plan
	var plan centerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve values from state
	var state centerResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_ = r.UpdateNested(ctx, &plan, &state, &resp.Diagnostics)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// update state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *centerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state centerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing center
	err := r.client.DeleteCenter(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Leostream center",
			"Could not delete center, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *centerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
