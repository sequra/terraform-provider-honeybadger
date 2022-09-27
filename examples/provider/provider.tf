terraform {
  required_providers {
    honeybadger = {
      version = "~> 1.0.0"
      source = "sequra/honeybadger"
    }
  }
}
provider "honeybadger" {
  api_key = "<INTRODUCE_YOUR_API_KEY>"
}
