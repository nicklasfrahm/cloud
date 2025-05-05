terraform {
  required_version = ">= 1.9.0"

  required_providers {
    talos = {
      source  = "siderolabs/talos"
      version = "0.8.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.17.0"
    }
  }
}
