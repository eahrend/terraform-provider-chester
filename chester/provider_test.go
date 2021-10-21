package chester

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"terraform-provider-chester": func() (*schema.Provider, error) { return testAccProvider, nil },
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CHESTER_HOST"); v == "" {
		t.Fatal("CHESTER_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("CHESTER_USERNAME"); v == "" {
		t.Fatal("CHESTER_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("CHESTER_PASSWORD"); v == "" {
		t.Fatal("CHESTER_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("CHESTER_CLIENT_ID"); v == "" {
		t.Fatal("CHESTER_CLIENT_ID must be set for acceptance tests")
	}
}
