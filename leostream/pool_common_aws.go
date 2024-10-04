// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"strconv"
	//"reflect"
	//"regexp"
	//"strings"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// poolResourceModel maps the resource schema data.
type awsPoolResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Display_name               types.String `tfsdk:"display_name"`
	Notes                      types.String `tfsdk:"notes"`
	Running_desktops_threshold types.Int64  `tfsdk:"running_desktops_threshold"`
	Pool_definition            types.Object `tfsdk:"pool_definition"`
	Provision                  types.Object `tfsdk:"provision"`
}

// // attrTypes - return attribute types for this model
// func (o poolResourceModel) attrTypes() map[string]attr.Type {
// 	return map[string]attr.Type{
// 		"id":                         types.StringType,
// 		"name":                       types.StringType,
// 		"display_name":               types.StringType,
// 		"notes":                      types.StringType,
// 		"running_desktops_threshold": types.Int64Type,
// 		"pool_definition":            types.ObjectType{AttrTypes: poolDefinitionModel{}.attrTypes()},
// 	}
// }

// defaultObject - return default object for this model
// func (o poolResourceModel) defaultObject() map[string]attr.Value {
// 	return map[string]attr.Value{
// 		"id":                         types.StringValue(""),
// 		"name":                       types.StringValue(""),
// 		"display_name":               types.StringValue(""),
// 		"notes":                      types.StringValue(""),
// 		"running_desktops_threshold": types.Int64Value(0),
// 		"pool_definition":            types.ObjectValueMust(poolDefinitionModel{}.attrTypes(), poolDefinitionModel{}.defaultObject()),
// 	}
// }

// nested attributes objects

// poolDefinitionModel maps filtering schema data
type awsPoolDefinitionModel struct {
	Restrict_by         types.String `tfsdk:"restrict_by"`
	Pool_attribute_join types.String `tfsdk:"pool_attribute_join"`
	Server_ids          types.List   `tfsdk:"server_ids"`
	Never_rogue         types.Int64  `tfsdk:"never_rogue"`
	Use_vmotion         types.Int64  `tfsdk:"use_vmotion"`
	Parent_pool_id      types.Int64  `tfsdk:"parent_pool_id"`
	Attributes          types.List   `tfsdk:"attributes"`
}

// attrTypes - return attribute types for this model
func (o awsPoolDefinitionModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"restrict_by":         types.StringType,
		"pool_attribute_join": types.StringType,
		"server_ids":          types.ListType{ElemType: types.Int64Type},
		"never_rogue":         types.Int64Type,
		"use_vmotion":         types.Int64Type,
		"parent_pool_id":      types.Int64Type,
		"attributes":          types.ListType{ElemType: types.ObjectType{AttrTypes: awsAttributesModel{}.attrTypes()}},
	}
}

// defaultObject - return default object for this model
func (o awsPoolDefinitionModel) defaultObject() map[string]attr.Value {
	bootstrap_serverids := convertToAttrInt64(CONFIG_POOL_SERVERIDS)

	return map[string]attr.Value{
		"restrict_by":         types.StringValue(CONFIG_POOL_RESTRICT_BY),
		"pool_attribute_join": types.StringValue(CONFIG_POOL_ATTRIBUTE_JOIN),
		// Let the default value be an empty list
		"server_ids":     types.ListValueMust(types.Int64Type, bootstrap_serverids),
		"never_rogue":    types.Int64Value(0),
		"use_vmotion":    types.Int64Value(0),
		"parent_pool_id": types.Int64Value(0),
		"attributes":     types.ListNull(types.ObjectType{AttrTypes: awsAttributesModel{}.attrTypes()}),
	}
}

