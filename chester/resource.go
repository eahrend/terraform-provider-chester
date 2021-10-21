package chester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	models "github.com/eahrend/chestermodels"
	chester "github.com/eahrend/terraform-provider-chester/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDatabaseRead,
		DeleteContext: resourceDatabaseDelete,
		CreateContext: resourceDatabaseCreate,
		UpdateContext: resourceDatabaseUpdate,
		Schema: map[string]*schema.Schema{
			"instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_ssl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"read_hostgroup": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"write_hostgroup": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			// not used, but I figured if I publish this, it wouldn't hurt to have
			"cert_data": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			// not used, but I figured if I publish this, it wouldn't hurt to have
			"key_data": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			// not used, but I figured if I publish this, it wouldn't hurt to have
			"ca_data": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"max_chester_instances": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"query_rules": &schema.Schema{
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: supressQueryRules,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"active": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"match_digest": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"destination_hostgroup": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"apply": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"comment": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// Nothing needs to be added here, map[string]interface allows us to add/remove as needed
			"master_instance": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},
			// TODO: once proxysql adds instance:ssl conifg we'll implement it here
			// 	need to make this a required variable, which may require some modifications
			// 	on zeus
			"read_replicas": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	mrb, err := json.Marshal(&db.MasterInstance)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed marshalling master instance %s", err.Error()),
		})
		return diags
	}
	err = json.NewDecoder(bytes.NewBuffer(mrb)).Decode(&mr)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed decoding master instance %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("master_instance", mr); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting master_instance with error %s", err.Error()),
		})
		return diags
	}

	if err := d.Set("max_chester_instances", db.ChesterMetaData.MaxChesterInstances); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting metadata with error %s", err.Error()),
		})
		return diags
	}
	rrs := []map[string]interface{}{}
	rrsb, err := json.Marshal(&db.ReadReplicas)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed mashalling read replicas %s", err.Error()),
		})
		return diags
	}
	err = json.NewDecoder(bytes.NewBuffer(rrsb)).Decode(&rrs)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed decoding read replicas %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("read_replicas", rrs); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting read_replicas with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("read_hostgroup", db.ReadHostGroup); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting read_replicas with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("write_hostgroup", db.WriteHostGroup); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting read_replicas with error %s", err.Error()),
		})
		return diags
	}
	if err := d.Set("enable_ssl", db.UseSSL); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed setting read_replicas with error %s", err.Error()),
		})
		return diags
	}
	d.SetId(d.Get("instance_name").(string))
	return diags
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*chester.Client)
	var diags []diag.Diagnostic
	diags = append(diags, diag.Diagnostic{
		Summary:  "Starting Resource Create",
		Severity: diag.Warning,
	})
	mi := d.Get("master_instance").(map[string]interface{})
	rris := d.Get("read_replicas").([]interface{})
	rrs := []models.AddDatabaseRequestDatabaseInformation{}
	for _, rr := range rris {
		newrr := rr.(map[string]interface{})
		b, err := json.Marshal(newrr)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("unable to modify query rule: %s", err.Error()),
			})
			return diags
		}
		read_replica := models.AddDatabaseRequestDatabaseInformation{}
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&read_replica)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("unable to modify query rule: %s", err.Error()),
			})
			return diags
		}
		rrs = append(rrs, read_replica)
	}
	cmd := models.ChesterMetaData{
		InstanceGroup:       d.Get("instance_name").(string),
		MaxChesterInstances: d.Get("max_chester_instances").(int),
	}
	// removing the cert/sa/key stuff here, and will re-add it once it becomes a feature of proxysql
	db := models.AddDatabaseRequest{
		EnableSSL:    0,
		Action:       "add",
		InstanceName: d.Get("instance_name").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		MasterInstance: models.AddDatabaseRequestDatabaseInformation{
			Name:      mi["name"].(string),
			IPAddress: mi["ip_address"].(string),
		},
		ReadReplicas:    rrs,
		ChesterMetaData: cmd,
	}
	_, err := c.AddDatabase(db)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to add database to datastore %s", err.Error()),
		})
		return diags
	}

	d.SetId(d.Get("instance_name").(string))
	resourceDatabaseRead(ctx, d, m)
	return diags
}

// TODO: Rework this to make one call, and add user/pass changes as needed
func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*chester.Client)
	instanceName := d.Get("instance_name").(string)
	var diags []diag.Diagnostic
	mdbr := models.ModifyDatabaseRequest{
		Action:       "modify",
		InstanceName: instanceName,
	}
	callChange := false
	if d.HasChange("username") {
		_, newUser := d.GetChange("username")
		mdbr.NewUsername = newUser.(string)
		callChange = true
	}
	if d.HasChange("password") {
		_, newPassword := d.GetChange("password")
		mdbr.NewPassword = newPassword.(string)
		callChange = true
	}
	if d.HasChange("query_rules") {
		qrs := d.Get("query_rules").([]interface{})
		for _, queryRule := range qrs {
			// check if username has changed
			queryRuleStruct := queryRule.(models.ProxySqlMySqlQueryRule)
			if d.HasChange("username") {
				_, newUser := d.GetChange("username")
				queryRuleStruct.Username = newUser.(string)
			}
			err := c.ModifyQueryRuleByID(queryRuleStruct)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("unable to modify query rule: %s", err.Error()),
				})
				return diags
			}
		}
	}
	if d.HasChange("read_replicas") {
		callChange = true
		rris := d.Get("read_replicas").([]interface{})
		rrs := []models.AddDatabaseRequestDatabaseInformation{}
		for _, rr := range rris {
			newrr := rr.(map[string]interface{})
			b, err := json.Marshal(newrr)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("unable to modify query rule: %s", err.Error()),
				})
				return diags
			}
			read_replica := models.AddDatabaseRequestDatabaseInformation{}
			err = json.NewDecoder(bytes.NewReader(b)).Decode(&read_replica)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("unable to modify query rule: %s", err.Error()),
				})
				return diags
			}
			rrs = append(rrs, read_replica)
		}
		mdbr.ReadReplicas = rrs
	}
	if d.HasChange("max_chester_instances") {
		callChange = true
		mdbr.ChesterMetaData = models.ChesterMetaData{
			InstanceGroup:       d.Get("instance_name").(string),
			MaxChesterInstances: d.Get("max_chester_instances").(int),
		}
	}
	if callChange {
		err := c.ModifyDatabase(mdbr)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
	}
	debugdiags := resourceDatabaseRead(ctx, d, m)
	diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: "Running read after update"})
	for _, dd := range debugdiags {
		diags = append(diags, dd)
	}
	return diags
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*chester.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Summary:  "Starting Resource Delete",
		Severity: diag.Warning,
	})
	instanceName := d.Get("instance_name").(string)
	userName := d.Get("username").(string)
	rdr := models.RemoveDatabaseRequest{
		Action:       "remove",
		InstanceName: instanceName,
		Username:     userName,
	}
	err := c.RemoveDatabase(rdr)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
