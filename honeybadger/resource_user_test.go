package honeybadger

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	hbc "terraform-provider-honeybadger/cli"
)

func TestAccHoneybadgerUserBasic(t *testing.T) {
	email := "test.sequra@sequra.es"
	isAdmin := true
	teamID := 1234

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHoneybadgerUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHoneybadgerUserConfigBasic(email, isAdmin, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHoneybadgerUserExists("honeybadger_user.test"),
				),
			},
		},
	})
}
func testAccCheckHoneybadgerUserConfigBasic(email string, isAdmin bool, teamID int) string {
	return fmt.Sprintf(`
	resource "honeybadger_user" "test" {
		email = %s
		admin = %t
 		team_id = [%d]
	}
	`, email, isAdmin, teamID)
}

func testAccCheckHoneybadgerUserDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*hbc.HoneybadgerClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "honeybadger_user" {
			continue
		}

		userID := rs.Primary.ID
		teams := rs.Primary.Attributes["team_id"]

		useridToString, err := strconv.Atoi(userID)
		if err != nil {
			return err
		}
		for _, team := range teams {
			err = c.DeleteUser(useridToString, int(team))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func testAccCheckHoneybadgerUserExists(n string) resource.TestCheckFunc {
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
