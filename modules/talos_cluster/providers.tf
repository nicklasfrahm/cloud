terraform {
  required_version = ">= 1.9.0"

  required_providers {
    local = {
      source = "hashicorp/local"
      version = "2.5.2"
    }
    talos = {
      source = "siderolabs/talos"
      version = "0.7.1"
    }
    random = {
      source = "hashicorp/random"
      version = "3.7.1"
    }
  }
}
