locals {
  patch_dir = "${path.module}/patches"
}

# Create the Talos secret bundle.
resource "talos_machine_secrets" "secret_bundle" {
  talos_version = var.global_config.talos.version
}

locals {
  // TODO: Use the IP address of the node.
  cluster_endpoint = var.config.spec.infrastructure.loadBalancer.host ? "https://${var.config.spec.infrastructure.loadBalancer.host}:${var.config.spec.infrastructure.loadBalancer.port}" : "https://${var.region.metadata.name}:6443"
}

# Define the controlplane configuration.
data "talos_machine_configuration" "controlplane" {
  cluster_name = var.config.metadata.name
  cluster_endpoint = local.cluster_endpoint
  machine_type = "controlplane"
  machine_secrets = talos_machine_secrets.this.secrets

  config_patches = [
    yamldecode(templatefile("${local.patch_dir}/cilium.yaml", {
      cilium_operator_replicas = length(var.config.spec.infrastructure.controlplanes) > 1 ? 2 : 1
    })),
    yamldecode(templatefile("${local.patch_dir}/scheduling.yaml", {
      allow_scheduling_on_control_planes = length(var.config.spec.infrastructure.workers) > 0 ? "false" : "true"
    })),
    // TODO: Add more patches as needed.
  ]
}

# Define the controlplane configuration.
data "talos_machine_configuration" "worker" {
  cluster_name = var.config.metadata.name
  cluster_endpoint = local.cluster_endpoint
  machine_type = "controlplane"
  machine_secrets = talos_machine_secrets.this.secrets
}
