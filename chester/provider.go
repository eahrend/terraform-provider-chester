package chester

import (
	"context"
	"fmt"

	chester "github.com/eahrend/terraform-provider-chester/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SqlAdminSvcChester, is a struct that wraps the chesterClient and
// any other services this needs. Previously used to contain a sqladmin
// client.
// TODO: Remove this and just use the API client
type SqlAdminSvcChester struct {
	ChesterClient *chester.Client
}

// Provider is the main terraform provider for chester.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHESTER_HOST", ""),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHESTER_USERNAME", ""),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Sensitive:   true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHESTER_PASSWORD", ""),
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHESTER_CLIENT_ID", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"chester_database": resourceDatabase(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"chester_database": dataSourceDatabase(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)
	audience := d.Get("client_id").(string)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	if (username != "") && (password != "") {
		c, err := chester.NewClient(host, username, password, audience)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to create Chester client %s", err.Error()),
			})
			return nil, diags
		}
		return c, diags
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Set user name and password",
			Detail:   "Set it dumb dumb",
		})
		return nil, diags
	}
}
