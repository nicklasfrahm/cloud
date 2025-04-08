output "talosconfig_admin" {
  description = "The Talos config for the cluster."
  depends_on  = [data.talos_client_configuration.this]
  value       = data.talos_client_configuration.this.talos_config
}

output "kubeconfig_user" {
  description = "The kubeconfig file for the cluster."
  depends_on  = [talos_cluster_kubeconfig.this]
  value       = local.kubeconfig
}

output "kubeconfig_credentials" {
  description = "The kubeconfig file for the cluster."
  depends_on  = [talos_cluster_kubeconfig.this]
  value = {
    host               = talos_cluster_kubeconfig.this.kubernetes_client_configuration.host
    ca_certificate     = talos_cluster_kubeconfig.this.kubernetes_client_configuration.ca_certificate
    client_certificate = sensitive(talos_cluster_kubeconfig.this.kubernetes_client_configuration.client_certificate)
    client_key         = sensitive(talos_cluster_kubeconfig.this.kubernetes_client_configuration.client_key)
  }
}
