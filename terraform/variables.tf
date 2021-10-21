## dummy email variable
variable "allowed_user_emails" {
  type = list(string)
  description = "email for allowed users"
}

## DNS for chester HTTP server
variable "chester_dns_name" {
  type = string
  description = "dns entry for chester"
}

## Required variables for google project factory
variable "google_project_factory_name" {
  type = string
  description = "name of the google project"
}


variable "google_project_factory_org_id" {
  type = string
  description = "The org ID where the project will reside"
}


variable "google_project_factory_billing_account" {
  type = string
  description = "billing account ID"
}

## Required variables for docker pull secret

variable "docker_pull_secret_string" {
  type = string
  description = "base64 encoded string for the docker pull secret"
}


## DNS related stuff
variable "dns_managed_zone_name" {
  type = string
  description = "name of the managed zone we want to get the details of"
}

variable "dns_managed_zone_project" {
  type = string
  description = "project of the dns managed zone you want to use"
}


### DEFAULT VARIABLES START HERE

## Default variables for google project factory
variable "google_project_factory_random_project_id" {
  type = bool
  description = "whether or not we give the project a random id"
  default = true
}

variable "google_project_factory_auto_create_network" {
  type = bool
  description = "whether or not to auto-create the network"
  default = true
}

variable "google_project_factory_activate_apis" {
  type = list(string)
  description = "list of apis to activate inside the project"
  default = [
    "compute.googleapis.com",
    "iap.googleapis.com",
    "cloudbuild.googleapis.com",
    "cloudasset.googleapis.com",
    "container.googleapis.com",
    "cloudkms.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "serviceusage.googleapis.com",
    "storage-api.googleapis.com",
    "redis.googleapis.com",
    "servicenetworking.googleapis.com",
    "datastore.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudkms.googleapis.com"
  ]
}

variable "inet_ranges" {
  type = list(string)
  description = "list of IPV4 CIDR ranges to nat"
  default = [
    "0.0.0.0/5",
    "8.0.0.0/7",
    "11.0.0.0/8",
    "12.0.0.0/6",
    "16.0.0.0/4",
    "32.0.0.0/3",
    "64.0.0.0/2",
    "128.0.0.0/3",
    "160.0.0.0/5",
    "168.0.0.0/6",
    "172.0.0.0/12",
    "172.32.0.0/11",
    "172.64.0.0/10",
    "172.128.0.0/9",
    "173.0.0.0/8",
    "174.0.0.0/7",
    "176.0.0.0/4",
    "192.0.0.0/9",
    "192.128.0.0/11",
    "192.160.0.0/13",
    "192.169.0.0/16",
    "192.170.0.0/15",
    "192.172.0.0/14",
    "192.176.0.0/12",
    "192.192.0.0/10",
    "193.0.0.0/8",
    "194.0.0.0/7",
    "196.0.0.0/6",
    "200.0.0.0/5",
    "208.0.0.0/4"]
}

## Network default variables

variable "network_auto_create_subnetworks" {
  type = bool
  description = "auto-create subnetworks or not"
  default = false
}

variable "network_description" {
  type = string
  description = "description of the description lmoa"
  default = "network description"
}

variable "network_delete_default_internet_gateway_routes" {
  type        = bool
  description = "If set, ensure that all routes within the network specified whose names begin with 'default-route' and with a next hop of 'default-internet-gateway' are deleted"
  default     = false
}

variable "network_routing_mode" {
  type = string
  description = "routing mode"
  default = "GLOBAL"
}

variable "shared_vpc_host" {
  type = bool
  description = "is this a host vpc"
  default = false
}

variable "vpc_name" {
  type = string
  description = "name of the vpc"
  default = "chester-test-vpc"
}

variable "vpc_subnets" {
  type        = list(map(string))
  default = [
    {
      subnet_name: "us-east1-test",
      subnet_region: "us-east1",
      cidr: "10.64.0.0/16",
      subnet_private_access: true,
      subnet_flow_logs: true,
      subnet_flow_logs_interval: "INTERVAL_15_MIN",
      subnet_flow_logs_sampling: 0.5,
      subnet_flow_logs_metadata: "INCLUDE_ALL_METADATA",
      secondary_ranges: "[]"
    },
    {
      subnet_name: "us-east1-test-db",
      subnet_region: "us-east1",
      cidr: "10.65.0.0/16",
      subnet_private_access: true,
      subnet_flow_logs: true,
      subnet_flow_logs_interval: "INTERVAL_15_MIN",
      subnet_flow_logs_sampling: 0.5,
      subnet_flow_logs_metadata: "INCLUDE_ALL_METADATA",
      secondary_ranges: "[]"
    },
    {
      subnet_name: "us-east1-test-compute",
      subnet_region: "us-east1",
      cidr: "10.66.0.0/17",
      subnet_private_access: true,
      subnet_flow_logs: true,
      subnet_flow_logs_interval: "INTERVAL_15_MIN",
      subnet_flow_logs_sampling: 0.5,
      subnet_flow_logs_metadata: "INCLUDE_ALL_METADATA",
      secondary_ranges: "[{range_name : us-east1-test-pods, ip_cidr_range : 10.66.128.0/18},{range_name : us-east1-test-services, ip_cidr_range : 10.66.192.0/18 }]"
    }]
}


