terraform {
  required_version = ">= 1.9.0"

  required_providers {
    talos = {
      source = "siderolabs/talos"
      version = "0.7.1"
    }
    dns = {
      source = "hashicorp/dns"
      version = "3.4.2"
    }
    http = {
      source = "hashicorp/http"
      version = "3.4.5"
    }
  }
}
