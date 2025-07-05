# cert-manager Helm Chart

A Helm chart for deploying cert-manager with Gateway API support enabled.

## Overview

This chart deploys cert-manager, a native Kubernetes certificate management controller that automates the issuance and renewal of TLS certificates. It includes:

- **cert-manager**: The core cert-manager controller
- **CRDs**: Custom Resource Definitions for Issuer, ClusterIssuer, and Certificate resources
- **Gateway API Support**: Enabled for automatic certificate provisioning with Gateway API

## Prerequisites

- Kubernetes 1.22+
- Helm 3.0+

## Installation

### Install the chart

```bash
helm install cert-manager oci://ghcr.io/nicklasfrahm/charts/cert-manager -n cert-manager --create-namespace
```

### Install with custom values

```bash
helm install cert-manager oci://ghcr.io/nicklasfrahm/charts/cert-manager -n cert-manager --create-namespace -f my-values.yaml
```

## Configuration

| Parameter                   | Description                    | Default |
| --------------------------- | ------------------------------ | ------- |
| `cert-manager.enabled`      | Enable cert-manager deployment | `true`  |
| `cert-manager.crds.enabled` | Install cert-manager CRDs      | `true`  |

## Features

- Automatic TLS certificate issuance and renewal
- Support for multiple certificate authorities (Let's Encrypt, HashiCorp Vault, etc.)
- Gateway API integration for seamless certificate management
- Custom Resource Definitions for declarative certificate management

## Next Steps

After installing cert-manager, you can:

1. Create certificate issuers (Issuer or ClusterIssuer resources)
2. Request certificates using Certificate resources
3. Use annotations for automatic certificate provisioning with Ingress or Gateway resources

For more information, visit the [cert-manager documentation](https://cert-manager.io/docs/).
