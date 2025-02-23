# Get the machine pools of the control planes.
data "http" "machinepools_controlplanes" {
  # Convert the machine pools into a list.
  for_each = [ for pool in var.config.spec.infrastructure.controlPlane.machinePools : pool.name ]

  url = "https://cloud.${var.domain}/v1beta1/machinepools/${each.value.name}.json"
}

# Get the machine pool nodes.
data "http" "nodes" {
  url = "https://cloud.${var.domain}/v1beta1/machines/index.json"
}

locals {
  # Get the control plane machine pools.
  controlplane_machinepools = [ for pool in jsondecode(data.http.machinepools_controlplanes["controlplane"].body) : pool ]

  # Get the control plane nodes.
  controlplane_nodes = [ for node in jsondecode(data.http.nodes.body) : node if node.metadata.labels["machinepool"] == "controlplane" ]
}
