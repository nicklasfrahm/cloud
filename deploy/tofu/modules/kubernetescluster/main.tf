locals {
  name = var.config.metadata.name
  talos_version = yamldecode(file("${path.cwd}/config.yaml")).talos.version
}

# Create the Talos secret bundle for the cluster.
resource "talos_machine_secrets" "this" {
  talos_version = local.talos_version
}

# Create the control plane configuration.
data "talos_machine_configuration" "this" {
  cluster_name     = local.name
  machine_type     = "controlplane"
  cluster_endpoint = "https://k8s.${local.name}.${var.domain}:6443"
  machine_secrets  = talos_machine_secrets.this.machine_secrets
}

# Get control-plane node IP addresses.
data "dns_a_record_set" "controlplanes" {
  host = "k8s.${local.name}.${var.domain}"
}

# data "talos_client_configuration" "this" {
#   cluster_name         = local.name
#   client_configuration = talos_machine_secrets.this.client_configuration
#   nodes                = data.dns_a_record_set.controlplanes.addrs
# }


# resource "talos_machine_configuration_apply" "this" {
#   client_configuration        = talos_machine_secrets.this.client_configuration
#   machine_configuration_input = data.talos_machine_configuration.this.machine_configuration
#   node                        = "10.5.0.2"
#   config_patches = [
#     yamlencode({
#       machine = {
#         install = {
#           disk = "/dev/sdd"
#         }
#       }
#     })
#   ]
# }

# resource "talos_machine_bootstrap" "this" {
#   depends_on = [
#     talos_machine_configuration_apply.this
#   ]
#   node                 = "10.5.0.2"
#   client_configuration = talos_machine_secrets.this.client_configuration
# }
