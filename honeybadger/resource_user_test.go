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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHoneybadgerUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHoneybadgerUserConfigBasic(email, isAdmin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHoneybadgerUserExists("honeybadger_user.new"),
				),
			},
		},
	})
}
func testAccCheckHoneybadgerUserConfigBasic(email string, isAdmin bool) string {
	return fmt.Sprintf(`
	resource "honeybadger_user" "new" {
		email = %s
		admin = %t
	}
	`, email, isAdmin)
}

func testAccCheckHoneybadgerUserDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*hbc.HoneybadgerClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "honeybadger_user" {
			continue
		}

		userID := rs.Primary.ID

		useridToString, err := strconv.Atoi(userID)
		if err != nil {
			return err
		}
		err = c.DeleteUser(useridToString)
		if err != nil {
			return err
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
