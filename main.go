package main

import (
	"context"

	"github.com/floydspace/terraform-provider-strava/strava"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name strava

func main() {
	providerserver.Serve(context.Background(), strava.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/floydspace/strava",
	})
}