// poolProvisionModel maps filtering schema data
type awsProvisionModel struct {
	Provision_on_off    types.Int64  `tfsdk:"provision_on_off"`
	Provision_max       types.Int64  `tfsdk:"provision_max"`
	Provision_vm_id     types.Int64  `tfsdk:"provision_vm_id"`
	Provision_server_id types.Int64  `tfsdk:"provision_server_id"`
	Provision_vm_name   types.String `tfsdk:"provision_vm_name"`
	Provision_threshold types.Int64  `tfsdk:"provision_threshold"`
	Provision_tenant_id types.Int64  `tfsdk:"provision_tenant_id"`
	//Provision_vm_name_next_value types.Int64  `tfsdk:"provision_vm_name_next_value"`
	Provision_vm_display_name types.String `tfsdk:"provision_vm_display_name"`
	Provision_url             types.String `tfsdk:"provision_url"`
	Provision_limits_enforce  types.Int64  `tfsdk:"provision_limits_enforce"`
	Mark_deletable            types.Int64  `tfsdk:"mark_deletable"`
	Center                    types.Object `tfsdk:"center"`
}

// attrTypes - return attribute types for this model
func (o awsProvisionModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"provision_on_off":          types.Int64Type,
		"provision_max":             types.Int64Type,
		"provision_vm_id":           types.Int64Type,
		"provision_server_id":       types.Int64Type,
		"provision_vm_name":         types.StringType,
		"provision_threshold":       types.Int64Type,
		"provision_tenant_id":       types.Int64Type,
		"provision_vm_display_name": types.StringType,
		"provision_url":             types.StringType,
		"provision_limits_enforce":  types.Int64Type,
		"mark_deletable":            types.Int64Type,
		"center":                    types.ObjectType{AttrTypes: awsCenterModel{}.attrTypes()},
	}
}

// defaultObject - return default object for this model representing the provision object
func (o awsProvisionModel) defaultObject() map[string]attr.Value {
	return map[string]attr.Value{
		"provision_on_off":          types.Int64Value(CONFIG_POOL_PROVISION_ON_OFF),
		"provision_max":             types.Int64Value(CONFIG_POOL_PROVISION_MAX),
		"provision_vm_id":           types.Int64Value(CONFIG_POOL_PROVISION_VM_ID),
		"provision_server_id":       types.Int64Value(CONFIG_POOL_PROVISION_SERVER_ID),
		"provision_vm_name":         types.StringValue(""),
		"provision_threshold":       types.Int64Value(CONFIG_POOL_PROVISION_THRESHOLD),
		"provision_tenant_id":       types.Int64Value(CONFIG_POOL_PROVISION_TENANT_ID),
		"provision_vm_display_name": types.StringValue(""),
		"provision_url":             types.StringValue(""),
		"provision_limits_enforce":  types.Int64Value(CONFIG_POOL_PROVISION_LIMITS_ENFORCE),
		"mark_deletable":            types.Int64Value(CONFIG_POOL_MARK_DELETABLE),
		"center":                    types.ObjectValueMust(awsCenterModel{}.attrTypes(), awsCenterModel{}.defaultObject()),
	}
}

// centerModel maps center schema data
type awsCenterModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	Provision_method types.String `tfsdk:"provision_method"`
	Aws_size         types.String `tfsdk:"aws_size"`
	Aws_iam_name     types.String `tfsdk:"aws_iam_name"`
	Aws_sub_net      types.String `tfsdk:"aws_sub_net"`
	Aws_sec_group    types.String `tfsdk:"aws_sec_group"`
	Aws_vpc_id       types.String `tfsdk:"aws_vpc_id"`
}

// attrTypes - return attribute types for this model
func (o awsCenterModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.Int64Type,
		"name":             types.StringType,
		"type":             types.StringType,
		"provision_method": types.StringType,
		"aws_size":         types.StringType,
		"aws_iam_name":     types.StringType,
		"aws_sub_net":      types.StringType,
		"aws_sec_group":    types.StringType,
		"aws_vpc_id":       types.StringType,
	}
}

// defaultObject - return default object for this model
func (o awsCenterModel) defaultObject() map[string]attr.Value {
	return map[string]attr.Value{
		"id":               types.Int64Value(0),
		"name":             types.StringValue(""),
		"type":             types.StringValue(""),
		"provision_method": types.StringValue("image"),
		"aws_size":         types.StringValue(""),
		"aws_iam_name":     types.StringValue(""),
		"aws_sub_net":      types.StringValue(""),
		"aws_sec_group":    types.StringValue(""),
		"aws_vpc_id":       types.StringValue(""),
	}
}

