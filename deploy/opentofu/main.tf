locals {
  # Fetch global configuration options, such as the Talos version.
  global_config = yamldecode(file("${path.cwd}/config.yaml"))

  # Fetch the list of all clusters.
  taloscluster_path  = "${path.cwd}/deploy/manifests/talosclusters"
  taloscluster_files = fileset(local.taloscluster_path, "*.yaml")
  taloscluster_configs = {
    for filename in local.taloscluster_files :
    replace(filename, ".yaml", "") => yamldecode(file("${local.taloscluster_path}/${filename}"))
  }

  # Fetch the list of all machines.
  machines_path  = "${path.cwd}/deploy/manifests/machines"
  machines_files = fileset(local.machines_path, "*.yaml")
  machines_configs = [
    for filename in local.machines_files :
    yamldecode(file("${local.machines_path}/${filename}"))
  ]
}

module "talos_cluster" {
  source = "../../modules/talos_cluster"

  for_each = local.taloscluster_configs

  global_config = local.global_config
  machines      = local.machines_configs

  config = local.taloscluster_configs[each.key]
}

output "kubeconfigs" {
  description = "A map of kubeconfigs for each cluster."
  depends_on  = [module.talos_cluster]
  value = {
    for cluster in local.taloscluster_configs :
    cluster.metadata.name => module.talos_cluster[cluster.metadata.name].kubeconfig_user != null ? module.talos_cluster[cluster.metadata.name].kubeconfig_user : ""
  }
}

provider "helm" {
  kubernetes {
    host                   = module.talos_cluster["hel01"].kubeconfig_credentials.host
    cluster_ca_certificate = base64decode(module.talos_cluster["hel01"].kubeconfig_credentials.ca_certificate)
    client_certificate     = base64decode(module.talos_cluster["hel01"].kubeconfig_credentials.client_certificate)
    client_key             = base64decode(module.talos_cluster["hel01"].kubeconfig_credentials.client_key)
  }
}

resource "helm_release" "authz" {
  name      = "authz"
  namespace = "kube-system"
  chart     = "${path.cwd}/charts/authz"
  atomic    = true

  values = [
    yamlencode({
      global = {
        clusterAdmins = [
          for admin in local.global_config.kubernetes.admins : admin.email
        ]
      }
    })
  ]
}
