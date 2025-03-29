terraform {
  required_version = ">= 1.9.0"

  required_providers {
    talos = {
      source = "siderolabs/talos"
      version = "0.7.1"
    }
  }
}
