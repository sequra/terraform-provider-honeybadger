# Create a new project
resource "honeybadger_project" "new_project" { # terraform import honeybadger_project.new_project 1234
  name = "Terraform project"
}