## GKE Default variables
variable "gke_kubernetes_version" {
  type = string
  description = "kubernetes version on masters"
  default = "latest"
}
variable "gke_premptible_nodes" {
  type = bool
  description = "whether to make the nodes preemptible"
  default = true
}
variable "gke_auto_upgrade" {
  type = bool
  description = "whether to auto-upgrade"
  default = false
}

variable "gke_auto_repair" {
  type = bool
  description = "whether to auto repair nodes"
  default = true
}

variable "gke_image_type" {
  type = string
  description = "image type we want on our nodes"
  default = "cos_containerd"
}

variable "gke_disk_type" {
  type = string
  description = "disk type of the node"
  default = "pd-standard"
}

variable "gke_disk_size_gb" {
  type = number
  description = "disk size on the nodes"
  default = 100
}

variable "gke_local_ssd_count" {
  type = number
  description = "number of ssds on the node"
  default = 0
}

variable "gke_node_max_count" {
  type = number
  description = "number of max nodes we want"
  default = 6
}

variable "gke_node_min_count" {
  type = number
  description = "number of minimal nodes we want"
  default = 3
}

variable "gke_node_machine_type" {
  type = string
  description = "machine type of the nodes"
  default = "n1-standard-4"
}

variable "gke_node_name" {
  type = string
  description = "name of the node pool"
  default = "node-pool"
}

variable "gke_cluster_name" {
  type = string
  description = "name of the cluster to create"
  default = "chester-test-cluster"
}

variable "gke_region" {
  type = string
  description = "The region to host the cluster in"
  default     = "us-east1"
}

variable "gke_zones" {
  type = list(string)
  description = "list of zones to deploy the nodes to"
  default = ["us-east1-b","us-east1-d", "us-east1-c"]
}

variable "gke_ip_range_pods_name" {
  type = string
  description = "The secondary ip range to use for pods"
  default     = "us-east1-test-pods"
}

variable "gke_ip_range_services_name" {
  type = string
  description = "The secondary ip range to use for services"
  default     = "us-east1-test-services"
}

variable "gke_http_load_balancing" {
  type = bool
  description = "enable http_load_balancing"
  default = true
}

variable "gke_enable_horizontal_pod_autoscaling" {
  type = bool
  description = "enable horizontal pod autoscaling"
  default = true
}

variable "gke_enable_network_policy" {
  type = bool
  description = "enable network policies"
  default = true
}

variable "gke_enable_private_endpoint" {
  type = bool
  description = "enable private endpoint"
  default = false
}

variable "gke_enable_private_nodes" {
  type = bool
  description = "enable private nodes"
  default = true
}

variable "gke_istio_enable" {
  type = bool
  description = "enable istio"
  default = false
}

variable "gke_cloudrun_enable" {
  type = bool
  description = "enable cloudrun"
  default = false
}

variable "gke_dns_cache_enable" {
  type = bool
  description = "enable dns caching"
  default = true
}

variable "gke_master_ipv4_cidr_block" {
  type = string
  description = "master ipv4 cidr block range"
  default = "10.0.0.0/28"
}

variable "gke_node_pool_network_tag" {
  type = string
  description = "network tag used for gke node pool"
  default = "gke"
}

variable "custom_auth_networks" {
  type        = list(object({ cidr_block = string, display_name = string }))
  description = "List of custom networks we want to access master from, in addition to the inet ranges"
  default = [
    {
      cidr_block = "10.64.0.0/16"
      display_name = "us-east1-test"
    },
    {
      cidr_block = "10.65.0.0/16"
      display_name = "us-east1-test-db"
    },
    {
      cidr_block = "10.66.0.0/17"
      display_name = "us-east1-test-compute"
    },
  ]
}

variable "mysql_cluster_additional_databases" {
  type = list(object({
    name=string,
    charset=string,
    collation=string
  }))
  description = "list of database data"
  default = [
    {
      name = "chester_dev"
      charset = "utf8"
      collation = "utf8_general_ci"
    }
  ]
}

