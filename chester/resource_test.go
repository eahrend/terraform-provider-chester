package chester

import (
	"errors"
	"fmt"
	"os"
	"testing"

	models "github.com/eahrend/chestermodels"
	"github.com/eahrend/terraform-provider-chester/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateDB(t *testing.T) {
	db := models.InstanceData{}
	instanceName := os.Getenv("INSTANCE_NAME")
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	readReplicaOneName := os.Getenv("RR_ONE_NAME")
	readReplicaTwoName := os.Getenv("RR_TWO_NAME")
	readReplicaOneIP := os.Getenv("RR_ONE_IP")
	readReplicaTwoIP := os.Getenv("RR_TWO_IP")
	writerInstanceName := os.Getenv("WR_NAME")
	writerInstanceIP := os.Getenv("WR_IP")
	sqlProjectID := os.Getenv("SQL_PROJECT_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config: createDBConfig(instanceName, userName, password, readReplicaOneName, readReplicaTwoName, readReplicaOneIP,
					readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccDBExists("chester_database.chester_proxysql", &db),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config: modifyMaxCount(instanceName, userName, password, readReplicaOneName, readReplicaTwoName, readReplicaOneIP,
					readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccDBMaxInstanceCountModified("chester_database.chester_proxysql", &db),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config: addReadReplica(instanceName, userName, password, readReplicaOneName, readReplicaTwoName, readReplicaOneIP,
					readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccAddReadReplica("chester_database.chester_proxysql", &db),
				),
			},
		},
	})
}

// createDBConfig is the creation of the database
// I'll have a way to configure this via env vars from the output
// of the terraform templates in ../terraform, but for now I can use static
// configuration.
func createDBConfig(instanceName, userName, password, readReplicaOneName, readReplicaTwoName, readReplicaOneIP,
	readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID string) string {
	return fmt.Sprintf(`
		resource "chester_database" "chester_proxysql" {
			instance_name = "%s"
			username = "%s"
			password = "%s"
		  	max_chester_instances = 4
			read_replicas {
				name = "%s"
				ip_address = "%s"
			}
			read_replicas {
				name = "%s"
				ip_address = "%s" 
			}
		  	enable_ssl = 0
		  	read_hostgroup = 10
		  	write_hostgroup = 5
		  	master_instance = {
				name = "%s"
				ip_address = "%s"
		  	}
		  	sql_project_id = "%s"
		}
	`, instanceName, userName, password, readReplicaOneName, readReplicaOneIP, readReplicaTwoName, readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID)
}

func addReadReplica(instanceName, userName, password, readReplicaOneName, readReplicaTwoName, readReplicaOneIP,
	readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID string) string {
	return fmt.Sprintf(`
		resource "chester_database" "chester_proxysql" {
			instance_name = "%s"
			username = "%s"
			password = "%s"
		  	max_chester_instances = 4
			read_replicas {
				name = "%s"
				ip_address = "%s"
			}
			read_replicas {
				name = "%s"
				ip_address = "%s" 
			}
			read_replicas {
				name = "chester-test-replica-3"
				ip_address = "10.149.0.8" 
			}
		  	enable_ssl = 0
		  	read_hostgroup = 10
		  	write_hostgroup = 5
		  	master_instance = {
				name = "%s"
				ip_address = "%s"
		  	}
		  	sql_project_id = "%s"
		}
	`, instanceName, userName, password, readReplicaOneName, readReplicaOneIP, readReplicaTwoName, readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID)
}

func modifyMaxCount(instanceName, userName, password, readReplicaOneName, readReplicaTwoName, readReplicaOneIP,
	readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID string) string {
	return fmt.Sprintf(`
		resource "chester_database" "chester_proxysql" {
			instance_name = "%s"
			username = "%s"
			password = "%s"
		  	max_chester_instances = 3
			read_replicas {
				name = "%s"
				ip_address = "%s"
			}
			read_replicas {
				name = "%s"
				ip_address = "%s" 
			}
		  	enable_ssl = 0
		  	read_hostgroup = 10
		  	write_hostgroup = 5
		  	master_instance = {
				name = "%s"
				ip_address = "%s"
		  	}
		  	sql_project_id = "%s"
		}
	`, instanceName, userName, password, readReplicaOneName, readReplicaOneIP, readReplicaTwoName, readReplicaTwoIP, writerInstanceName, writerInstanceIP, sqlProjectID)
}

func testAccDBMaxInstanceCountModified(resourceName string, db *models.InstanceData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return errors.New("failed to get resource")
		}
		client := testAccProvider.Meta().(*api.Client)
		resp, err := client.GetDatabase(val.Primary.ID)
		if err != nil {
			return err
		}
		if resp.ChesterMetaData.MaxChesterInstances != 3 {
			return errors.New("max chester instance mismatch")
		}
		return nil
	}
}

func testAccAddReadReplica(resourceName string, db *models.InstanceData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return errors.New("failed to get resource")
		}
		client := testAccProvider.Meta().(*api.Client)
		resp, err := client.GetDatabase(val.Primary.ID)
		if err != nil {
			return err
		}
		if len(resp.ReadReplicas) != 3 {
			return errors.New("failed to add instance")
		}
		return nil
	}
}

func testAccDBExists(resourceName string, db *models.InstanceData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return errors.New("failed to get resource")
		}
		if val.Primary.ID == "" {
			return errors.New("failed to get id")
		}
		client := testAccProvider.Meta().(*api.Client)
		resp, err := client.GetDatabase(val.Primary.ID)
		if err != nil {
			return err
		}
		if resp.InstanceName != val.Primary.ID {
			return errors.New("id mismatch")
		}
		return nil
	}
}
