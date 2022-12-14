
---
layout: ""
page_title: "Honeybadger: honeybadger_project"
description: |-
  Creates and manages projects within your Honeybadger organization
---

# honeybadger_team (Resource)

This resource allows you to create and manage projects within your Honeybadger organization.


## Example Usage

```terraform
# Create a new project
resource "honeybadger_project" "new_project" { # terraform import honeybadger_project.new_project 1234
  name = "Terraform project"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `language` (String)
- `last_updated` (String)

### Read-Only

- `id` (String) The ID of this resource.


# Import

Repositories can be imported using the project id, e.g.

```
$ terraform import honeybadger_project.new_team 1234
```
