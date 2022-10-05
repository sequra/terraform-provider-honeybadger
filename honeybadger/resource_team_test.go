package honeybadger

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	hbc "terraform-provider-honeybadger/cli"
)

func TestAccHoneybadgerTeamBasic(t *testing.T) {
	teamName := "Test Team"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHoneybadgerTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHoneybadgerTeamConfigBasic(teamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHoneybadgerTeamExists("honeybadger_team.test"),
				),
			},
		},
	})
}
func testAccCheckHoneybadgerTeamConfigBasic(teamName string) string {
	return fmt.Sprintf(`
	resource "honeybadger_team" "test" {
		name = %s
	}
	`, teamName)
}

func testAccCheckHoneybadgerTeamDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*hbc.HoneybadgerClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "honeybadger_team" {
			continue
		}

		userID := rs.Primary.ID

		teamIDToString, err := strconv.Atoi(userID)
		if err != nil {
			return err
		}
		err = c.DeleteTeam(teamIDToString)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckHoneybadgerTeamExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No UserID set")
		}

		return nil
	}
}
