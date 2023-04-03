# Manage example push subscription.
resource "strava_push_subscription" "example" {
  callback_url = "http://a-valid.com/url"
  verify_token = "STRAVA"
}
