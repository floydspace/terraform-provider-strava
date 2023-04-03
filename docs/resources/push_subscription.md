---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "strava_push_subscription Resource - strava"
subcategory: ""
description: |-
  Manages a Strava push subscription.
---

# strava_push_subscription (Resource)

Manages a Strava push subscription.

## Example Usage

```terraform
# Manage example push subscription.
resource "strava_push_subscription" "example" {
  callback_url = "http://a-valid.com/url"
  verify_token = "STRAVA"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `callback_url` (String) Address where webhook events will be sent; maximum length of 255 characters.
- `verify_token` (String, Sensitive) String chosen by the application owner for client security. An identical string will be included in the validation request made by Strava's subscription service.

### Read-Only

- `application_id` (Number) Strava API application ID.
- `created_at` (String) Date and time the subscription was created.
- `id` (Number) Push subscription ID.
- `last_updated` (String) Timestamp of the last Terraform update to the push subscription.
- `resource_state` (Number) State of the push subscription.
- `updated_at` (String) Date and time the subscription was last updated.

## Import

Import is supported using the following syntax:

```shell
# Push subscription can be imported by specifying the subscription identifier and the verify token.
terraform import strava_push_subscription.example 12345,STRAVA
```