# Create a new user
resource "honeybadger_user" "new" {
  email = "test.sequra@sequra.es"
  team {
    id = honeybadger_team.new_team.id
    is_admin = false
  }
}
