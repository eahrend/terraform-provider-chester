module "project-factory" {
  source               = "terraform-google-modules/project-factory/google"
  version              = "~> 10.1"
  activate_apis        = var.google_project_factory_activate_apis
  name                 = var.google_project_factory_name
  random_project_id    = var.google_project_factory_random_project_id
  org_id               = var.google_project_factory_org_id
  billing_account      = var.google_project_factory_billing_account
  auto_create_network  = var.google_project_factory_auto_create_network
}


resource "google_compute_network" "network" {
  depends_on                      = [module.project-factory]
  name                            = format("%s-%s", module.project-factory.project_id, var.vpc_name)
  auto_create_subnetworks         = var.network_auto_create_subnetworks
  routing_mode                    = var.network_routing_mode
  project                         = module.project-factory.project_id
  description                     = var.network_description
  delete_default_routes_on_create = var.network_delete_default_internet_gateway_routes
}

resource "google_compute_subnetwork" "subnetwork" {
  depends_on = [module.project-factory]
  for_each                 = local.subnets
  name                     = each.value.subnet_name
  ip_cidr_range            = each.value.cidr
  region                   = each.value.subnet_region
  private_ip_google_access = lookup(each.value, "subnet_private_access", "false")
  dynamic "log_config" {
    for_each = lookup(each.value, "subnet_flow_logs", false) ? [{
      aggregation_interval = lookup(each.value, "subnet_flow_logs_interval", "INTERVAL_5_SEC")
      flow_sampling        = lookup(each.value, "subnet_flow_logs_sampling", "0.5")
      metadata             = lookup(each.value, "subnet_flow_logs_metadata", "INCLUDE_ALL_METADATA")
    }] : []
    content {
      aggregation_interval = log_config.value.aggregation_interval
      flow_sampling        = log_config.value.flow_sampling
      metadata             = log_config.value.metadata
    }
  }
  network     = google_compute_network.network.id
  project     = module.project-factory.project_id
  description = lookup(each.value, "description", null)
  secondary_ip_range = yamldecode(each.value.secondary_ranges)
}

locals {
  depends_on = [module.project-factory]
  subnets = {
  for x in var.vpc_subnets :
    "${x.subnet_region}/${x.subnet_name}" => x
  }
}







module "gke" {
  depends_on = [module.project-factory,google_compute_subnetwork.subnetwork,google_compute_network.network]
  version                    = "16.1.0"
  source                     = "terraform-google-modules/kubernetes-engine/google//modules/beta-private-cluster"
  project_id                 = module.project-factory.project_id
  name                       = var.gke_cluster_name
  region                     = var.gke_region
  zones                      = var.gke_zones
  network                    = google_compute_network.network.name
  network_project_id         = module.project-factory.project_id
  subnetwork                 = "us-east1-test-compute"
  ip_range_pods              = var.gke_ip_range_pods_name
  ip_range_services          = var.gke_ip_range_services_name
  http_load_balancing        = var.gke_http_load_balancing
  horizontal_pod_autoscaling = var.gke_enable_horizontal_pod_autoscaling
  network_policy             = var.gke_enable_network_policy
  enable_private_endpoint    = var.gke_enable_private_endpoint
  enable_private_nodes       = var.gke_enable_private_nodes
  master_ipv4_cidr_block     = var.gke_master_ipv4_cidr_block
  kubernetes_version         = var.gke_kubernetes_version
  istio                      = var.gke_istio_enable
  cloudrun                   = var.gke_cloudrun_enable
  dns_cache                  = var.gke_dns_cache_enable
  node_pools = [
    {
      name               = var.gke_node_name
      machine_type       = var.gke_node_machine_type
      min_count          = var.gke_node_min_count
      max_count          = var.gke_node_max_count
      local_ssd_count    = var.gke_local_ssd_count
      disk_size_gb       = var.gke_disk_size_gb
      disk_type          = var.gke_disk_type
      image_type         = var.gke_image_type
      auto_repair        = var.gke_auto_repair
      auto_upgrade       = var.gke_auto_upgrade
      preemptible        = var.gke_premptible_nodes
      initial_node_count = 3
      version             = var.gke_kubernetes_version
    },
  ]

