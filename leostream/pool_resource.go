package leostream

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault" // Add this line
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

//todo
// make center_name and center_type outputs of the center resource

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &poolResource{}
	_ resource.ResourceWithConfigure   = &poolResource{}
	_ resource.ResourceWithImportState = &poolResource{}
)

// NewPoolResource is a helper function to simplify the provider implementation.
func NewPoolResource() resource.Resource {
	return &poolResource{}
}

// poolResource is the resource implementation.
type poolResource struct {
	client *leostream.Client
}

// Metadata returns the resource type name.
func (r *poolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool"
}

// Schema defines the schema for the resource.
func (r *poolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"display_name": schema.StringAttribute{
				Optional: true,
			},
			"notes": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"running_desktops_threshold": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"pool_definition": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					poolDefinitionModel{}.attrTypes(), poolDefinitionModel{}.defaultObject()),
				),
				Attributes: map[string]schema.Attribute{
					"restrict_by": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("C"),
					},
					"server_ids": schema.ListAttribute{
						ElementType: types.Int64Type,
						Optional:    true,
						Computed:    true,
						//Set Default to be an empty list
						Default: listdefault.StaticValue(types.ListNull(types.Int64Type)),
					},
					"never_rogue": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"use_vmotion": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"parent_pool_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(1),
					},
					"pool_attribute_join": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("A"),
					},
					"attributes": schema.ListNestedAttribute{
						Description: "Attributes to match",
						Optional:    true,
						Computed:    true,
						//Default:     setdefault.StaticValue(types.SetNull(types.ObjectType{AttrTypes: attributesModel{}.attrTypes()})),
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
								"vm_table_field": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"ad_attribute_field": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"vm_gpu_field": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = true

									}, "", "")},
								},
								"text_to_match": schema.StringAttribute{
									Optional: true,
									Computed: false,
									//Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"condition_type": schema.StringAttribute{
									Optional: true,
									Computed: false,
									//Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
							},
						},
					},
				},
			},
			"provision": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					provisionModel{}.attrTypes(), provisionModel{}.defaultObject()),
				),
				Attributes: map[string]schema.Attribute{
					"provision_on_off": schema.Int64Attribute{
						Optional: true,
					},
					"provision_max": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"provision_vm_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"provision_server_id": schema.Int64Attribute{
						Optional: true,
					},
					"provision_threshold": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"provision_tenant_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"provision_vm_name_next_value": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"provision_limits_enforce": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"mark_deletable": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"provision_url": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"provision_vm_display_name": schema.StringAttribute{
						Optional: true,
					},
					"provision_vm_name": schema.StringAttribute{
						Optional: true,
					},
					"center": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Default: objectdefault.StaticValue(types.ObjectValueMust(
							centerModel{}.attrTypes(), centerModel{}.defaultObject()),
						),
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Optional: true,
							},
							"name": schema.StringAttribute{
								Optional: true,
								Computed: false,
								//Default:  stringdefault.StaticString(""),
							},
							"type": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"aws_size": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"provision_method": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("image"),
							},
							"aws_iam_name": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"aws_sub_net": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"aws_sec_group": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"aws_vpc_id": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *poolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *poolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan

	var plan poolResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// empty state as it's a create operation
	var state poolResourceModel

	// defer to common function to create or update the resource

	PlStored := r.CreateNested(ctx, &plan, &state, &resp.Diagnostics)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Map response body to schema and populate Computed attribute values
	// convert int64 to string
	plan.ID = types.StringValue(strconv.FormatInt(PlStored.Stored_data.ID, 10))

	// set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *poolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve values from state
	var state poolResourceModel
	tflog.Info(ctx, "Performing state get on pool resource")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Performing Read on pool resource")

	// // use common model for state
	var newState poolResourceModel
	// use common Read function
	newState.Read(ctx, *r.client, &resp.Diagnostics, "resource", state.ID.ValueString())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// populate internal fields into new state
	newState.ID = state.ID

	//set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *poolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from plan
	var plan poolResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve values from state
	var state poolResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.UpdateNested(ctx, &plan, &state, &resp.Diagnostics)
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

func (r *poolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state poolResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing pool
	err := r.client.DeletePool(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Leostream Pool",
			"Could not delete pool, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *poolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