// attributesModel maps pool definition attribute schema data
type awsAttributesModel struct {
	Vm_table_field     types.String `tfsdk:"vm_table_field"`
	Ad_attribute_field types.String `tfsdk:"ad_attribute_field"`
	Vm_gpu_field       types.String `tfsdk:"vm_gpu_field"`
	Text_to_match      types.String `tfsdk:"text_to_match"`
	Condition_type     types.String `tfsdk:"condition_type"`
}

// attrTypes - return attribute types for this model
func (o awsAttributesModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vm_table_field":     types.StringType,
		"ad_attribute_field": types.StringType,
		"vm_gpu_field":       types.StringType,
		"text_to_match":      types.StringType,
		"condition_type":     types.StringType,
	}
}

// // defaultObject - return default object for this model
// func (o attributesModel) defaultObject() map[string]attr.Value {
// 	return map[string]attr.Value{
// 		"vm_table_field":     types.StringValue(""),
// 		"ad_attribute_field": types.StringValue(""),
// 		"vm_gpu_field":       types.StringValue(""),
// 		"text_to_match":      types.StringValue(""),
// 		"condition_type":     types.StringValue(""),
// 	}
// }

// common `Read` function for both data source and resource
func (o *awsPoolResourceModel) Read(ctx context.Context, client leostream.Client, diags *diag.Diagnostics, rtype string, id string) {
	//Pool CONFIG
	//get refreshed pool config value from Leostream API
	poolConfig, err := client.GetPool(id)

	if err != nil {
		diags.AddError(
			"Unable to read Pool Configuration",
			err.Error(),
		)
		return
	}

	// Map pool config to state
	o.Name = types.StringValue(poolConfig.Name)
	o.Display_name = types.StringValue(poolConfig.Display_name)
	o.Notes = types.StringValue(poolConfig.Notes)
	o.Running_desktops_threshold = types.Int64Value(poolConfig.Running_desktops_threshold)

	// Map pool definition to state
	var statePoolDefinition awsPoolDefinitionModel
	statePoolDefinition.Restrict_by = types.StringValue(poolConfig.Pool_definition.Restrict_by)
	statePoolDefinition.Pool_attribute_join = types.StringValue(poolConfig.Pool_definition.Pool_attribute_join)
	statePoolDefinition.Server_ids, *diags = types.ListValueFrom(ctx, types.Int64Type, poolConfig.Pool_definition.Server_ids)
	statePoolDefinition.Never_rogue = types.Int64Value(poolConfig.Pool_definition.Never_rogue)
	statePoolDefinition.Use_vmotion = types.Int64Value(poolConfig.Pool_definition.Use_vmotion)
	statePoolDefinition.Parent_pool_id = types.Int64Value(poolConfig.Pool_definition.Parent_pool_id)

	// Create a slice of attributesModel called statePoolDefinitionAttributes
	var statePoolDefinitionAttributes []awsAttributesModel
	// Loop through the poolConfig.Pool_definition.Attributes and assign the values to the stateAttributes
	for _, attribute := range poolConfig.Pool_definition.Attributes {
		var stateAttributes awsAttributesModel
		stateAttributes.Vm_table_field = types.StringValue(attribute.Vm_table_field)
		stateAttributes.Ad_attribute_field = types.StringValue(attribute.Ad_attribute_field)
		stateAttributes.Vm_gpu_field = types.StringValue(attribute.Vm_gpu_field)
		stateAttributes.Text_to_match = types.StringValue(attribute.Text_to_match)
		stateAttributes.Condition_type = types.StringValue(attribute.Condition_type)
		// Append the stateAttributes to the statePoolDefinitionAttributes
		statePoolDefinitionAttributes = append(statePoolDefinitionAttributes, stateAttributes)
	}

	// Assign the list to the statePoolDefinitionAttributes list value in the statePoolDefinition
	// convert to a list
	statePoolDefinition.Attributes, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: awsAttributesModel{}.attrTypes()}, statePoolDefinitionAttributes)

	//Add pool definition to pool model
	o.Pool_definition, _ = types.ObjectValueFrom(ctx, awsPoolDefinitionModel{}.attrTypes(), &statePoolDefinition)

	// Handle provision attribute
	var stateProvision awsProvisionModel
	stateProvision.Provision_on_off = types.Int64Value(poolConfig.Provision.Provision_on_off)
	stateProvision.Provision_max = types.Int64Value(poolConfig.Provision.Provision_max)
	stateProvision.Provision_vm_id = types.Int64Value(poolConfig.Provision.Provision_vm_id)
	stateProvision.Provision_server_id = types.Int64Value(poolConfig.Provision.Provision_server_id)
	stateProvision.Provision_vm_name = types.StringValue(poolConfig.Provision.Provision_vm_name)
	stateProvision.Provision_threshold = types.Int64Value(poolConfig.Provision.Provision_threshold)
	stateProvision.Provision_tenant_id = types.Int64Value(poolConfig.Provision.Provision_tenant_id)
	//stateProvision.Provision_vm_name_next_value = types.Int64Value(poolConfig.Provision.Provision_vm_name_next_value)
	stateProvision.Provision_vm_display_name = types.StringValue(poolConfig.Provision.Provision_vm_display_name)
	stateProvision.Provision_url = types.StringValue(poolConfig.Provision.Provision_url)
	stateProvision.Provision_limits_enforce = types.Int64Value(poolConfig.Provision.Provision_limits_enforce)
	stateProvision.Mark_deletable = types.Int64Value(poolConfig.Provision.Mark_deletable)

	// if poolConfig.Provision.Center is not null, then unpack the center attributes
	if poolConfig.Provision.Center != nil {
		// Handle center attribute
		var stateCenter awsCenterModel
		stateCenter.ID = types.Int64Value(poolConfig.Provision.Center.ID)
		stateCenter.Name = types.StringValue(poolConfig.Provision.Center.Name)
		stateCenter.Type = types.StringValue(poolConfig.Provision.Center.Type)
		stateCenter.Provision_method = types.StringValue(poolConfig.Provision.Center.Provision_method)
		stateCenter.Aws_size = types.StringValue(poolConfig.Provision.Center.Aws_size)
		stateCenter.Aws_iam_name = types.StringValue(poolConfig.Provision.Center.Aws_iam_name)
		stateCenter.Aws_sub_net = types.StringValue(poolConfig.Provision.Center.Aws_sub_net)
		stateCenter.Aws_sec_group = types.StringValue(poolConfig.Provision.Center.Aws_sec_group)
		stateCenter.Aws_vpc_id = types.StringValue(poolConfig.Provision.Center.Aws_vpc_id)

		// Add center to provision model
		stateProvision.Center, _ = types.ObjectValueFrom(ctx, awsCenterModel{}.attrTypes(), &stateCenter)

	}

	// Add provision to pool model
	o.Provision, _ = types.ObjectValueFrom(ctx, awsProvisionModel{}.attrTypes(), &stateProvision)

}

