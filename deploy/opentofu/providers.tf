terraform {
  required_version = ">= 1.9.0"

  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "3.0.2"
    }
    talos = {
      source  = "siderolabs/talos"
      version = "0.9.0"
    }
  }

  backend "gcs" {
    bucket = "nicklasfrahm"
    prefix = "tofu/state/homelab"
  }
}
