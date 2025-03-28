terraform {
  required_version = ">= 1.9.0"

  backend "gcs" {
    bucket = "nicklasfrahm"
    prefix  = "tofu/state/homelab"
  }
}
