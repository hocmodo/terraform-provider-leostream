// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.hocmodo.nl/community/leostream-client-go"
)

// poolAssignmentResourceModel maps the resource schema data.
type poolAssignmentResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Pool_id               types.Int64  `tfsdk:"pool_id"`
	Policy_id             types.Int64  `tfsdk:"policy_id"`
	Offer_filter          types.String `tfsdk:"offer_filter"`
	Offer_filter_json     types.Object `tfsdk:"offer_filter_json"`
	Plan_protocol_id      types.Int64  `tfsdk:"plan_protocol_id"`
	Plan_power_control_id types.Int64  `tfsdk:"plan_power_control_id"`
	Plan_release_id       types.Int64  `tfsdk:"plan_release_id"`
	Offer_quantity        types.Int64  `tfsdk:"offer_quantity"`
	Display_mode          types.String `tfsdk:"display_mode"`
	Start_if_stopped      types.Int64  `tfsdk:"start_if_stopped"`
}

// offerFilterModel maps filtering schema data
type offerFilterJsonModel struct {
	Join    types.String `tfsdk:"join"`
	Filters types.List   `tfsdk:"filters"`
}

// filterModel maps filtering schema data
type filterModel struct {
	Offer_filter_attribute types.String `tfsdk:"offer_filter_attribute"`
	Offer_filter_condition types.String `tfsdk:"offer_filter_condition"`
	Offer_filter_value     types.String `tfsdk:"offer_filter_value"`
}

// attrTypes - return attribute types for this model
func (o filterModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"offer_filter_attribute": types.StringType,
		"offer_filter_condition": types.StringType,
		"offer_filter_value":     types.StringType,
	}
}

// attrTypes - return attribute types for this model
func (o offerFilterJsonModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"join":    types.StringType,
		"filters": types.ListType{ElemType: types.ObjectType{AttrTypes: filterModel{}.attrTypes()}},
	}
}

// defaultObject - return default object for this model
func (o offerFilterJsonModel) defaultObject() map[string]attr.Value {
	return map[string]attr.Value{
		"join":    types.StringValue("O"),
		"filters": types.ListNull(types.ObjectType{AttrTypes: filterModel{}.attrTypes()}),
	}
}

// common `Read` function for both data source and resource
func (o *poolAssignmentResourceModel) Read(ctx context.Context, client leostream.Client, diags *diag.Diagnostics, rtype string, policy_id string, id string) {

	tflog.Info(ctx, "Performing read nested")

	//poolassignment CONFIG
	//get refreshed poolassignment config value from Leostream API
	poolassignmentConfig, err := client.GetPoolAssignment(policy_id, id)

	if err != nil {
		diags.AddError(
			"Unable to read poolassignment configuration",
			err.Error(),
		)
		return
	}

	o.ID = types.StringValue(strconv.FormatInt(poolassignmentConfig.ID, 10))

	// Map poolassignment fields to state
	o.Display_mode = types.StringValue(poolassignmentConfig.Display_mode)
	o.Offer_filter = types.StringValue(poolassignmentConfig.Offer_filter)
	o.Offer_quantity = types.Int64Value(poolassignmentConfig.Offer_quantity)
	o.Plan_power_control_id = types.Int64Value(poolassignmentConfig.Plan_power_control_id)
	o.Plan_protocol_id = types.Int64Value(poolassignmentConfig.Plan_protocol_id)
	o.Plan_release_id = types.Int64Value(poolassignmentConfig.Plan_release_id)
	o.Policy_id = types.Int64Value(poolassignmentConfig.Policy_id)
	o.Pool_id = types.Int64Value(poolassignmentConfig.Pool_id)
	o.Start_if_stopped = types.Int64Value(poolassignmentConfig.Start_if_stopped)

	// Map offer filter json to state
	var stateOfferFilter offerFilterJsonModel
	stateOfferFilter.Join = types.StringValue(poolassignmentConfig.Offer_filter_json.Join)

	// Create a slice of filterModel called stateOfferJsonFilters
	var stateOfferJsonFilters []filterModel
	// Loop through the filters list and assign the values to filterModel
	for _, filter := range poolassignmentConfig.Offer_filter_json.Filters {
		var stateFilters filterModel
		stateFilters.Offer_filter_attribute = types.StringValue(filter.Offer_filter_attribute)
		stateFilters.Offer_filter_condition = types.StringValue(filter.Offer_filter_condition)
		stateFilters.Offer_filter_value = types.StringValue(filter.Offer_filter_value)
		// Append the stateFilters to the stateOfferJsonFilters
		stateOfferJsonFilters = append(stateOfferJsonFilters, stateFilters)
	}

	// Assign the list to the Filters value in the stateOfferFilter
	stateOfferFilter.Filters, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: filterModel{}.attrTypes()}, stateOfferJsonFilters)

	//Add stateOfferFilter to poolassignment model
	o.Offer_filter_json, _ = types.ObjectValueFrom(ctx, offerFilterJsonModel{}.attrTypes(), &stateOfferFilter)

}

