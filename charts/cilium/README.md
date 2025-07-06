# Cilium

This chart is not official. Instead, it is an opinionated umbrella chart of the
official Cilium Helm chart. It is tuned for use with Talos and Kubernetes
clusters that are not using kube-proxy.

## Features

- **Capability tuning**  
  This chart limits the capabilities to a set allowed by Talos.

- **No kube-proxy**  
  This chart is disables kube-proxy by default and configures Cilium to
  function without it.

- **Built-in load balancer**  
  This chart enables `nodeipam` to allow the creation of `LoadBalancer`
  services without external infrastructure.

- **Gateway API support**  
  This chart enables the Cilium Gateway API controller to allow the creation of
  `Gateway` and `HTTPRoute` resources.

- **Gateway API CRDs**  
  This chart installs the Gateway API CRDs, which are required for the
  Gateway API controller to function. It does this using a `pre-install`
  hook to ensure that the CRDs are installed before the controller is started.

- **Hubble support**  
  This chart enables Hubble for network observability, including the Hubble UI
  and Hubble Relay. This allows you to monitor network traffic and events in
  real-time.
