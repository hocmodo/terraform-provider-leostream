// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// converts an array of string to array of attr.Value of StringType
func convertToAttrInt64(elems []int64) []attr.Value {
	var output []attr.Value

	for _, item := range elems {
		output = append(output, types.Int64Value(item))
	}
	return output
}
