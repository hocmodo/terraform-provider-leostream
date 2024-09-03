package leostream

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// centerResourceModel maps the resource schema data.
type centerResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Center_definition types.Object `tfsdk:"center_definition"`
}

// attrTypes - return attribute types for this model
// func (o centerResourceModel) attrTypes() map[string]attr.Type {
// 	return map[string]attr.Type{
// 		"id":                types.StringType,
// 		"center_definition": types.ObjectType{AttrTypes: centerDefinitionModel{}.attrTypes()},
// 	}
// }

// defaultObject - return default object for this model
// func (o centerResourceModel) defaultObject() map[string]attr.Value {
// 	return map[string]attr.Value{
// 		"id":                types.StringValue(""),
// 		"center_definition": types.ObjectValueMust(centerDefinitionModel{}.attrTypes(), centerDefinitionModel{}.defaultObject()),
// 	}
// }

// nested attributes objects

// centerDefinitionModel maps filtering schema data
type centerDefinitionModel struct {
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
func (o centerDefinitionModel) attrTypes() map[string]attr.Type {
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

// defaultObject - return default object for this model
func (o centerDefinitionModel) defaultObject() map[string]attr.Value {
	return map[string]attr.Value{
		"name":                  types.StringValue(""),
		"allow_rogue":           types.Int64Value(0),
		"allow_rogue_policy_id": types.Int64Value(0),
		"continuous_autotag":    types.Int64Value(0),
		"init_unavailable":      types.Int64Value(0),
		"new_as_deletable":      types.Int64Value(0),
		"notes":                 types.StringValue(""),
		"offer_vms":             types.Int64Value(1),
		"poll_interval":         types.Int64Value(1),
		"proxy_address":         types.StringValue(""),
		"type":                  types.StringValue("amazon"),
		"vc_auth_method":        types.StringValue(""),
		"vc_datacenter":         types.StringValue(""),
		"vc_name":               types.StringValue(""),
		"vc_password":           types.StringValue("**********"),
		"wait_inst_status":      types.Int64Value(0),
		"wait_sys_status":       types.Int64Value(0),
	}
}

// common `Read` function for both data source and resource
func (o *centerResourceModel) Read(ctx context.Context, client leostream.Client, diags *diag.Diagnostics, rtype string, id string) {
	//center CONFIG
	//get refreshed center config value from Leostream API
	centerConfig, err := client.GetCenter(id)

	if err != nil {
		diags.AddError(
			"Unable to read center Configuration",
			err.Error(),
		)
		return
	}

	o.ID = types.StringValue(strconv.FormatInt(centerConfig.ID, 10))

	// Map center definition to state
	var statecenterDefinition centerDefinitionModel
	statecenterDefinition.Name = types.StringValue(centerConfig.Center_definition.Name)
	statecenterDefinition.Allow_rogue = types.Int64Value(centerConfig.Center_definition.Allow_rogue)
	statecenterDefinition.Allow_rogue_policy_id = types.Int64Value(centerConfig.Center_definition.Allow_rogue_policy_id)
	statecenterDefinition.Continuous_autotag = types.Int64Value(centerConfig.Center_definition.Continuous_autotag)
	statecenterDefinition.Init_unavailable = types.Int64Value(centerConfig.Center_definition.Init_unavailable)
	statecenterDefinition.New_as_deletable = types.Int64Value(centerConfig.Center_definition.New_as_deletable)
	statecenterDefinition.Notes = types.StringValue(centerConfig.Center_definition.Notes)
	statecenterDefinition.Offer_vms = types.Int64Value(centerConfig.Center_definition.Offer_vms)
	statecenterDefinition.Poll_interval = types.Int64Value(centerConfig.Center_definition.Poll_interval)
	statecenterDefinition.Proxy_address = types.StringValue(centerConfig.Center_definition.Proxy_address)
	statecenterDefinition.Type = types.StringValue(centerConfig.Center_definition.Type)
	statecenterDefinition.Vc_auth_method = types.StringValue(centerConfig.Center_definition.Vc_auth_method)
	statecenterDefinition.Vc_datacenter = types.StringValue(centerConfig.Center_definition.Vc_datacenter)
	statecenterDefinition.Vc_name = types.StringValue(centerConfig.Center_definition.Vc_name)
	statecenterDefinition.Vc_password = types.StringValue(centerConfig.Center_definition.Vc_password)
	statecenterDefinition.Wait_inst_status = types.Int64Value(centerConfig.Center_definition.Wait_inst_status)
	statecenterDefinition.Wait_sys_status = types.Int64Value(centerConfig.Center_definition.Wait_sys_status)

	//Add center definition to center model
	o.Center_definition, _ = types.ObjectValueFrom(ctx, centerDefinitionModel{}.attrTypes(), &statecenterDefinition)

}

// `Create` function for the resource
func (r *centerResource) CreateNested(ctx context.Context, plan *centerResourceModel, state *centerResourceModel, diags *diag.Diagnostics) *leostream.CenterStored {
	// center CONFIG

	// Instantiate empty object for storing plan data
	var centerConfig leostream.Center

	// Unpack nested attributes from plan for the center definition
	var plancenterDefinition centerDefinitionModel
	*diags = plan.Center_definition.As(ctx, &plancenterDefinition, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}

	// Instantiate empty object for storing plan data for the center definition object in the center config
	var centerDefinitionConfig leostream.CenterDefinition

	// Populate center_definition field restrict_by in empty object from plan
	centerDefinitionConfig.Name = plancenterDefinition.Name.ValueString()
	centerDefinitionConfig.Allow_rogue = plancenterDefinition.Allow_rogue.ValueInt64()
	centerDefinitionConfig.Allow_rogue_policy_id = plancenterDefinition.Allow_rogue_policy_id.ValueInt64()
	centerDefinitionConfig.Continuous_autotag = plancenterDefinition.Continuous_autotag.ValueInt64()
	centerDefinitionConfig.Init_unavailable = plancenterDefinition.Init_unavailable.ValueInt64()
	centerDefinitionConfig.New_as_deletable = plancenterDefinition.New_as_deletable.ValueInt64()
	centerDefinitionConfig.Notes = plancenterDefinition.Notes.ValueString()
	centerDefinitionConfig.Offer_vms = plancenterDefinition.Offer_vms.ValueInt64()
	centerDefinitionConfig.Poll_interval = plancenterDefinition.Poll_interval.ValueInt64()
	centerDefinitionConfig.Proxy_address = plancenterDefinition.Proxy_address.ValueString()
	centerDefinitionConfig.Type = plancenterDefinition.Type.ValueString()
	centerDefinitionConfig.Vc_auth_method = plancenterDefinition.Vc_auth_method.ValueString()
	centerDefinitionConfig.Vc_datacenter = plancenterDefinition.Vc_datacenter.ValueString()
	centerDefinitionConfig.Vc_name = plancenterDefinition.Vc_name.ValueString()
	centerDefinitionConfig.Vc_password = plancenterDefinition.Vc_password.ValueString()
	centerDefinitionConfig.Wait_inst_status = plancenterDefinition.Wait_inst_status.ValueInt64()
	centerDefinitionConfig.Wait_sys_status = plancenterDefinition.Wait_sys_status.ValueInt64()

	// Assign the center definition config to the center config
	centerConfig.Center_definition = centerDefinitionConfig

	// Create new center
	centersStored, err := r.client.CreateCenter(centerConfig, nil)

	if err != nil {
		diags.AddError(
			"Unable to Create center",
			err.Error(),
		)
		return nil
	} else {
		return centersStored
	}
}

// `Update` function for the resource
func (r *centerResource) UpdateNested(ctx context.Context, plan *centerResourceModel, state *centerResourceModel, diags *diag.Diagnostics) *leostream.CenterStored {
	// center CONFIG

	// Instantiate empty object for storing plan data
	var centerConfig leostream.Center

	// Unpack nested attributes from plan for the center definition
	var plancenterDefinition centerDefinitionModel
	*diags = plan.Center_definition.As(ctx, &plancenterDefinition, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}
	// Instantiate empty object for storing plan data for the center definition object in the center config
	var centerDefinitionConfig leostream.CenterDefinition

	// Populate center definition config from plan
	centerDefinitionConfig.Name = plancenterDefinition.Name.ValueString()
	centerDefinitionConfig.Allow_rogue = plancenterDefinition.Allow_rogue.ValueInt64()
	centerDefinitionConfig.Allow_rogue_policy_id = plancenterDefinition.Allow_rogue_policy_id.ValueInt64()
	centerDefinitionConfig.Continuous_autotag = plancenterDefinition.Continuous_autotag.ValueInt64()
	centerDefinitionConfig.Init_unavailable = plancenterDefinition.Init_unavailable.ValueInt64()
	centerDefinitionConfig.New_as_deletable = plancenterDefinition.New_as_deletable.ValueInt64()
	centerDefinitionConfig.Notes = plancenterDefinition.Notes.ValueString()
	centerDefinitionConfig.Offer_vms = plancenterDefinition.Offer_vms.ValueInt64()
	centerDefinitionConfig.Poll_interval = plancenterDefinition.Poll_interval.ValueInt64()
	centerDefinitionConfig.Proxy_address = plancenterDefinition.Proxy_address.ValueString()
	centerDefinitionConfig.Type = plancenterDefinition.Type.ValueString()
	centerDefinitionConfig.Vc_auth_method = plancenterDefinition.Vc_auth_method.ValueString()
	centerDefinitionConfig.Vc_datacenter = plancenterDefinition.Vc_datacenter.ValueString()
	centerDefinitionConfig.Vc_name = plancenterDefinition.Vc_name.ValueString()
	//centerDefinitionConfig.Vc_password = plancenterDefinition.Vc_password.ValueString()
	centerDefinitionConfig.Wait_inst_status = plancenterDefinition.Wait_inst_status.ValueInt64()
	centerDefinitionConfig.Wait_sys_status = plancenterDefinition.Wait_sys_status.ValueInt64()

	// Assign the center definition config to the center config
	centerConfig.Center_definition = centerDefinitionConfig

	// Update center
	centersStored, err := r.client.UpdateCenter(plan.ID.ValueString(), centerConfig, nil)

	if err != nil {
		diags.AddError(
			"Unable to modify center",
			err.Error(),
		)
		return nil
	} else {
		return centersStored
	}
}
