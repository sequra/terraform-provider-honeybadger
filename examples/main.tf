terraform {
  required_providers {
    honeybadger = {
      version = "~> 0.1.0"
      source  = "sequra.com/providers/honeybadger"
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
  admin = false
  team_id = [honeybadger_team.new_team.id]
}


# Update user - before you need to import using: terraform import honeybadger_user.TestSequra2 test.sequra.page2@sequra.es
resource "honeybadger_user" "TestSequra2" {
  email = "test.sequra.page2@sequra.es"
  admin = true
  team_id = [honeybadger_team.new_team.id]
}