variable "mysql_cluster_disk_size" {
  description = "size of the master disk"
  type = string
  default = "10"
}

variable "mysql_cluster_tier" {
  description = "the tier for the database"
  type = string
  default = "db-n1-standard-1"
}

variable "mysql_cluster_user_labels" {
  description = "labels for master instance"
  default = {
    "kill" = "me"
  }
  type = map(string)
}

variable "mysql_cluster_sql_user_name" {
  description = "the user name created in sql"
  default = "chester-user"
  type = string
}


variable "mysql_cluster_activation_policy" {
  description = "activation policy"
  default = "ALWAYS"
  type = string
}

variable "mysql_cluster_database_flags" {
  description = "list of database flags"
  default = []
  type = list(object({ name=string, value=string }))
}


variable "mysql_cluster_database_version" {
  type = string
  description = "version of the database"
  default = "MYSQL_5_7"
}

variable "mysql_cluster_deletion_protection" {
  type = bool
  description = "keep the database from being destroyed, set to false cause this is just a test instance"
  default = false
}

variable "mysql_cluster_database_name" {
  type = string
  description = "cluster name"
  default = "chester-test"
}

variable "mysql_cluster_database_region" {
  type = string
  description = "region the db sits in"
  default = "us-east1"
}

variable "mysql_cluster_read_replicas" {
  description = "list of read replicas"
  default = []
  type = list(
  object(
  {
    name=string
    tier=string
    ip_configuration=object(
    {
      ipv4_enabled=bool
      private_network=string
      require_ssl=bool
      authorized_networks=list(
      object(
      {
        value=string
      })
      )
    }
    )
    zone=string
    disk_autoresize=string
    disk_type=string
    disk_size=number
    user_labels=map(string)
    database_flags=list(object({
      name=string
      value=string
    }))
  }))
}

variable "mysql_cluster_require_ssl" {
  type = bool
  description = "require ssl"
  default = false
}



## Tiller default variables
variable "tiller_namespace" {
  description = "name of namespace where we're deploying tiller"
  type = string
  default = "kube-system"
}


variable "tiller_deployment_name" {
  description = "The name to use for the Kubernetes Deployment resource. This should be unique to the Namespace if you plan on having multiple Tiller Deployments in a single Namespace."
  type        = string
  default     = "tiller-deploy"
}

variable "tiller_deployment_labels" {
  description = "Any labels to attach to the Kubernetes Deployment resource."
  type        = map(string)
  default     = {}
}

variable "tiller_deployment_annotations" {
  description = "Any annotations to attach to the Kubernetes Deployment resource."
  type        = map(string)
  default     = {}
}

variable "tiller_deployment_replicas" {
  description = "The number of Pods to use for Tiller. 1 should be sufficient for most use cases."
  type        = number
  default     = 1
}

variable "tiller_service_name" {
  description = "The name to use for the Kubernetes Service resource. This should be unique to the Namespace if you plan on having multiple Tiller Deployments in a single Namespace."
  type        = string
  default     = "tiller-deploy"
}

variable "tiller_service_labels" {
  description = "Any labels to attach to the Kubernetes Service resource."
  type        = map(string)
  default     = {}
}

variable "tiller_service_annotations" {
  description = "Any annotations to attach to the Kubernetes Service resource."
  type        = map(string)
  default     = {}
}

variable "tiller_image" {
  description = "The container image to use for the Tiller Pods."
  type        = string
  default     = "ghcr.io/helm/tiller"
}

variable "tiller_image_version" {
  description = "The version of the container image to use for the Tiller Pods."
  type        = string
  default     = "v2.16.1"
}

variable "tiller_image_pull_policy" {
  description = "Policy for pulling the container image used for the Tiller Pods. Use `Always` if the image tag is mutable (e.g latest)"
  type        = string
  default     = "IfNotPresent"
}

variable "tiller_listen_localhost" {
  description = "If Enabled, Tiller will only listen on localhost within the container."
  type        = bool
  default     = true
}

variable "tiller_history_max" {
  description = "The maximum number of revisions saved per release. Use 0 for no limit."
  type        = number
  default     = 0
}

variable "tiller_sa_name" {
  description = "the name of the tiller service account"
  type = string
  default = "tiller"
}

variable "tiller_cluster_role_binding_name" {
  description = "name of the cluster role binding for tiller"
  type = string
  default = "tiller"
}

## Proxysql helm install defaults

variable "helm_release_namespace" {
  description = "name of the namespace that helm will deploy proxysql"
  type = string
  default = "proxysql"
}

variable "helm_release_values" {
  type = list(string)
  description = "yaml decoded string of the values"
  default = []
}


variable "helm_sets" {
  type = list(map(any))
  description = "list of set values"
  default = []
}