  node_pools_oauth_scopes = {
    all = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
    default-node-pool = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
  node_pools_labels = {
    all = {
      default-node-pool = false
    }
  }
  node_pools_metadata = {
    all = {
      node-pool-metadata-custom-value = "my-node-pool"
    }
  }
  node_pools_tags = {
    all = [
      var.gke_node_pool_network_tag,
    ]
  }
  master_authorized_networks = local.custom_auth_networks
}

locals {
  inet_mapping = [
  for k,v in null_resource.inet_mappings.*.triggers:
  {
    cidr_block = v.cidr_block
    display_name = v.display_name
  }
  ]
  auth_network_mapping = [
  for k,v in var.custom_auth_networks: {
    cidr_block = v.cidr_block
    display_name = v.display_name
  }
  ]
  custom_auth_networks = concat(local.inet_mapping,local.auth_network_mapping)
}

resource "null_resource" "inet_mappings" {
  count = length(var.inet_ranges)

  triggers = {
    cidr_block      = element(var.inet_ranges, count.index)
    display_name    = element(var.inet_ranges, count.index)
  }
}
resource "google_compute_global_address" "private_ip_alloc" {
  project       = module.project-factory.project_id
  name          = "chester-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.network.self_link
}


resource "google_service_networking_connection" "peering_network_connection" {
  network                 = google_compute_network.network.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}


module "mysql_cluster" {
  depends_on = [google_service_networking_connection.peering_network_connection]
  source  = "GoogleCloudPlatform/sql-db/google//modules/mysql"
  version = "5.1.0"
  additional_databases      = var.mysql_cluster_additional_databases
  disk_size                 = var.mysql_cluster_disk_size
  tier                      = var.mysql_cluster_tier
  user_labels               = var.mysql_cluster_user_labels
  user_name                 = var.mysql_cluster_sql_user_name
  activation_policy         = var.mysql_cluster_activation_policy
  database_flags            = var.mysql_cluster_database_flags
  database_version          = var.mysql_cluster_database_version
  deletion_protection       = var.mysql_cluster_deletion_protection
  name                      = var.mysql_cluster_database_name
  project_id                = module.project-factory.project_id
  read_replica_name_suffix  = ""
  region                    = var.mysql_cluster_database_region
  # just gonna throw master in B cause it works
  zone                      = "${var.mysql_cluster_database_region}-b"
  read_replicas             = [
    {
      name = "-1"
      tier = var.mysql_cluster_tier
      ip_configuration = {
        private_network = format("projects/%s/global/networks/%s", module.project-factory.project_id, google_compute_network.network.name)
        ipv4_enabled    = false
        require_ssl     = false
        authorized_networks = [{
          value = "0.0.0.0/0"
        }]
      }
      zone = "us-east1-d"
      disk_type = "PD_SSD"
      disk_autoresize = true
      disk_size = var.mysql_cluster_disk_size
      database_flags = []
      user_labels = var.mysql_cluster_user_labels
    },
    {
      database_flags = []
      disk_type = "PD_SSD"
      name = "-2"
      tier = var.mysql_cluster_tier
      ip_configuration = {
        private_network = format("projects/%s/global/networks/%s", module.project-factory.project_id, google_compute_network.network.name)
        ipv4_enabled    = false
        require_ssl     = false
        authorized_networks = [{
          value = "0.0.0.0/0"
        }]
      }
      zone = "us-east1-c"
      disk_autoresize = true
      disk_size = var.mysql_cluster_disk_size
      user_labels = var.mysql_cluster_user_labels
    }]
  ip_configuration = {
    private_network = format("projects/%s/global/networks/%s", module.project-factory.project_id, google_compute_network.network.name)
    ipv4_enabled    = false
    require_ssl     = var.mysql_cluster_require_ssl
    authorized_networks = [{
      value = "0.0.0.0/0"
    }]
  }
}


locals {
  tiller_listen_localhost_arg = var.tiller_listen_localhost ? ["--listen=localhost:44134"] : []
}





resource "helm_release" "proxysql" {
  depends_on = [module.gke,kubernetes_namespace.proxysql]
  name       = "proxysql"
  chart      = "./helm/proxysql/Chart/proxysql"
  namespace  = var.helm_release_namespace
  values     = [
    "${file("./helm/proxysql/values/test.yaml")}"
  ]
  set {
    name  = "sql_writer"
    value = module.mysql_cluster.private_ip_address
  }

  set {
    name  = "sql_reader_one"
    value = module.mysql_cluster.replicas[0].private_ip_address
  }
  set {
    name  = "sql_reader_two"
    value = module.mysql_cluster.replicas[1].private_ip_address
  }
  set {
    name = "sql_username"
    value = "chester-user"
  }
  set {
    name  = "sql_password"
    value = module.mysql_cluster.generated_user_password
  }
  force_update = true
  dynamic "set" {
    for_each = [for s in var.helm_sets: {
      name   = s.name
      value =  s.value
    }]
    content {
      name  = set.value.name
      value = set.value.value
    }
  }
}



module "cloud_nat" {
  source        = "terraform-google-modules/cloud-nat/google"
  version       = "~> 1.2"
  project_id    = module.project-factory.project_id
  region        = var.cloud_nat_region
  router        = var.cloud_nat_router
  create_router = var.cloud_nat_create_router
  name          = "${var.cloud_nat_region}-nat-gateway"
  network       = google_compute_network.network.self_link
}


## Pubsub stuff

module "pubsub" {
  source  = "terraform-google-modules/pubsub/google"
  version = "~> 1.4.0"
  topic              = var.chester_pubsub_topic_name
  topic_labels       = var.chester_topic_labels
  project_id         = module.project-factory.project_id
  pull_subscriptions = var.chester_pull_subscriptions
}

## KMS stuff that we use to decrypt data in datastore
resource "google_kms_key_ring" "keyring" {
  name      = var.chester_ring_name
  location  = var.chester_ring_location
  project   = module.project-factory.project_id
}



resource "google_kms_crypto_key" "key" {
  name              = var.chester_key_name
  key_ring          = google_kms_key_ring.keyring.id
  rotation_period   = var.chester_key_rotation_period != "0" ? var.chester_key_rotation_period : null
  purpose           = var.chester_key_purpose
  version_template {
    algorithm = var.chester_key_algorithm
  }
  lifecycle {
    ## normally we don't want to prevent destroy, but this stuff ain't coming up anytime soon
    prevent_destroy = false
  }
}




module "service_accounts" {
  source        = "terraform-google-modules/service-accounts/google"
  version       = "3.0.1"
  names         = [var.chester_sa_display_name]
  project_id    = module.project-factory.project_id
  prefix        = "sa"
  generate_keys = true
}

resource "google_pubsub_subscription_iam_member" "pubsub_subscription_iam_additive_viewer" {
  project      = module.project-factory.project_id
  subscription = module.pubsub.subscription_names[0]
  role         = "roles/pubsub.viewer"
  member       = "serviceAccount:${module.service_accounts.email}"
}

resource "google_pubsub_subscription_iam_member" "pubsub_subscription_iam_additive_editor" {
  project      = module.project-factory.project_id
  subscription = module.pubsub.subscription_names[0]
  role         = "roles/pubsub.editor"
  member       = "serviceAccount:${module.service_accounts.email}"
}


resource "google_pubsub_topic_iam_member" "pubsub_topic_iam_additive_viewer" {
  topic        = module.pubsub.topic
  project      = module.project-factory.project_id
  role         = "roles/pubsub.publisher"
  member       = "serviceAccount:${module.service_accounts.email}"
}

resource "google_pubsub_topic_iam_member" "pubsub_topic_iam_additive_publisher" {
  topic        = module.pubsub.topic
  project      = module.project-factory.project_id
  role         = "roles/pubsub.viewer"
  member       = "serviceAccount:${module.service_accounts.email}"
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_viewer" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.publicKeyViewer"
  members       = ["serviceAccount:${module.service_accounts.email}"]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_encrypter" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"
  members       = ["serviceAccount:${module.service_accounts.email}"]
}
resource "google_kms_crypto_key_iam_binding" "cryptoKeyDecrypter" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.publicKeyViewer"
  members       = ["serviceAccount:${module.service_accounts.email}"]
}




