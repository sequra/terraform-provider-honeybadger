terraform {
  required_providers {
    honeybadger = {
      source = "sequra/honeybadger"
      version = "~> 1.0.0"
    }
  }
}
provider "honeybadger" {
  api_key = "sdfsd89342"
}


# Return all teams & users
data "honeybadger_teams" "all" {}
output "all_teams" {
  value = data.honeybadger_teams.all.teams
}

# Create a new team
resource "honeybadger_team" "new_team" { # terraform import honeybadger_team.new_team 1234
  name = "Terraform team 2"
}


# Create a new user
resource "honeybadger_user" "new" {
  email = "test.sequra@sequra.es"
  team {
    id = honeybadger_team.new_team.id
    is_admin = false
  }
}


# Update user - before you need to import using: terraform import honeybadger_user.TestSequra2 test.sequra.page2@sequra.es
resource "honeybadger_user" "TestSequra2" {
  email = "test.sequra.page2@sequra.es"
  team {
    id = honeybadger_team.new_team.id
    is_admin = false
  }
}