## Cloud Nat Default variables
variable "cloud_nat_create_router" {
  type = bool
  description = "Create router instead of using an existing one, uses 'router' variable for new resource name."
  default     = true
}

variable "cloud_nat_region" {
  type = string
  description = "region for the cloud nat to server"
  default = "us-east1"
}

variable "cloud_nat_router" {
  type = string
  description = "The name of the router in which this NAT will be configured. Changing this forces a new NAT to be created."
  default =  "chester-us-east1-router"
}

## Pubsub default variables
variable "chester_pull_subscriptions" {
  type = list(map(string))
  description = "a list of pull subscriptions and their options"
  default =  [{
    name                  = "chester-test-subscription"
    ack_deadline_seconds  = 10
    max_delivery_attempts = 5
  }]
}


variable "chester_pubsub_topic_name" {
  type = string
  description = "name of the pubsub topic"
  default = "chester-test-topic"
}

variable "chester_topic_labels" {
  type = map(string)
  description = "a map of topic labels"
  default = {
    "someone" = "didnt"
    "set"   = "up"
    "the"   = "labels"
  }
}

## Service Account name
variable "chester_sa_display_name" {
  type = string
  description = "Service Account Display Name"
  default = "chester-test"
}

## K8S Secret stuff
variable "chester_api_secret_name" {
  type = string
  description = "Chester API Secrets name"
  default = "chester-api-secrets"
}


variable "chester_namespace" {
  type = string
  description = "namespace where chester exists in k8s"
  default = "chester"
}

variable "chester_api_secret_type" {
  type = string
  description = "type of the secret to be stored"
  default = "opaque"
}

## KMS Stuff
variable "chester_ring_name" {
  type = string
  description = "name of the chester"
  default = "chester-test-ring"
}
variable "chester_ring_location" {
  type = string
  description = "where the ring is located"
  default = "global"
}
variable "chester_key_name" {
  type = string
  description = "what the key is called"
  default = "chester-test-key"
}
variable "chester_key_rotation_period" {
  type = string
  default = "0"
  description = "time until the key rotates"
}
variable "chester_key_purpose" {
  type = string
  default = "ASYMMETRIC_DECRYPT"
  description = "Purpose, defined by https://cloud.google.com/kms/docs/algorithms#key_purposes"
}
variable "chester_key_algorithm" {
  type = string
  description = "algorithm to be used"
  default = "RSA_DECRYPT_OAEP_4096_SHA512"
}

## IAP Stuff
variable "chester_iap_display_name" {
  type = string
  description = "name of the iap client"
  default = "chester-iap-client"
}

## IP Address stuff
variable "static_ip_name" {
  type = string
  description = "chester IP address name"
  default = "chester-test-external-ip"
}

variable "static_ip_description" {
  type = string
  description = "description of the IP address"
  default = "chester test external ip"
}

variable "static_ip_address_type" {
  type = string
  description = "public or private"
  default = "EXTERNAL"
}

## Firewall stuff
variable "chester_firewall_name" {
  type = string
  description = "name of the firewall rule"
  default = "chester-test-iap"
}

variable "chester_firewall_priority" {
  type = number
  description = "priority of the firewall rule"
  default = 1000
}

variable "chester_firewall_allow_protocols" {
  description = "list of objects that contain data about allow rules"
  type = list(object({
    protocol = string,
    ports = list(string)
  }))
  default = [
    {
      protocol = "tcp",
      ports = ["80","443","31278"]
    }
  ]
}

variable "chester_firewall_deny_protocols" {
  description = "list of objects that contain data about allow rules"
  type = list(object({
    protocol = string,
    ports = list(string)
  }))
  default = []
}


variable "chester_firewall_target_tags" {
  description = "target tags we're allowing access to"
  type = list(string)
  default = ["gke"]
}

variable "chester_firewall_source_tags" {
  description = "incoming tags"
  type = list(string)
  default = []
}

variable "chester_firewall_source_ranges" {
  description = "CIDR ranges we allow access to"
  type = list(string)
  default = ["130.211.0.0/22","35.191.0.0/16","35.235.240.0/20"]
}

variable "chester_firewall_service_accounts" {
  description = "service accounts that have access to this"
  type = list(string)
  default = []
}


## Pull secret default variables
variable "docker_pull_secret_name" {
  type = string
  default = "gcr"
  description = "name of the pull secret"
}

## Chester API Helm stuff
variable "chester_api_image_name" {
  type = string
  description = "full image name of the chester-api server"
}

variable "chester_api_image_pull_policy" {
  type = string
  description = "pull policy for the chester image"
  default = "IfNotPresent"
}







