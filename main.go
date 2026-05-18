package main

import (
	"context"
	"flag"

	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set this to true if you want to debug the code using device")
	flag.Parse()

	ctx := context.Background()

	providerserver.Serve(ctx, taskmanager.New, providerserver.ServeOpts{
		Debug:   false,
		Address: "terraform.examplr.com/local/taskmanager",
	})
}
