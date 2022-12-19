package honeybadger

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	hbc "terraform-provider-honeybadger/cli"
)

func TestAccHoneybadgerProjectBasic(t *testing.T) {
	projectName := "New Project"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHoneybadgerProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHoneybadgerProjectConfigBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHoneybadgerProjectExists("honeybadger_project.test"),
				),
			},
		},
	})
}
func testAccCheckHoneybadgerProjectConfigBasic(projectName string) string {
	return fmt.Sprintf(`
	resource "honeybadger_project" "test" {
		name = %s
	}
	`, projectName)
}

func testAccCheckHoneybadgerProjectDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*hbc.HoneybadgerClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "honeybadger_project" {
			continue
		}

		projectID := rs.Primary.ID

		projectIDToString, err := strconv.Atoi(projectID)
		if err != nil {
			return err
		}
		err = c.DeleteTeam(projectIDToString)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckHoneybadgerProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No projectID set")
		}

		return nil
	}
}
