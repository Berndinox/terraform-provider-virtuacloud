package main

import (
	"context"
	"flag"
	"log"

	"github.com/Berndinox/tf-provider-virtua-cloud/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "0.1.0"
)

func main() {
	debugFlag := flag.Bool("debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.opentofu.org/Berndinox/virtuacloud",
		Debug:   *debugFlag,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
