// Copyright (c) HashiCorp, Inc.

package main

import (
	"context"
	"flag"
	"terraform-provider-leostream/leostream"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	providerserver.Serve(context.Background(), leostream.New, providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/hocmodo/leostream",
	})
}