// `Create` function for the resource
func (r *awsPoolResource) CreateNested(ctx context.Context, plan *awsPoolResourceModel, state *awsPoolResourceModel, diags *diag.Diagnostics) *leostream.PoolsStored {
	// Pool CONFIG

	// Instantiate empty object for storing plan data
	var poolConfig leostream.Pool

	// Populate pool config from plan
	poolConfig.Name = plan.Name.ValueString()
	poolConfig.Display_name = plan.Display_name.ValueString()
	poolConfig.Notes = plan.Notes.ValueString()
	poolConfig.Running_desktops_threshold = plan.Running_desktops_threshold.ValueInt64()

	// Unpack nested attributes from plan for the pool definition
	var planPoolDefinition awsPoolDefinitionModel
	*diags = plan.Pool_definition.As(ctx, &planPoolDefinition, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}

	// Instantiate empty object for storing plan data for the pool definition object in the pool config
	var poolDefinitionConfig leostream.PoolDefinition

	// Populate pool_definition field restrict_by in empty object from plan
	poolDefinitionConfig.Restrict_by = planPoolDefinition.Restrict_by.ValueString()

	// Populate pool_definition Pool_attribute_join field in empty object from plan
	poolDefinitionConfig.Pool_attribute_join = planPoolDefinition.Pool_attribute_join.ValueString()

	// Populate pool_definition Server_ids field in empty object from plan (but only if it is not empty)
	// todo: what is the default value for server_ids? 0?
	//if len(planPoolDefinition.Server_ids.Elements()) > 0 {
	for _, server_id := range planPoolDefinition.Server_ids.Elements() {
		// Convert the server_id to an int64 using an intermediary variable
		server_id_int64, _ := strconv.ParseInt(server_id.String(), 10, 32)
		// Append the server_id_int64 to the poolDefinitionConfig.Server_ids
		poolDefinitionConfig.Server_ids = append(poolDefinitionConfig.Server_ids, server_id_int64)
	}
	*diags = planPoolDefinition.Server_ids.ElementsAs(ctx, &poolDefinitionConfig.Server_ids, false)
	if diags.HasError() {
		return nil
	}
	//}
	// Populate pool_definition Never_rogue field in empty object from plan
	poolDefinitionConfig.Never_rogue = planPoolDefinition.Never_rogue.ValueInt64()

	// Populate pool_definition  Use_vmotion in empty object from plan
	poolDefinitionConfig.Use_vmotion = planPoolDefinition.Use_vmotion.ValueInt64()

	// Populate pool_definition Parent_pool_id in empty object from plan
	poolDefinitionConfig.Parent_pool_id = planPoolDefinition.Parent_pool_id.ValueInt64()

	// Instantiate empty object for storing plan data for the attributes object in the pool definition object in the pool config
	var planAttributes []awsAttributesModel

	// Populate pool_definition Attributes field in empty object from plan (but only if it exists)
	// todo: what is the default value for Attributes? empty null object?
	if !planPoolDefinition.Attributes.IsNull() {

		*diags = planPoolDefinition.Attributes.ElementsAs(ctx, &planAttributes, false)
		if diags.HasError() {
			return nil
		}
	}

	//Object for storing plan data for the attributes list in the pooldefinition object of the pool config
	var attributesConfig []leostream.PoolAttributes

	// Loop through the planAttributes and assign the values to the attributesConfig
	for _, attribute := range planAttributes {
		var attributeConfig leostream.PoolAttributes
		attributeConfig.Vm_table_field = attribute.Vm_table_field.ValueString()
		attributeConfig.Ad_attribute_field = attribute.Ad_attribute_field.ValueString()
		attributeConfig.Vm_gpu_field = attribute.Vm_gpu_field.ValueString()
		attributeConfig.Text_to_match = attribute.Text_to_match.ValueString()
		attributeConfig.Condition_type = attribute.Condition_type.ValueString()

		// Append the attributeConfig to the attributesConfig
		attributesConfig = append(attributesConfig, attributeConfig)

	}

	// Assign the attributesConfig to the poolDefinitionConfig.Attributes
	poolDefinitionConfig.Attributes = attributesConfig

	// Assign the pool definition config to the pool config
	poolConfig.Pool_definition = &poolDefinitionConfig

	// Unpack nested attributes from plan for the provision object in the pool config
	var planProvision awsProvisionModel
	*diags = plan.Provision.As(ctx, &planProvision, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}
	// Object for storing plan data for the provision object in the pool config
	var provisionConfig leostream.Provision

	// Populate pool provision config from plan
	provisionConfig.Provision_on_off = planProvision.Provision_on_off.ValueInt64()
	provisionConfig.Provision_max = planProvision.Provision_max.ValueInt64()
	provisionConfig.Provision_vm_id = planProvision.Provision_vm_id.ValueInt64()
	provisionConfig.Provision_server_id = planProvision.Provision_server_id.ValueInt64()
	provisionConfig.Provision_vm_name = planProvision.Provision_vm_name.ValueString()
	provisionConfig.Provision_threshold = planProvision.Provision_threshold.ValueInt64()
	provisionConfig.Provision_tenant_id = planProvision.Provision_tenant_id.ValueInt64()
	//_value = planProvision.Provision_vm_name_next_value.ValueInt64()
	provisionConfig.Provision_vm_display_name = planProvision.Provision_vm_display_name.ValueString()
	provisionConfig.Provision_url = planProvision.Provision_url.ValueString()
	provisionConfig.Provision_limits_enforce = planProvision.Provision_limits_enforce.ValueInt64()
	provisionConfig.Mark_deletable = planProvision.Mark_deletable.ValueInt64()

	// Object for storing plan data for the center object in the provision object of the pool config
	var centerConfig leostream.PoolAwsCenter

	var planCenter awsCenterModel
	if !planProvision.Center.IsNull() {
		*diags = planProvision.Center.As(ctx, &planCenter, basetypes.ObjectAsOptions{})
	}

	centerConfig.ID = int64(planCenter.ID.ValueInt64())
	centerConfig.Name = planCenter.Name.ValueString()
	centerConfig.Type = planCenter.Type.ValueString()
	centerConfig.Provision_method = planCenter.Provision_method.ValueString()
	centerConfig.Aws_size = planCenter.Aws_size.ValueString()
	centerConfig.Aws_iam_name = planCenter.Aws_iam_name.ValueString()
	centerConfig.Aws_sub_net = planCenter.Aws_sub_net.ValueString()
	centerConfig.Aws_sec_group = planCenter.Aws_sec_group.ValueString()
	centerConfig.Aws_vpc_id = planCenter.Aws_vpc_id.ValueString()

	provisionConfig.Center = &centerConfig

	poolConfig.Provision = &provisionConfig

	// Create new pool
	PoolsStored, err := r.client.CreatePool(poolConfig, nil)

	if err != nil {
		diags.AddError(
			"Unable to Create Pool",
			err.Error(),
		)
		return nil
	} else {
		return PoolsStored
	}
}

