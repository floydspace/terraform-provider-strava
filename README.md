# terraform-provider-strava

This is a terraform provider for managing your strava.

## Example

```tf
provider "strava" {
  client_id     = "5"
  client_secret = "7b2946535949ae70f015d696d8ac602830ece412"
}

resource "strava_push_subscription" "example" {
  callback_url = "http://a-valid.com/url"
  verify_token = "STRAVA"
}

output "subscription_id" {
  value = strava_push_subscription.example.id
}
```

## Installation

Add the following to your terraform configuration

```tf
terraform {
  required_providers {
    strava = {
      source  = "floydspace/strava"
      version = "~> 0.1.0"
    }
  }
}
```
