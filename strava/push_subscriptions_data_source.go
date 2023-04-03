package strava

import (
	"context"

	"github.com/floydspace/strava-webhook-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &pushSubscriptionsDataSource{}
	_ datasource.DataSourceWithConfigure = &pushSubscriptionsDataSource{}
)

// NewPushSubscriptionsDataSource is a helper function to simplify the provider implementation.
func NewPushSubscriptionsDataSource() datasource.DataSource {
	return &pushSubscriptionsDataSource{}
}

// pushSubscriptionsDataSource is the data source implementation.
type pushSubscriptionsDataSource struct {
	client *strava.Client
}

// pushSubscriptionsDataSourceModel maps the data source schema data.
type pushSubscriptionsDataSourceModel struct {
	ID                types.String             `tfsdk:"id"`
	PushSubscriptions []pushSubscriptionsModel `tfsdk:"push_subscriptions"`
}

// pushSubscriptionsModel maps pushSubscriptions schema data.
type pushSubscriptionsModel struct {
	ID            types.Int64  `tfsdk:"id"`
	ResourceState types.Int64  `tfsdk:"resource_state"`
	ApplicationID types.Int64  `tfsdk:"application_id"`
	CallbackURL   types.String `tfsdk:"callback_url"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *pushSubscriptionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_subscriptions"
}

// Schema defines the schema for the data source.
func (d *pushSubscriptionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of push subscriptions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"push_subscriptions": schema.ListNestedAttribute{
				Description: "List of push subscriptions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Push subscription ID.",
							Computed:    true,
						},
						"resource_state": schema.Int64Attribute{
							Description: "State of the subscription.",
							Computed:    true,
						},
						"application_id": schema.Int64Attribute{
							Description: "Strava API application ID.",
							Computed:    true,
						},
						"callback_url": schema.StringAttribute{
							Description: "Address where webhook events will be sent; maximum length of 255 characters.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Date and time the subscription was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Date and time the subscription was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *pushSubscriptionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pushSubscriptionsDataSourceModel

	pushSubscriptions, err := d.client.GetAllSubscriptions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Strava Subscriptions",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, subscription := range *pushSubscriptions {
		subscriptionState := pushSubscriptionsModel{
			ID:            types.Int64Value(int64(subscription.ID)),
			ResourceState: types.Int64Value(int64(subscription.ResourceState)),
			ApplicationID: types.Int64Value(int64(subscription.ApplicationID)),
			CallbackURL:   types.StringValue(subscription.CallbackURL),
			CreatedAt:     types.StringValue(subscription.CreatedAt),
			UpdatedAt:     types.StringValue(subscription.UpdatedAt),
		}

		state.PushSubscriptions = append(state.PushSubscriptions, subscriptionState)
	}

	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *pushSubscriptionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*strava.Client)
}