// `Update` function for the resource
func (r *awsPoolResource) UpdateNested(ctx context.Context, plan *awsPoolResourceModel, state *awsPoolResourceModel, diags *diag.Diagnostics) *leostream.PoolsStored {
	// Pool CONFIG

	// Instantiate empty object for storing plan data
	var poolConfig leostream.Pool

	// Populate pool config from plan
	poolConfig.Name = plan.Name.ValueString()
	poolConfig.Display_name = plan.Display_name.ValueString()
	poolConfig.Notes = plan.Notes.ValueString()
	poolConfig.Running_desktops_threshold = plan.Running_desktops_threshold.ValueInt64()

	// Unpack nested attributes from plan for the pool definition
	var planPoolDefinition awsPoolDefinitionModel
	*diags = plan.Pool_definition.As(ctx, &planPoolDefinition, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}
	// Instantiate empty object for storing plan data for the pool definition object in the pool config
	var poolDefinitionConfig leostream.PoolDefinition

	// Populate pool definition config from plan
	poolDefinitionConfig.Restrict_by = planPoolDefinition.Restrict_by.ValueString()
	if len(planPoolDefinition.Server_ids.Elements()) > 0 {
		*diags = planPoolDefinition.Server_ids.ElementsAs(ctx, &poolDefinitionConfig.Server_ids, false)
		if diags.HasError() {
			return nil
		}
	}

	// Assign the value of planPoolDefinition.Pool_attribute_join to poolDefinitionConfig.Pool_attribute_join
	poolDefinitionConfig.Pool_attribute_join = planPoolDefinition.Pool_attribute_join.ValueString()
	// Assign the value of planPoolDefinition.Never_rogue to poolDefinitionConfig.Never_rogue
	poolDefinitionConfig.Never_rogue = planPoolDefinition.Never_rogue.ValueInt64()
	// Assign the value of planPoolDefinition.Use_vmotion to poolDefinitionConfig.Use_vmotion
	poolDefinitionConfig.Use_vmotion = planPoolDefinition.Use_vmotion.ValueInt64()
	// Assign the value of planPoolDefinition.Parent_pool_id to poolDefinitionConfig.Parent_pool_id
	poolDefinitionConfig.Parent_pool_id = planPoolDefinition.Parent_pool_id.ValueInt64()

	// Instantiate empty object for storing plan data for the attributes object in the pool definition object in the pool config
	var planAttributes []awsAttributesModel

	// Populate pool_definition Attributes field in empty object from plan (but only if it exists)
	// todo: what is the default value for Attributes? empty null object?
	if !planPoolDefinition.Attributes.IsNull() {

		*diags = planPoolDefinition.Attributes.ElementsAs(ctx, &planAttributes, false)
		if diags.HasError() {
			return nil
		}
	}

	//Object for storing plan data for the attributes list in the pooldefinition object of the pool config
	var attributesConfig []leostream.PoolAttributes

	// Loop through the planAttributes and assign the values to the attributesConfig
	for _, attribute := range planAttributes {
		var attributeConfig leostream.PoolAttributes
		attributeConfig.Vm_table_field = attribute.Vm_table_field.ValueString()
		attributeConfig.Ad_attribute_field = attribute.Ad_attribute_field.ValueString()
		attributeConfig.Vm_gpu_field = attribute.Vm_gpu_field.ValueString()
		attributeConfig.Text_to_match = attribute.Text_to_match.ValueString()
		attributeConfig.Condition_type = attribute.Condition_type.ValueString()

		// Append the attributeConfig to the attributesConfig
		attributesConfig = append(attributesConfig, attributeConfig)

	}

	// Assign the attributesConfig to the poolDefinitionConfig.Attributes
	poolDefinitionConfig.Attributes = attributesConfig

	// Assign the pool definition config to the pool config
	poolConfig.Pool_definition = &poolDefinitionConfig

	// unpack nested attributes from plan for the provision object in the pool config
	var planProvision awsProvisionModel
	*diags = plan.Provision.As(ctx, &planProvision, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}

	// instantiate empty object for storing plan data for the provision object in the pool config
	var provisionConfig leostream.Provision
	// populate pool definition config from plan
	provisionConfig.Provision_on_off = planProvision.Provision_on_off.ValueInt64()
	provisionConfig.Provision_max = planProvision.Provision_max.ValueInt64()
	provisionConfig.Provision_vm_id = planProvision.Provision_vm_id.ValueInt64()
	provisionConfig.Provision_server_id = planProvision.Provision_server_id.ValueInt64()
	provisionConfig.Provision_vm_name = planProvision.Provision_vm_name.ValueString()
	provisionConfig.Provision_threshold = planProvision.Provision_threshold.ValueInt64()
	provisionConfig.Provision_tenant_id = planProvision.Provision_tenant_id.ValueInt64()
	//provisionConfig.Provision_vm_name_next_value = planProvision.Provision_vm_name_next_value.ValueInt64()
	provisionConfig.Provision_vm_display_name = planProvision.Provision_vm_display_name.ValueString()
	provisionConfig.Provision_url = planProvision.Provision_url.ValueString()
	provisionConfig.Provision_limits_enforce = planProvision.Provision_limits_enforce.ValueInt64()
	provisionConfig.Mark_deletable = planProvision.Mark_deletable.ValueInt64()

	// unpack nested attributes from plan for the center object in provision object in the pool config
	var centerConfig leostream.PoolAwsCenter

	var planCenter awsCenterModel
	if !planProvision.Center.IsNull() {
		*diags = planProvision.Center.As(ctx, &planCenter, basetypes.ObjectAsOptions{})
	}

	centerConfig.ID = int64(planCenter.ID.ValueInt64())
	centerConfig.Name = planCenter.Name.ValueString()
	centerConfig.Type = planCenter.Type.ValueString()
	centerConfig.Provision_method = planCenter.Provision_method.ValueString()
	centerConfig.Aws_size = planCenter.Aws_size.ValueString()
	centerConfig.Aws_iam_name = planCenter.Aws_iam_name.ValueString()
	centerConfig.Aws_sub_net = planCenter.Aws_sub_net.ValueString()
	centerConfig.Aws_sec_group = planCenter.Aws_sec_group.ValueString()
	centerConfig.Aws_vpc_id = planCenter.Aws_vpc_id.ValueString()

	provisionConfig.Center = &centerConfig

	poolConfig.Provision = &provisionConfig

	tflog.Info(ctx, "Performing Update via pool common")

	// Update pool
	PoolsStored, err := r.client.UpdatePool(plan.ID.ValueString(), poolConfig, nil)

	if err != nil {
		diags.AddError(
			"Unable to Create Pool",
			err.Error(),
		)
		return nil
	} else {
		return PoolsStored
	}
}
