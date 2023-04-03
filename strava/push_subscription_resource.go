package strava

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/floydspace/strava-webhook-client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pushSubscriptionResource{}
	_ resource.ResourceWithConfigure   = &pushSubscriptionResource{}
	_ resource.ResourceWithImportState = &pushSubscriptionResource{}
)

// NewPushSubscriptionResource is a helper function to simplify the provider implementation.
func NewPushSubscriptionResource() resource.Resource {
	return &pushSubscriptionResource{}
}

// pushSubscriptionResource is the resource implementation.
type pushSubscriptionResource struct {
	client *strava.Client
}

// pushSubscriptionResourceModel maps the resource schema data.
type pushSubscriptionResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	ResourceState types.Int64  `tfsdk:"resource_state"`
	ApplicationID types.Int64  `tfsdk:"application_id"`
	CallbackURL   types.String `tfsdk:"callback_url"`
	VerifyToken   types.String `tfsdk:"verify_token"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *pushSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_subscription"
}

// Schema defines the schema for the resource.
func (r *pushSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Strava push subscription.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Push subscription ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update to the push subscription.",
				Computed:    true,
			},
			"resource_state": schema.Int64Attribute{
				Description: "State of the push subscription.",
				Computed:    true,
			},
			"application_id": schema.Int64Attribute{
				Description: "Strava API application ID.",
				Computed:    true,
			},
			"callback_url": schema.StringAttribute{
				Description: "Address where webhook events will be sent; maximum length of 255 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"verify_token": schema.StringAttribute{
				Description: "String chosen by the application owner for client security. An identical string will be included in the validation request made by Strava's subscription service.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	}
}

// Create a new resource
func (r *pushSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pushSubscriptionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subscription
	pushSubscription, err := r.client.CreateSubscription(strava.SubscriptionItem{
		CallbackURL: plan.CallbackURL.ValueString(),
		VerifyToken: plan.VerifyToken.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subscription",
			"Could not create subscription, unexpected error: "+err.Error(),
		)
		return
	}

	pushSubscription, err = r.client.GetSubscription(int(pushSubscription.ID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created subscription",
			"Could not read created subscription, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(pushSubscription.ID))
	plan.ResourceState = types.Int64Value(int64(pushSubscription.ResourceState))
	plan.ApplicationID = types.Int64Value(int64(pushSubscription.ApplicationID))
	plan.CallbackURL = types.StringValue(pushSubscription.CallbackURL)
	plan.CreatedAt = types.StringValue(pushSubscription.CreatedAt)
	plan.UpdatedAt = types.StringValue(pushSubscription.UpdatedAt)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *pushSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pushSubscriptionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed push subscription value from Strava
	pushSubscription, err := r.client.GetSubscription(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Strava Subscription",
			"Could not read Strava subscription ID "+string(rune(state.ID.ValueInt64()))+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ResourceState = types.Int64Value(int64(pushSubscription.ResourceState))
	state.ApplicationID = types.Int64Value(int64(pushSubscription.ApplicationID))
	state.CallbackURL = types.StringValue(pushSubscription.CallbackURL)
	state.CreatedAt = types.StringValue(pushSubscription.CreatedAt)
	state.UpdatedAt = types.StringValue(pushSubscription.UpdatedAt)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pushSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pushSubscriptionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subscription
	err := r.client.DeleteSubscription(int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Strava Subscription",
			"Could not delete subscription, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new subscription
	pushSubscription, err := r.client.CreateSubscription(strava.SubscriptionItem{
		CallbackURL: plan.CallbackURL.ValueString(),
		VerifyToken: plan.VerifyToken.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subscription",
			"Could not create subscription, unexpected error: "+err.Error(),
		)
		return
	}

	pushSubscription, err = r.client.GetSubscription(int(pushSubscription.ID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created subscription",
			"Could not read created subscription, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.Int64Value(int64(pushSubscription.ID))
	plan.ResourceState = types.Int64Value(int64(pushSubscription.ResourceState))
	plan.ApplicationID = types.Int64Value(int64(pushSubscription.ApplicationID))
	plan.CallbackURL = types.StringValue(pushSubscription.CallbackURL)
	plan.CreatedAt = types.StringValue(pushSubscription.CreatedAt)
	plan.UpdatedAt = types.StringValue(pushSubscription.UpdatedAt)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pushSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pushSubscriptionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subscription
	err := r.client.DeleteSubscription(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Strava Subscription",
			"Could not delete subscription, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *pushSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*strava.Client)
}

func (r *pushSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")

	if len(parts) < 2 {
		resp.Diagnostics.AddError(
			"Error importing item",
			"Could not import item, unexpected error (ID should be in the format <id>,<verify_token>): "+req.ID,
		)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	verifyToken := parts[1]

	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing item",
			"Could not import item, unexpected error (The <id> part should be an integer): "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("verify_token"), verifyToken)...)
}
