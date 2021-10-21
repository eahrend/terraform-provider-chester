package chester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	chester "github.com/eahrend/terraform-provider-chester/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO: Once proxysql enables instance:ssl config, add ssl config to each hostgroup
func dataSourceDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcesDatabaseRead,
		Schema: map[string]*schema.Schema{
			"instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"read_hostgroup": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"write_hostgroup": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_ssl": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query_rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"match_digest": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_hostgroup": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"apply": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"comment": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"master_instance": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"read_replicas": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcesDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*chester.Client)
	// Warning or errors can be collected in a slice type
	databaseName := d.Get("instance_name").(string)
	var diags diag.Diagnostics
	// Warning or errors can be collected in a slice type
	db, err := c.GetDatabase(databaseName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed getting database with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("username", db.Username); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting username with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("password", db.Password); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting password with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("instance_name", db.InstanceName); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting instance_name with error %s", err.Error()),
		})
		return diags
	}
	queryRules := flattenQueryRules(db.QueryRules)
	if err := d.Set("query_rules", queryRules); err != nil {
		return diag.FromErr(err)
	}
	mr := map[string]interface{}{}
	mrb, _ := json.Marshal(&db.MasterInstance)
	json.NewDecoder(bytes.NewBuffer(mrb)).Decode(&mr)
	if err := d.Set("master_instance", mr); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting master_instance with error %s", err.Error()),
		})
		return diags
	}
	rrs := []map[string]interface{}{}
	rrsb, _ := json.Marshal(&db.ReadReplicas)
	json.NewDecoder(bytes.NewBuffer(rrsb)).Decode(&rrs)
	if err := d.Set("read_replicas", rrs); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting read_replicas with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("username", db.Username); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting username with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("password", db.Password); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting password with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("read_hostgroup", db.ReadHostGroup); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting read_hostgroup with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("write_hostgroup", db.WriteHostGroup); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting write_hostgroup with error %s", err.Error()),
		})
		return diags
	}
	// removing cert/ca/key fields here, will readd them once they become relevant again
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
