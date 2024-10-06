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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
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
	_ resource.Resource                = &awsPoolResource{}
	_ resource.ResourceWithConfigure   = &awsPoolResource{}
	_ resource.ResourceWithImportState = &awsPoolResource{}
)

// NewPoolResource is a helper function to simplify the provider implementation.
func NewAwsPoolResource() resource.Resource {
	return &awsPoolResource{}
}

// poolResource is the resource implementation.
type awsPoolResource struct {
	client *leostream.Client
}

// Metadata returns the resource type name.
func (r *awsPoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_pool"
}

// Schema defines the schema for the resource.
func (r *awsPoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The AWS pool resource allows you to manage Leostream pools. AWS Pools are used to group desktops in AWS together for management and provisioning.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the pool.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the pool.",
				Optional:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Display name of the pool.",
				Optional:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes for the pool.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"running_desktops_threshold": schema.Int64Attribute{
				Description: "Running and available desktops in the pool.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"pool_definition": schema.SingleNestedAttribute{
				Description: "Pool definition",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					awsPoolDefinitionModel{}.attrTypes(), awsPoolDefinitionModel{}.defaultObject()),
				),
				Attributes: map[string]schema.Attribute{
					"restrict_by": schema.StringAttribute{
						Description: `Restrict by:
						A = by attribute (default)
						T = by tag
						C = by centers
						E = vSphere hosts
						L = vSphere clusters
						V = vSphere resource pools
						Z = LDAP attributes
						H = ad hoc list (selection from parent pool)
						`,
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("C"),
					},
					"server_ids": schema.ListAttribute{
						Description: "List of tag IDs defining this pool",
						ElementType: types.Int64Type,
						Optional:    true,
						Computed:    true,
						//Set Default to be an empty list
						Default: listdefault.StaticValue(types.ListNull(types.Int64Type)),
					},
					"never_rogue": schema.Int64Attribute{
						Description: "0 or 1: A boolean field indicating if desktops in this pool treat any user as the assigned user",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"use_vmotion": schema.Int64Attribute{
						Description: "0 or 1: A boolean field indicating whether VMs of this pool will vMotion to new host",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"parent_pool_id": schema.Int64Attribute{
						Description: "ID of the parent pool",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(1),
					},
					"pool_attribute_join": schema.StringAttribute{
						Description: `A or O: How do the pool attributes get joined:
						A = And
						O = Or
						`,
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("A"),
					},
					"attributes": schema.ListNestedAttribute{
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
								"vm_table_field": schema.StringAttribute{
									Description: `The machine's attribute to search; must be a column in the vm table. Cannot exist if ad_attribute_field or vm_gpu_field is populated.
									name - Name;
									display_name - Display name;
									windows_name - Machine name;
									ip - Hostname or IP address;
									partition_names - Disk partition name;
									partition_mount_points - Partition mount point;
									guest_os - Operating system;
									os_version - Operating system version;
									installed_protocols - Installed protocols;
									vc_memory_mb - Memory (in MB);
									vc_num_cpu - Number of CPUs;
									vc_num_ethernet_cards - Number of NICs;
									num_disks - Number of disks;
									computer_model - Computer model;
									bios_serial_number - BIOS serial number;
									max_clock_speed - CPU speed (GHz);
									notes - Notes;
									vc_annotation - Center "Notes";
									tag_filter - Tags;
									server_id - Servers.
									`,
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"ad_attribute_field": schema.StringAttribute{
									Description: `Desktop attribute, mandatory for LDAP attributes,
									see possible values for an AD Center in centers.get response, field ldap_attributes.
									annot exist if vm_table_field or vm_gpu_field is populated.`,
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"vm_gpu_field": schema.StringAttribute{
									Description: "The GPU field to search; must be a column in the vm_gpu table. Cannot exist if vm_table_field or ad_attribute_field is populated.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = true

									}, "", "")},
								},
								"text_to_match": schema.StringAttribute{
									Description: "The free form text attribute",
									Optional:    true,
									Computed:    false,
									//Default:  stringdefault.StaticString(""),
									PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
										// If the plan has a value for the nested object, we need to replace

										resp.RequiresReplace = false

									}, "", "")},
								},
								"condition_type": schema.StringAttribute{
									Description: `The search conditional:
									ip - "matches (CIDR notation)";
									np - "does not match (CIDR)";
									eq - "is equal to";
									ne - "is not equal to";
									gt - "is greater than";
									lt - "is less than";
									ct - "contains";
									nc - "does not contain";
									bw - "begins with";
									ew - "ends with".`,
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
				Description: `Container for parameters related to Provisioning.
				Provisioning parameters depends on what Centers are defined in the Connection Broker and which sets of values in every Center type (e.g. Azure, AWS, etc.) are defined.
				`,
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					awsProvisionModel{}.attrTypes(), awsProvisionModel{}.defaultObject()),
				),
				Attributes: map[string]schema.Attribute{
					"provision_on_off": schema.Int64Attribute{
						Description: `
						A boolean field indicating if state of provisioning for this pool is:
						Running - provision according to thresholds
						Stopped - disabled by user or the Broker by error
						`,
						Optional: true,
					},
					"provision_max": schema.Int64Attribute{
						Description: "The maximum number of new machines that will be provisioned when the threshold is reached.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"provision_vm_id": schema.Int64Attribute{
						Description: "The ID of the server which will do the provisioning, or 0 if URL notification only",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"provision_server_id": schema.Int64Attribute{
						Description: "The ID of the server which will do the provisioning, or 0 if URL notification only",
						Optional:    true,
					},
					"provision_threshold": schema.Int64Attribute{
						Description: "Minimum number of available VMs before triggering provisioning.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"provision_tenant_id": schema.Int64Attribute{
						Description: "The tenant to provision into",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"provision_limits_enforce": schema.Int64Attribute{
						Description: "0 or 1: A boolean field indicating if Broker creates and deletes virtual machines to meet the start and max threshold.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"mark_deletable": schema.Int64Attribute{
						Description: "0 or 1: Specifies whether to initialize newly-provisioned desktops as 'deletable'.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
					"provision_url": schema.StringAttribute{
						Description: "The URL to notify when a new machine is provisioned.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"provision_vm_display_name": schema.StringAttribute{
						Description: "The display name of the VM to be provisioned.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"provision_vm_name": schema.StringAttribute{
						Description: "The name of the VM to be provisioned.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"center": schema.SingleNestedAttribute{
						Description: `Container for parameters related to Center. Offically:
						Google (object)
						or RHEV (object)
						or Scale (object)
						or OpenStack (object)
						or (Amazon AWS (Provision from image (object)
						or Provision from launch template (object)))
						or Azure (object)
						or vCenter (object)
						or ProvisionCenter (null) (ProvisionCenter).
						!This versoim of the provider only supports AWS.`,
						Optional: false,
						Required: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Description: "Unique identifier for the center.",
								Optional:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the center.",
								Optional:    true,
								Computed:    false,
							},
							"type": schema.StringAttribute{
								Description: "Type of the center. Currently only AWS is supported: amazon",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
							"aws_size": schema.StringAttribute{
								Description: `The size of the instance to provision.
								eg. t2.micro`,
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"provision_method": schema.StringAttribute{
								Description: "The method of provisioning. Currently only 'image' is supported.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("image"),
							},
							"aws_iam_name": schema.StringAttribute{
								Description: "The name of the IAM role to use for the instance.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
							"aws_sub_net": schema.StringAttribute{
								Description: "The subnet ID to use for the instance.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
							"aws_sec_group": schema.StringAttribute{
								Description: "The security group name to use for the instance.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
							"aws_vpc_id": schema.StringAttribute{
								Description: "The VPC ID to use for the instance.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString(""),
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *awsPoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *awsPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan

	var plan awsPoolResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// empty state as it's a create operation
	var state awsPoolResourceModel

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
func (r *awsPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve values from state
	var state awsPoolResourceModel
	tflog.Info(ctx, "Performing state get on pool resource")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Performing Read on pool resource")

	// // use common model for state
	var newState awsPoolResourceModel
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

func (r *awsPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from plan
	var plan awsPoolResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve values from state
	var state awsPoolResourceModel
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

func (r *awsPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state awsPoolResourceModel
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

func (r *awsPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