// `Create` function for the resource
func (r *poolAssignmentResource) CreateNested(ctx context.Context, plan *poolAssignmentResourceModel, state *poolAssignmentResourceModel, diags *diag.Diagnostics, policy_id string) *leostream.PoolAssignmentssStored {

	tflog.Info(ctx, "Performing create nested")

	// center CONFIG

	// Instantiate empty object for storing plan data
	var poolassignmentConfig leostream.PoolAssignment

	// Instantiate empty object for storing plan data for the offer_filter_json object in the poolassignment config
	var offerFilterJsonConfig leostream.PoolAssignmentFilters

	// Populate poolassignment config from plan
	poolassignmentConfig.Policy_id, _ = strconv.ParseInt(policy_id, 10, 64)
	poolassignmentConfig.Pool_id = plan.Pool_id.ValueInt64()
	poolassignmentConfig.Offer_filter = plan.Offer_filter.ValueString()
	poolassignmentConfig.Plan_protocol_id = plan.Plan_protocol_id.ValueInt64()
	poolassignmentConfig.Plan_power_control_id = plan.Plan_power_control_id.ValueInt64()
	poolassignmentConfig.Plan_release_id = plan.Plan_release_id.ValueInt64()
	poolassignmentConfig.Offer_quantity = plan.Offer_quantity.ValueInt64()
	poolassignmentConfig.Display_mode = plan.Display_mode.ValueString()
	poolassignmentConfig.Start_if_stopped = plan.Start_if_stopped.ValueInt64()

	// Unpack nested attributes from plan for the pool offer_filter_json
	var planOfferFilterJson offerFilterJsonModel

	*diags = plan.Offer_filter_json.As(ctx, &planOfferFilterJson, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}

	offerFilterJsonConfig.Join = planOfferFilterJson.Join.ValueString()

	// Instantiate empty object for storing plan data for the attributes object in the pool definition object in the pool config
	var planFilters []filterModel

	if !planOfferFilterJson.Filters.IsNull() {

		*diags = planOfferFilterJson.Filters.ElementsAs(ctx, &planFilters, false)
		if diags.HasError() {
			return nil
		}
	}

	//Object for storing plan data for the filters list in the filtersConfig object of the poolassignment config
	var filtersConfig []leostream.Filters

	// Loop through the planFilters and assign the values to the filterConfig instance
	for _, filter := range planFilters {
		var filterConfig leostream.Filters
		filterConfig.Offer_filter_attribute = filter.Offer_filter_attribute.ValueString()
		filterConfig.Offer_filter_condition = filter.Offer_filter_condition.ValueString()
		filterConfig.Offer_filter_value = filter.Offer_filter_value.ValueString()

		// Append the filterConfig instance to the list
		filtersConfig = append(filtersConfig, filterConfig)

	}

	// Assign the center definition config to the center config
	offerFilterJsonConfig.Filters = filtersConfig

	poolassignmentConfig.Offer_filter_json = &offerFilterJsonConfig

	// Create new poolassignment
	poolAssignmentsStored, err := r.client.CreatePoolAssignment(poolassignmentConfig, policy_id, nil)

	if err != nil {
		diags.AddError(
			"Unable to create poolassignment",
			err.Error(),
		)
		return nil
	} else {
		return poolAssignmentsStored
	}
}

