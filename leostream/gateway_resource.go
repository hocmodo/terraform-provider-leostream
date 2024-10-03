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
	_ resource.Resource              = &gatewayResource{}
	_ resource.ResourceWithConfigure = &gatewayResource{}
)

// NewGatewayResource is a helper function to simplify the provider implementation.
func NewGatewayResource() resource.Resource {
	return &gatewayResource{}
}

// gatewayResource is the resource implementation.
type gatewayResource struct {
	client *leostream.Client
}

// gatewayResourceModel maps the resource schema data.
type gatewayResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Address          types.String `tfsdk:"address"`
	Address_private  types.String `tfsdk:"address_private"`
	Load_balancer_id types.Int64  `tfsdk:"load_balancer_id"`
	Use_src_ip       types.Int64  `tfsdk:"use_src_ip"`
	Notes            types.String `tfsdk:"notes"`
}

// Metadata returns the resource type name.
func (r *gatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

// Schema defines the schema for the resource.
func (r *gatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the gateway.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Display name of the gateway.",
				Optional: true,
			},
			"address": schema.StringAttribute{
				Description: "Public IP address of the gateway.",
				Optional: true,
			},
			"address_private": schema.StringAttribute{
				Description: "Private IP address of the gateway.",
				Optional: true,
			},
			"load_balancer_id": schema.Int64Attribute{
				Description: "ID of the cluster associated with the gateway.",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
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
				Description: "Notes for the gateway.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *gatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *gatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gatewayResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var gw leostream.Gateway
	gw.Address_private = plan.Address_private.ValueString()
	gw.Address = plan.Address.ValueString()
	gw.Load_balancer_id = plan.Load_balancer_id.ValueInt64()
	gw.Name = plan.Name.ValueString()
	gw.Notes = plan.Notes.ValueString()
	gw.Use_src_ip = plan.Use_src_ip.ValueInt64()

	// Create new gateway
	GwStored, err := r.client.CreateGateway(gw, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating gateway",
			"Could not create gateway, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	// convert int64 to string
	plan.ID = types.StringValue(strconv.FormatInt(GwStored.Stored_data.ID, 10))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *gatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state gatewayResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed gateway value from Leostream
	gateway, err := r.client.GetGateway(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Leostream Gateway",
			"Could not read Leostream gateway ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to model

	state.ID = types.StringValue(strconv.FormatInt(gateway.ID, 10))
	state.Name = types.StringValue(gateway.Name)
	state.Address = types.StringValue(gateway.Address)
	state.Address_private = types.StringValue(gateway.Address_private)
	state.Load_balancer_id = types.Int64Value(int64(gateway.Load_balancer_id))
	state.Use_src_ip = types.Int64Value(int64(gateway.Use_src_ip))
	state.Notes = types.StringValue(gateway.Notes)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *gatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan gatewayResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var gw leostream.Gateway
	gw.Address_private = plan.Address_private.ValueString()
	gw.Address = plan.Address.ValueString()
	gw.Load_balancer_id = plan.Load_balancer_id.ValueInt64()
	gw.Name = plan.Name.ValueString()
	gw.Notes = plan.Notes.ValueString()
	gw.Use_src_ip = plan.Use_src_ip.ValueInt64()

	ctx = tflog.SetField(ctx, "Plan ID", plan.ID.ValueString())
	ctx = tflog.SetField(ctx, "Address", plan.Address.ValueString())
	ctx = tflog.SetField(ctx, "Address_private", plan.Address_private.ValueString())
	ctx = tflog.SetField(ctx, "Load_balancer_id", plan.Load_balancer_id.ValueInt64())
	ctx = tflog.SetField(ctx, "Notes", plan.Notes.ValueString())
	ctx = tflog.SetField(ctx, "Use_src_ip", plan.Use_src_ip.ValueInt64())
	ctx = tflog.SetField(ctx, "Name", plan.Name.ValueString())
	tflog.Info(ctx, "Updating Leostream Gateway")

	// Update existing gateway
	_, err := r.client.UpdateGateway(plan.ID.ValueString(), gw, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Leostream Gateway",
			"Could not update gateway, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *gatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state gatewayResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "Plan ID", state.ID.ValueString())
	tflog.Info(ctx, "Deleting Leostream Gateway")

	// Delete existing gateway
	err := r.client.DeleteGateway(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Leostream Gateway",
			"Could not delete gateway, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *gatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
