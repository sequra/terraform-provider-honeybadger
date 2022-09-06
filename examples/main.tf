terraform {
  required_providers {
    honeybadger = {
      version = "~> 0.1.0"
      source  = "sequra.com/providers/honeybadger"
    }
  }
}
provider "honeybadger" {
  host = "http://localhost:8080"
  team_id = 1234
  api_key = "sdfsd89342"
}


# Returns all users

data "honeybadger_users" "all" {}
output "all_users" {
  value = data.honeybadger_users.all.users
}

# Create a new user
resource "honeybadger_user" "new" {
  email = "test.sequra@sequra.es"
  admin = false
}


# Update user - before you need to import using: terraform import honeybadger_user.TestSequra2 test.sequra.page2@sequra.es
resource "honeybadger_user" "TestSequra2" {
  email = "test.sequra.page2@sequra.es"
  admin = true
}

