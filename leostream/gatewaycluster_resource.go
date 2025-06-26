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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &gatewayClusterResource{}
	_ resource.ResourceWithConfigure = &gatewayClusterResource{}
)

// NewGatewayResource is a helper function to simplify the provider implementation.
func NewGatewayClusterResource() resource.Resource {
	return &gatewayClusterResource{}
}

// gatewayResource is the resource implementation.
type gatewayClusterResource struct {
	client *leostream.Client
}

// gatewayResourceModel maps the resource schema data.
type gatewayClusterResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Address          types.String `tfsdk:"address"`
	Load_balance_via types.String `tfsdk:"load_balance_via"`
	Use_src_ip       types.Int64  `tfsdk:"use_src_ip"`
	Notes            types.String `tfsdk:"notes"`
}

// Metadata returns the resource type name.
func (r *gatewayClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gatewaycluster"
}

// Schema defines the schema for the resource.
func (r *gatewayClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The gatewaycluster resource allows you to create, read, update, and delete gatewayclusters in Leostream.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the gatewaycluster.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Display name of the gatewaycluster.",
				Optional:    true,
			},
			"address": schema.StringAttribute{
				Description: "Public IP address of the gatewaycluster.",
				Optional:    true,
			},
			"load_balance_via": schema.StringAttribute{
				Description: `Enum: E,P,L
				E: All gateways in the cluster are used for load balancing
				P: Only the login gateway will be used
				L: Only the gateway with the least number of connections will be used`,
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"use_src_ip": schema.Int64Attribute{
				Description: `Method of source IP filtering
				0: do not use source IP filtering, but random port(default)
				1: use source IP filtering, but same port on gateway and desktop
				2: use source IP filtering, but random port on gateway
				`,
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"notes": schema.StringAttribute{
				Description: "Notes for the gatewaycluster.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *gatewayClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*leostream.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configuration Type",
			fmt.Sprintf("Expected *leostream.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create a new resource.
func (r *gatewayClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Info(ctx, "Performing update")

	// Retrieve values from plan
	var plan gatewayClusterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var gwcluster leostream.GatewayCluster
	gwcluster.Address = plan.Address.ValueString()
	gwcluster.Load_balance_via = plan.Load_balance_via.ValueString()
	gwcluster.Name = plan.Name.ValueString()
	gwcluster.Notes = plan.Notes.ValueString()
	gwcluster.Use_src_ip = plan.Use_src_ip.ValueInt64()

	// Create new gateway
	GwClusterStored, err := r.client.CreateGatewayCluster(gwcluster, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating gatewaycluster",
			"Could not create gatewaycluster, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	// convert int64 to string
	plan.ID = types.StringValue(strconv.FormatInt(GwClusterStored.Stored_data.ID, 10))

	tflog.Info(ctx, "Updating state")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *gatewayClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Info(ctx, "Performing read")

	// Get current state
	var state gatewayClusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed gatewaycluster value from Leostream
	gatewaycluster, err := r.client.GetGatewayCluster(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading  gatewaycluster configuration",
			"Could not read Leostream gatewaycluster ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to model

	state.ID = types.StringValue(strconv.FormatInt(gatewaycluster.ID, 10))
	state.Name = types.StringValue(gatewaycluster.Name)
	state.Address = types.StringValue(gatewaycluster.Address)
	state.Load_balance_via = types.StringValue(gatewaycluster.Load_balance_via)
	state.Use_src_ip = types.Int64Value(int64(gatewaycluster.Use_src_ip))
	state.Notes = types.StringValue(gatewaycluster.Notes)

	tflog.Info(ctx, "Updating state")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *gatewayClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	tflog.Info(ctx, "Performing state")

	// Retrieve values from plan
	var plan gatewayClusterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var gwcluster leostream.GatewayCluster
	gwcluster.Address = plan.Address.ValueString()
	gwcluster.Load_balance_via = plan.Load_balance_via.ValueString()
	gwcluster.Name = plan.Name.ValueString()
	gwcluster.Notes = plan.Notes.ValueString()
	gwcluster.Use_src_ip = plan.Use_src_ip.ValueInt64()

	ctx = tflog.SetField(ctx, "Plan ID", plan.ID.ValueString())
	ctx = tflog.SetField(ctx, "Address", plan.Address.ValueString())
	ctx = tflog.SetField(ctx, "Load_balance_via", plan.Load_balance_via.ValueString())
	ctx = tflog.SetField(ctx, "Notes", plan.Notes.ValueString())
	ctx = tflog.SetField(ctx, "Use_src_ip", plan.Use_src_ip.ValueInt64())
	ctx = tflog.SetField(ctx, "Name", plan.Name.ValueString())
	tflog.Info(ctx, "Updating Leostream GatewayCluster")

	// Update existing gatewaycluster
	_, err := r.client.UpdateGatewayCluster(plan.ID.ValueString(), gwcluster, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Leostream GatewayCluster",
			"Could not update gatewaycluster, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Updating state")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *gatewayClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Info(ctx, "Performing delete")

	// Retrieve values from state
	var state gatewayClusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "Plan ID", state.ID.ValueString())
	tflog.Info(ctx, "Deleting gatewaycluster")

	// Delete existing gatewaycluster
	err := r.client.DeleteGatewayCluster(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting gatewaycluster",
			"Could not delete gatewaycluster, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *gatewayClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
