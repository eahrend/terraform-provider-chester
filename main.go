package main

import (
	"github.com/eahrend/terraform-provider-chester/chester"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// main is the entry point to the terraform provider
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return chester.Provider()
		},
	})
}
