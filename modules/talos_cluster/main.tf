locals {
  cluster_name = var.config.metadata.name
  patch_dir = "${path.module}/patches"

  controlplane_names = [ for controlplane in var.config.spec.infrastructure.controlplanes : controlplane.name ]
  controlplane_machines = {
    for machine in var.machines :
    machine.metadata.name => machine if contains(local.controlplane_names, machine.metadata.name)
  }

  worker_names = [ for worker in var.config.spec.infrastructure.workers : worker.name ]
  worker_machines = {
    for machine in var.machines :
    machine.metadata.name => machine if contains(local.worker_names, machine.metadata.name)
  }

  # Determine the cluster endpoint based on the load balancer configuration.
  cluster_endpoint = (
    var.config.spec.infrastructure.loadBalancer.host != "" ?
      "https://${var.config.spec.infrastructure.loadBalancer.host}:${var.config.spec.infrastructure.loadBalancer.port}" :
      "https://${local.controlplane_names[0]}.${var.global_config.dns.zone}:6443"
  )

  # Configuration patches for the cluster.
  config_patches = [
    templatefile("${local.patch_dir}/base.yaml", {
      cilium_operator_replicas = length(local.controlplane_names) > 1 ? 2 : 1,
      allow_scheduling_on_control_planes = length(local.worker_names) > 0 ? "false" : "true",
    }),
    // TODO: Add more patches as needed.
  ]
}

# Create the Talos secret bundle.
resource "talos_machine_secrets" "secret_bundle" {
  talos_version = var.global_config.talos.version
}

# Create Talos client configuration.
data "talos_client_configuration" "this" {
  cluster_name = local.cluster_name
  client_configuration = talos_machine_secrets.secret_bundle.client_configuration
}

moved {
  from = talos_machine_secrets.this
  to = talos_machine_secrets.secret_bundle
}

# Define the controlplane configuration.
data "talos_machine_configuration" "controlplane" {
  cluster_name = local.cluster_name
  cluster_endpoint = local.cluster_endpoint
  machine_type = "controlplane"
  machine_secrets = talos_machine_secrets.secret_bundle.machine_secrets

  config_patches = local.config_patches

  for_each = local.controlplane_machines
}

# Define the controlplane configuration.
data "talos_machine_configuration" "worker" {
  cluster_name = local.cluster_name
  cluster_endpoint = local.cluster_endpoint
  machine_type = "worker"
  machine_secrets = talos_machine_secrets.secret_bundle.machine_secrets

  config_patches = local.config_patches

  for_each = local.worker_machines
}