resource "google_project_iam_binding" "datastore_user" {
  project = module.project-factory.project_id
  role    = "roles/datastore.user"
  members = [
    "serviceAccount:${module.service_accounts.email}",
  ]
}

resource "google_project_iam_binding" "sql_admin" {
  project = module.project-factory.project_id
  role    = "roles/cloudsql.admin"
  members = [
    "serviceAccount:${module.service_accounts.email}",
  ]
}

resource "kubernetes_namespace" "proxysql" {
  depends_on = [module.gke]
  metadata {
    annotations = {
      name = "proxysql"
    }

    name = "proxysql"
  }
}

## K8S secret resource
resource "kubernetes_namespace" "chester" {
  depends_on = [module.gke]
  metadata {
    annotations = {
      name = var.chester_namespace
    }

    name = var.chester_namespace
  }
}


resource "kubernetes_secret" "chester_api" {
  depends_on = [module.gke,kubernetes_namespace.chester]
  metadata {
    name = var.chester_api_secret_name
    namespace = var.chester_namespace
  }
  data = {
    topic_key = module.service_accounts.key
    subscription_key = module.service_accounts.key
    sqladmin_key = module.service_accounts.key
    datastore_key = module.service_accounts.key
    kms_key = module.service_accounts.key
  }
  type = var.chester_api_secret_type
}
## IAP Stuff
resource "google_iap_brand" "chester_brand" {
  support_email     = var.allowed_user_emails[0]
  application_title = "Chester App"
  project           = module.project-factory.project_id
}