// `Update` function for the resource
func (r *poolAssignmentResource) UpdateNested(ctx context.Context, plan *poolAssignmentResourceModel, state *poolAssignmentResourceModel, diags *diag.Diagnostics, policy_id string) *leostream.PoolAssignmentssStored {

	tflog.Info(ctx, "Performing update nested")

	// center CONFIG

	// Instantiate empty object for storing plan data
	var poolassignmentConfig leostream.PoolAssignment

	// Instantiate empty object for storing plan data for the offer_filter_json object in the poolassignment config
	var offerFilterJsonConfig leostream.PoolAssignmentFilters

	// Populate poolassignment config from plan
	poolassignmentConfig.Policy_id, _ = strconv.ParseInt(policy_id, 10, 64)
	poolassignmentConfig.Pool_id = plan.Pool_id.ValueInt64()
	poolassignmentConfig.Offer_filter = plan.Offer_filter.ValueString()
	poolassignmentConfig.Plan_protocol_id = plan.Plan_protocol_id.ValueInt64()
	poolassignmentConfig.Plan_power_control_id = plan.Plan_power_control_id.ValueInt64()
	poolassignmentConfig.Plan_release_id = plan.Plan_release_id.ValueInt64()
	poolassignmentConfig.Offer_quantity = plan.Offer_quantity.ValueInt64()
	poolassignmentConfig.Display_mode = plan.Display_mode.ValueString()
	poolassignmentConfig.Start_if_stopped = plan.Start_if_stopped.ValueInt64()

	// Unpack nested attributes from plan for the pool offer_filter_json
	var planOfferFilterJson offerFilterJsonModel
	*diags = plan.Offer_filter_json.As(ctx, &planOfferFilterJson, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil
	}

	offerFilterJsonConfig.Join = planOfferFilterJson.Join.ValueString()

	// Instantiate empty object for storing plan data for the filters list in the filterJson object in the poolassignment config
	var planFilters []filterModel

	// Populate poolassignment_definition json filter field in empty object from plan (but only if it exists)
	if !planOfferFilterJson.Filters.IsNull() {

		*diags = planOfferFilterJson.Filters.ElementsAs(ctx, &planFilters, false)
		if diags.HasError() {
			return nil
		}
	}

	//Object for storing plan data for the filter list in the Offer_filter_json object of the poolassignment
	var filtersConfig []leostream.Filters

	// Loop through the planAttributes and assign the values to the attributesConfig
	for _, filter := range planFilters {
		var filterConfig leostream.Filters
		filterConfig.Offer_filter_attribute = filter.Offer_filter_attribute.ValueString()
		filterConfig.Offer_filter_condition = filter.Offer_filter_condition.ValueString()
		filterConfig.Offer_filter_value = filter.Offer_filter_value.ValueString()

		// Append the filterConfig to the list
		filtersConfig = append(filtersConfig, filterConfig)

	}

	// Assign the filters  config to the offerFilterJsonConfig
	offerFilterJsonConfig.Filters = filtersConfig

	poolassignmentConfig.Offer_filter_json = &offerFilterJsonConfig

	// Update pool
	PoolAssignmentsStored, err := r.client.UpdatePoolAssignment(plan.ID.ValueString(), poolassignmentConfig, policy_id, nil)

	if err != nil {
		diags.AddError(
			"Unable to update poolassignment",
			err.Error(),
		)
		return nil
	} else {
		return PoolAssignmentsStored
	}

}
