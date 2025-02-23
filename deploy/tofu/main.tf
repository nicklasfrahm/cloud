locals {
  # Load all region configuration files.
  kubernetescluster_path = "${path.cwd}/deploy/manifests/kubernetesclusters"
  kubernetescluster_files = fileset(local.kubernetescluster_path, "*.yaml")
  kubernetescluster_configs = {
    for filename in local.kubernetescluster_files :
    replace(filename, ".yaml", "") => yamldecode(file("${local.kubernetescluster_path}/${filename}"))
  }

  config = yamldecode(file("${path.cwd}/config.yaml"))
}

module "kubernetescluster" {
  source = "./modules/kubernetescluster"

  for_each = local.kubernetescluster_configs

  config = each.value

  domain = local.config.dns.domain
}
