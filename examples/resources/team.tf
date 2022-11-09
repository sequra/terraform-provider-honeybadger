# Create a new team
resource "honeybadger_team" "new_team" { # terraform import honeybadger_team.new_team 1234
  name = "Terraform team 2"
}
