package strava

import (
	"context"
	"os"

	"github.com/floydspace/strava-webhook-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &stravaProvider{}
)

// stravaProvider is the provider implementation.
type stravaProvider struct{}

// stravaProviderModel maps provider schema data to a Go type.
type stravaProviderModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &stravaProvider{}
}

// Metadata returns the provider type name.
func (p *stravaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "strava"
}

// Schema defines the provider-level schema for configuration data.
func (p *stravaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Strava.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "Strava API application ID. May also be provided via the STRAVA_CLIENT_ID environment variable.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Strava API application secret. May also be provided via the STRAVA_CLIENT_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a strava API client for data sources and resources.
func (p *stravaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Strava client")

	// Retrieve provider data from configuration
	var config stravaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Strava API Client ID",
			"The provider cannot create the Strava API client as there is an unknown configuration value for the Strava API Client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STRAVA_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Strava API Client Secret",
			"The provider cannot create the Strava API client as there is an unknown configuration value for the Strava API Client Secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STRAVA_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	clientId := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if clientId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Strava API Client ID",
			"The provider cannot create the Strava API client as there is a missing or empty value for the Strava API Client ID. "+
				"Set the client_id value in the configuration or use the STRAVA_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Strava API Client Secret",
			"The provider cannot create the Strava API client as there is a missing or empty value for the Strava API Client Secret. "+
				"Set the client_secret value in the configuration or use the STRAVA_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "strava_client_id", clientId)
	ctx = tflog.SetField(ctx, "strava_client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "strava_client_secret")

	tflog.Debug(ctx, "Creating Strava client")

	// Create a new Strava client using the configuration values
	client, err := strava.NewClient(nil, &clientId, &clientSecret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Strava API Client",
			"An unexpected error occurred when creating the Strava API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Strava Client Error: "+err.Error(),
		)
		return
	}

	// Make the Strava client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Strava client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *stravaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPushSubscriptionsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *stravaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPushSubscriptionResource,
	}
}
