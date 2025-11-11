// Copyright (c) HashiCorp, Inc.

package main

import (
	"fmt"
	"context"
	"flag"
	"terraform-provider-leostream/leostream"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"bufio"
    "os"
    "os/exec"
    "strings"
)

func main() {

	var debug bool
	var userInput string
	// CodeQL test: hardcoded credential
	var secret = "12345"
	fmt.Printf("secret: %v\n", secret)


	// Prefer arg if present, otherwise prompt.
    if len(os.Args) > 1 {
        userInput = os.Args[1]
    } else {
        fmt.Print("Enter input: ")
        reader := bufio.NewReader(os.Stdin)
        line, _ := reader.ReadString('\n')
        userInput = strings.TrimSpace(line)
    }
	cmd := exec.Command("sh", "-c", userInput) // command injection
	out,_ := cmd.CombinedOutput()
	fmt.Println(string(out))

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	providerserver.Serve(context.Background(), leostream.New, providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/hocmodo/leostream",
	})
}
