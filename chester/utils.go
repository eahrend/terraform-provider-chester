package chester

import (
	models "github.com/eahrend/chestermodels"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func supressQueryRules(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("query_rules") {
		return false
	} else {
		return true
	}
}

func flattenQueryRules(queryRules []models.ProxySqlMySqlQueryRule) []interface{} {
	if queryRules != nil {
		qrs := make([]interface{}, len(queryRules), len(queryRules))
		for i, queryRule := range queryRules {
			qr := make(map[string]interface{})
			qr["rule_id"] = queryRule.RuleID
			qr["username"] = queryRule.Username
			qr["active"] = queryRule.Active
			qr["match_digest"] = queryRule.MatchDigest
			qr["destination_hostgroup"] = queryRule.DestinationHostgroup
			qr["apply"] = queryRule.Apply
			qr["comment"] = queryRule.Comment
			qrs[i] = qr
		}
		return qrs
	}
	return make([]interface{}, 0)
}