resource "google_iap_client" "chester_client" {
  display_name = var.chester_iap_display_name
  brand        = google_iap_brand.chester_brand.name
}



## Static IP
resource "google_compute_global_address" "chester_ip_address" {
  name = var.static_ip_name
  project = module.project-factory.project_id
  description = var.static_ip_description
  address_type = var.static_ip_address_type
}

## Firewall rule for allowing IAP access to the cluster
resource "google_compute_firewall" "chester_firewall" {
  project = module.project-factory.project_id
  name    = var.chester_firewall_name
  network = google_compute_network.network.self_link
  priority = var.chester_firewall_priority
  dynamic "allow" {
    for_each = var.chester_firewall_allow_protocols
    content {
      protocol = allow.value["protocol"]
      ports = allow.value["ports"]
    }
  }
  dynamic "deny" {
    for_each = var.chester_firewall_deny_protocols
    content {
      protocol = deny.value["protocol"]
      ports = deny.value["ports"]
    }
  }
  target_tags               = length(var.chester_firewall_target_tags)      > 0 ? var.chester_firewall_target_tags      : null
  source_ranges             = length(var.chester_firewall_source_ranges)    > 0 ? var.chester_firewall_source_ranges    : null
}

## Docker pull secret
resource "kubernetes_secret" "docker_pull_secret" {
  depends_on = [module.gke,kubernetes_namespace.chester]
  metadata {
    name = var.docker_pull_secret_name
    namespace = var.chester_namespace
  }

  data = {
    ".dockerconfigjson" = base64decode(var.docker_pull_secret_string)
  }

  type = "kubernetes.io/dockerconfigjson"
}

## Chester API release
## TODO: Will dynamically get service ID
resource "helm_release" "chester_api" {
  depends_on = [module.gke,kubernetes_namespace.chester]
  name       = "chester-api"
  chart      = "./helm/chester-api/chart/chester-api/"
  namespace  = var.chester_namespace
  values     = [
    "${file("./helm/chester-api/values/test.yaml")}"
  ]
  set {
    name = "c_env.BACKEND_SERVICE_ID"
    # this will fail, you need to manually put this in
    value = "000000000"
  }
  set {
    name = "c_env.BASIC_AUTH_ENABLED"
    value = "false"
  }
  set {
    name = "service.omitClusterIP"
    value = "true"
  }
  set {
    name = "c_env.PROJECT_ID"
    value = module.project-factory.project_id
  }
  set {
    name = "c_env.PUBSUB_TOPIC"
    value = module.pubsub.topic
  }
  set {
    name = "c_env.PUBSUB_SUBSCRIPTION"
    value = module.pubsub.subscription_names[0]
  }
  set {
    name = "c_env.PROJECT_NUMBER"
    value = module.project-factory.project_number
  }
  set {
    name = "c_env.KMS_KEY_NAME"
    value = google_kms_crypto_key.key.name
  }
  set {
    name = "c_env.KMS_KEY_RING"
    value = google_kms_key_ring.keyring.name
  }
  set {
    name = "c_env.KMS_KEY_VERSION"
    value = "1"
  }
  set {
    name = "c_env.KMS_KEY_LOCATION"
    value = "global"
  }
  set {
    name = "namespace"
    value = var.chester_namespace
  }
  set{
    name = "secret.name"
    value = var.chester_api_secret_name
  }
  set {
    name  = "image.name"
    value = var.chester_api_image_name
  }
  set {
    name = "image.pullSecret"
    value = var.docker_pull_secret_name
  }
  set {
    name = "image.pullPolicy"
    value = var.chester_api_image_pull_policy
  }
  set {
    name = "ingress.annotations.networking\\.gke\\.io/managed-certificates"
    value = "chester-test"
  }
  set {
    name = "ingress.annotations.kubernetes\\.io/ingress\\.global-static-ip-name"
    value = google_compute_global_address.chester_ip_address.name
  }
  set {
    name = "oauthsecret"
    value = kubernetes_secret.oauth_secret.metadata.0.name
  }
  force_update = false
}

