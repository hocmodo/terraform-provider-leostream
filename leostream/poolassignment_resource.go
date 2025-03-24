// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &poolAssignmentResource{}
	_ resource.ResourceWithConfigure   = &poolAssignmentResource{}
	_ resource.ResourceWithImportState = &poolAssignmentResource{}
)

// NewpoolAssignmentResource is a helper function to simplify the provider implementation.
func NewpoolAssignmentResource() resource.Resource {
	return &poolAssignmentResource{}
}

// poolAssignmentResource is the resource implementation.
type poolAssignmentResource struct {
	client *leostream.Client
}

// Metadata returns the resource type name.
func (r *poolAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_assignment"
}

// Schema defines the schema for the resource.
func (r *poolAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The poolassignment resource allows you to create, read, update, and delete poolassignments in Leostream.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the poolassignment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_id": schema.Int64Attribute{
				Description: "Pool ID for which this assignment is valid",
				Required:    true,
			},
			"policy_id": schema.Int64Attribute{
				Description: "Policy ID for which this assignment is valid",
				Required:    true,
			},
			"offer_filter": schema.StringAttribute{
				Description: `The method used to decide whether desktops from this pool will be included in the offer.
				0: Included for all users (default)
				1: Only included if user's AD record matches the offer_filter_json criteria
				2: Only included if current date and time is within the offer_filter_json time ranges`,
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("0"),
			},
			"offer_filter_json": schema.SingleNestedAttribute{
				Description: "offer_filter_json",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					offerFilterJsonModel{}.attrTypes(), offerFilterJsonModel{}.defaultObject()),
				),
				Attributes: map[string]schema.Attribute{
					"join": schema.StringAttribute{
						Description: "And/Or join condition - A or O",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("O"),
					},
					"filters": schema.ListNestedAttribute{
						Description: "Array container for Pool attributes (restrict_by is 'A') or for LDAP attributes (restrict_by is 'Z', requires Active Directory Centers).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
								// If the plan has a value for the nested object, we need to replace

								resp.RequiresReplace = false

							}, "", ""),
						},
						NestedObject: schema.NestedAttributeObject{
							//Add a PlanModifier to the NestedObject
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
									// If the plan has a value for the nested object, we need to replace

									resp.RequiresReplace = false

								}, "", ""),
							},
							Attributes: map[string]schema.Attribute{
								"offer_filter_attribute": schema.StringAttribute{
									Description: "Offer_filter_attribute",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"offer_filter_condition": schema.StringAttribute{
									Description: "Offer_filter_condition",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"offer_filter_value": schema.StringAttribute{
									Description: "Offer_filter_value",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = true

									}, "", "")},
								},
							},
						},
					},
				},
			},
			"plan_protocol_id": schema.Int64Attribute{
				Description: "ID of protocol plan to assign",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"plan_power_control_id": schema.Int64Attribute{
				Description: "ID of power plan to assign",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"plan_release_id": schema.Int64Attribute{
				Description: "ID of release plan to assign",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"offer_quantity": schema.Int64Attribute{
				Description: "The number of VMs to offer to a user at login",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"display_mode": schema.StringAttribute{
				Description: `How to describe the offered machines:
				0 = Desktop name (default)
				5 = Desktop display name
				1 = Machine name
				2 = Pool name
				3 = Pool name : Desktop name
				6 = Pool name : Desktop display name
				4 = Pool name : Machine name
				7 = Pool display name
				8 = Pool display name : Desktop name
				9 = Pool display name : Desktop display name
				10= Pool display name : Machine name`,
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("0"),
			},
			"start_if_stopped": schema.Int64Attribute{
				Description: "A boolean field indicating whether to attempt to power on a machine if it's currently stopped/suspended.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *poolAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *poolAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan

	var plan poolAssignmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// empty state as it's a create operation
	var state poolAssignmentResourceModel

	CrStored := r.CreateNested(ctx, &plan, &state, &resp.Diagnostics, plan.Policy_id.String())
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
func (r *poolAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve values from state
	var state poolAssignmentResourceModel
	tflog.Info(ctx, "Performing state get on center resource")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Performing Read on center resource")

	// // use common model for state
	var newState poolAssignmentResourceModel
	// use common Read function
	newState.Read(ctx, *r.client, &resp.Diagnostics, "resource", state.Policy_id.String(), state.ID.ValueString())
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

func (r *poolAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from plan
	var plan poolAssignmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve values from state
	var state poolAssignmentResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_ = r.UpdateNested(ctx, &plan, &state, &resp.Diagnostics, state.Policy_id.String())
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

func (r *poolAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state poolAssignmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing center
	err := r.client.DeletePoolAssignment(state.ID.ValueString(), state.Policy_id.String(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Leostream poolassignment",
			"Could not delete poolassignment, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *poolAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: <policy_id>:<pool_assignment_id>
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import identifier with format: <policy_id>:<pool_assignment_id>. Got: %q", req.ID),
		)
		return
	}

	policyID, err := strconv.ParseInt(idParts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Policy ID",
			fmt.Sprintf("Cannot parse policy_id as integer: %v", err),
		)
		return
	}

	poolAssignmentID := idParts[1]

	// Set both required attributes in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), policyID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), poolAssignmentID)...)
}
