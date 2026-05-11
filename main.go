package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/uladzimirSTR/terraform-provider-trino/internal/provider"
)

var version string = "dev"
var name string = "trino"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/ulstr/trino",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version, name), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