## Chester Daemon Release
resource "helm_release" "chester_daemon" {
  depends_on = [module.gke,kubernetes_namespace.chester]
  name       = "chester-daemon"
  chart      = "./helm/chester-daemon/chart/chester-daemon/"
  namespace  = var.chester_namespace
  values     = [
    "${file("./helm/chester-daemon/values/test.yaml")}"
  ]
  set {
    name = "c_env.PROJECT_ID"
    value = module.project-factory.project_id
  }
  set {
    name = "c_env.PUBSUB_TOPIC"
    value = module.pubsub.topic
  }
  set {
    name = "c_env.PUBSUB_SUBSCRIPTION"
    value = module.pubsub.subscription_names[0]
  }
  set {
    name = "c_env.NETWORK_PROJECT_ID"
    value = module.project-factory.project_id
  }
  set {
    name = "c_env.NETWORK_NAME"
    value = google_compute_network.network.name
  }
  set {
    name = "namespace"
    value = var.chester_namespace
  }
  set{
    name = "secret.name"
    value = var.chester_api_secret_name
  }
  set {
    name  = "image.name"
    value = var.chester_api_image_name
  }
  set {
    name = "image.pullSecret"
    value = var.docker_pull_secret_name
  }
  set {
    name = "image.pullPolicy"
    value = var.chester_api_image_pull_policy
  }
  force_update = true
}



## Get DNS parent entry

data "google_dns_managed_zone" "managed_zone" {
  name = var.dns_managed_zone_name
  project = var.dns_managed_zone_project
}

## Set DNS record
resource "google_dns_record_set" "resource_recordset" {
  project = var.dns_managed_zone_project
  managed_zone = data.google_dns_managed_zone.managed_zone.name
  name         = var.chester_dns_name
  type         = "A"
  rrdatas      = [google_compute_global_address.chester_ip_address.address]
  ttl          = 300
}


## Seeing if I can turn on datastore this way
resource "google_app_engine_application" "datastore" {
  project     = module.project-factory.project_id
  location_id = "us-east1"
  database_type = "CLOUD_DATASTORE_COMPATIBILITY"
}




// set up secret:
resource "kubernetes_secret" "oauth_secret" {
  metadata {
    name = "iap-secret"
    namespace = var.chester_namespace
  }

  data = {
    client_id = google_iap_client.chester_client.client_id
    client_secret = google_iap_client.chester_client.secret
  }

  type = "Opaque"
}


resource "google_iap_web_backend_service_iam_binding" "binding" {
  # This needs to be filled in manually
  web_backend_service = "000000"
  project = module.project-factory.project_id
  role = "roles/iap.httpsResourceAccessor"
  members = var.allowed_user_emails
}


## Data source, commenting out cause it's currently bugged, see https://github.com/hashicorp/terraform-provider-kubernetes/issues/1400
/*
data "kubernetes_ingress" "chester_ingress" {
  depends_on = [module.gke.ca_certificate,helm_release.chester_api]
  metadata {
    name = "chester-api-ing"
    namespace = var.chester_namespace
    annotations = {
      "ingress.kubernetes.io/backends" = ""
    }
  }
}
*/