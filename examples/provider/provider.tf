terraform {
  required_providers {
    honeybadger = {
      version = "~> 0.1.0"
      source  = "sequra.com/providers/honeybadger"
    }
  }
}
provider "honeybadger" {
  api_key = "<INTRODUCE_YOUR_API_KEY>"
}
