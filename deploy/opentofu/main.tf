locals {
  # Fetch global configuration options, such as the Talos version.
  global_config = yamldecode(file("${path.cwd}/config.yaml"))

  # Fetch the list of all clusters.
  taloscluster_path = "${path.cwd}/deploy/manifests/talosclusters"
  taloscluster_files = fileset(local.taloscluster_path, "*.yaml")
  taloscluster_configs = {
    for filename in local.taloscluster_files :
    replace(filename, ".yaml", "") => yamldecode(file("${local.taloscluster_path}/${filename}"))
  }

  # Fetch the list of all machines.
  machines_path = "${path.cwd}/deploy/manifests/machines"
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
  machines = local.machines_configs

  config = each.value
}
