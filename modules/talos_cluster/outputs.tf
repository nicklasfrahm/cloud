output "kubeconfig" {
  description = "The kubeconfig file for the cluster."
  depends_on = [ talos_cluster_kubeconfig.this ]
  value = local.kubeconfig
}
