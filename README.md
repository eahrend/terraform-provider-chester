# Terraform-Provider-Chester

Terraform interface for the proxysql autoscaling solution chester.



## Resource Inputs
| Input Name      	| Type                                                                                                              	| Required 	| Default 	| Sensitive 	| Description                                                                                                                                                                     	|
|-----------------	|-------------------------------------------------------------------------------------------------------------------	|----------	|---------	|-----------	|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|
| instance_name   	| string                                                                                                            	| true     	| N/A     	| false     	| Cloud SQL instance name                                                                                                                                                         	|
| sql_project_id  	| string                                                                                                            	| true     	| N/A     	| false     	| Project ID of the SQL instance this will belong to                                                                                                                              	|
| enable_ssl      	| boolean                                                                                                           	| true     	| N/A     	| false     	| NOTE: This is largely ignored because of the lack of ProxySQL cert mappings                                                                                                     	|
| username        	| string                                                                                                            	| true     	| N/A     	| false     	| Cloud SQL instance username                                                                                                                                                     	|
| password        	| string                                                                                                            	| true     	| N/A     	| true      	| Cloud SQL instance password                                                                                                                                                     	|
| read_hostgroup  	| int                                                                                                               	| true     	| N/A     	| false     	| Hostgroup number for the read replicas on the proxysql instance                                                                                                                 	|
| write_hostgroup 	| int                                                                                                               	| true     	| N/A     	| false     	| Hostgroup number for the write replica on the proxysql instance                                                                                                                 	|
| query_rules     	| list(obj({<br>username: string,<br>active: int,<br>match_digest: string,<br>destination_hostgroup: int,<br><br>}) 	| false    	| N/A     	| false     	| Query rules, if not specified it uses the default based on your read/write hostgroups. Details can be found: https://proxysql.com/documentation/main-runtime/#mysql_query_rules 	|
| master_instance 	| obj({<br>name: string,<br>ip_address: string,<br>})                                                               	| true     	| N/A     	| false     	| Details about the master instance                                                                                                                                               	|
| read_replicas   	| list(obj({<br>name: string,<br>ip_address: string,<br>})                                                          	| true     	| N/A     	| false     	| Details about the read replicas                                                                                                                                                 	|                                                    	|



## Provider Input
| Input    	| Type   	| Required 	| Default 	| Sensitive 	| Description                                                                                                                                            	|
|----------	|--------	|----------	|---------	|-----------	|--------------------------------------------------------------------------------------------------------------------------------------------------------	|
| host     	| string 	| true     	| N/A     	| false     	| Url of the chester-api instance, example: http://0.0.0.0                                                                                               	|
| username 	| string 	| true     	| N/A     	| false     	| Username of the chester-api instance NOTE: This will be deprecated in favor of IAP authentication	|
| password 	| string 	| true     	| N/A     	| true      	| Password of the chester-api instance NOTE: This will be deprecated in favor of IAP authentication	|

## Example Usage
```hcl-terraform
resource "chester_database" "chester_proxysql" {
  instance_name = "database-name"
  username = "sqlUserName"
  password = "notasecurepassword"
  read_hostgroup = "10"
  write_hostgroup = "5"
  master_instance = {
    name = "database-name-master"
    ip_address = "1.2.3.4"
  }
  sql_project_id = "project-name"
}

terraform {
  required_providers {
    chester = {
      source = "eahrend.com/eahrend/chester"
      version = "0.0.1-alpha"
    }
  }
}

provider "chester" {
  host = "http://0.0.0.0:80"
  username = "chester"
  password = "bison"
}
```

## Installation
Download your OS/Arch from here:
https://github.com/eahrend/terraform-chester-provider/releases

Unzip to your terraform provider folder.

On mac OS it will be located at `~/.terraform.d/plugins/eahrend.com/eahrend/chester/<version number>/<os>_<arch>`

## Testing
1. Run terraform apply in ./terraform, after the resources are created you'll need to manually create the IAP resource via the security panel in the console.
2. After the IAP resource is created, you shouldn't need to re-configure anything else. 
3. Get the configs required in ./chester and set the env vars then run go run test -v 
4. You'll need to set the application default credentials to a service account, since it simplifies programmatic access. 


## Notes
1. While the setup in `./terraform` is a good basis for a network/gke setup, it should not be considered "production ready".
2. Please, please, please, please, use RBAC in your GKE clusters and a service mesh to limit communication between services and encrypt your traffic. I don't here because this is just for testing purposes.