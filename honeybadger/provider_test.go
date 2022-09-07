package honeybadger

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"honeybadger": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("HONEYBADGER_HOST"); err == "" {
		t.Fatal("HONEYBADGER_HOST must be set for acceptance tests")
	}
	if err := os.Getenv("HONEYBADGER_API_KEY"); err == "" {
		t.Fatal("HONEYBADGER_API_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("HONEYBADGER_TEAM_ID"); err == "" {
		t.Fatal("HONEYBADGER_TEAM_ID must be set for acceptance tests")
	}
}
