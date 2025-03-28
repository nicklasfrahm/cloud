locals {
  name = var.region.metadata.name
  talos_version = yamldecode(file("${path.cwd}/config.yaml")).talos.version
}

# Create the Talos secret bundle.
resource "talos_machine_secrets" "this" {
  talos_version = local.talos_version
}

resource "random_id" "this" {
  byte_length = 8
}

# Create a file with the Talos version.
resource "local_file" "version" {
  filename = "${path.cwd}/deploy/tofu/out/${local.name}"
  content = "${local.talos_version}/${random_id.this.hex}"
}
