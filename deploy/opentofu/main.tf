locals {
  # Load all region configuration files.
  region_path = "${path.cwd}/deploy/manifests/regions"
  region_files = fileset(local.region_path, "*.yaml")
  region_configs = {
    for filename in local.region_files :
    replace(filename, ".yaml", "") => yamldecode(file("${local.region_path}/${filename}"))
  }
}

module "talos_cluster" {
  source = "../../modules/talos_cluster"

  for_each = local.region_configs

  region = each.value
}
