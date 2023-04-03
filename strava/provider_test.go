package strava

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Strava client is properly configured.
	// It is also possible to use the STRAVA_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "strava" {
	client_id     = "5"
	client_secret = "7b2946535949ae70f015d696d8ac602830ece412"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"strava": providerserver.NewProtocol6WithError(New()),
	}
)
