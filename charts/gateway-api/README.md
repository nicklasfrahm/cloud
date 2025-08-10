# Gateway API Helm Chart

![Gateway API Logo](https://www.kubeblog.com/wp-content/uploads/2024/09/gateway-api-logo-1.png)

A Helm chart for installing Gateway API Custom Resource Definitions (CRDs) in your Kubernetes cluster.

## Overview

The Gateway API represents the next generation of Kubernetes Ingress, Load Balancing, and Service Mesh APIs. This Helm chart provides a simple way to install the Gateway API CRDs as a prerequisite for other Gateway API implementations.

## Features

- **Idempotent Installation**: Safe to run multiple times without removing existing CRDs
- **Version Control**: Specify exact Gateway API version to install
- **Channel Support**: Choose between standard and experimental CRDs
- **Resource Selection**: Install only the CRDs you need
- **Pre-install Hook**: Automatically installs CRDs before other resources

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installation

### Add the repository (if applicable)

```bash
# Add your helm repository here if published
helm repo add your-repo <repository-url>
helm repo update
```

### Install the chart

```bash
helm install gateway-api ./charts/gateway-api
```

### Install with custom values

```bash
helm install gateway-api ./charts/gateway-api -f custom-values.yaml
```

## Configuration

The following table lists the configurable parameters of the gateway-api chart and their default values.

| Parameter                              | Description                                  | Default               |
| -------------------------------------- | -------------------------------------------- | --------------------- |
| `gatewayCRDInstaller.enabled`          | Enable the installation of Gateway API CRDs  | `true`                |
| `gatewayCRDInstaller.image.repository` | Container image repository for CRD installer | `bitnami/kubectl`     |
| `gatewayCRDInstaller.image.tag`        | Container image tag for CRD installer        | `1.33.2-debian-12-r2` |
| `gatewayCRDInstaller.version`          | Version of Gateway API CRDs to install       | `1.2.0`               |
| `gatewayCRDInstaller.channel`          | Channel for CRDs (standard or experimental)  | `experimental`        |
| `gatewayCRDInstaller.resources`        | List of specific CRDs to install             | See values.yaml       |

### Available CRDs

The chart can install the following Gateway API CRDs:

- `gateway.networking.k8s.io_gatewayclasses` - Gateway class definitions
- `gateway.networking.k8s.io_gateways` - Gateway resources
- `gateway.networking.k8s.io_httproutes` - HTTP routing rules
- `gateway.networking.k8s.io_referencegrants` - Cross-namespace references
- `gateway.networking.k8s.io_grpcroutes` - gRPC routing rules (experimental)
- `gateway.networking.k8s.io_tlsroutes` - TLS routing rules (experimental)

### Example Custom Values

```yaml
gatewayCRDInstaller:
  enabled: true
  version: "1.2.0"
  channel: "standard" # Use standard channel instead of experimental
  image:
    repository: "bitnami/kubectl"
    tag: "1.33.2-debian-12-r2"
  resources:
    - gateway.networking.k8s.io_gatewayclasses
    - gateway.networking.k8s.io_gateways
    - gateway.networking.k8s.io_httproutes
    - gateway.networking.k8s.io_referencegrants
```

## How it Works

This chart uses Helm hooks to install Gateway API CRDs before any other resources:

1. **Pre-install Hook**: A Kubernetes Job runs before chart installation
2. **CRD Installation**: The job uses `kubectl` to apply Gateway API CRDs
3. **Idempotent**: Multiple runs won't conflict with existing CRDs
4. **Cleanup**: The job is automatically cleaned up after successful completion

## Upgrading

### To a new chart version

```bash
helm upgrade gateway-api ./charts/gateway-api
```

### To a new Gateway API version

Update the `gatewayCRDInstaller.version` value and upgrade:

```yaml
gatewayCRDInstaller:
  version: "1.3.0" # New version
```

```bash
helm upgrade gateway-api ./charts/gateway-api -f values.yaml
```

## Uninstallation

```bash
helm uninstall gateway-api
```

**Note**: This will not remove the CRDs themselves, as they may be used by other applications. To remove CRDs manually:

```bash
kubectl delete crd gatewayclasses.gateway.networking.k8s.io
kubectl delete crd gateways.gateway.networking.k8s.io
kubectl delete crd httproutes.gateway.networking.k8s.io
kubectl delete crd referencegrants.gateway.networking.k8s.io
kubectl delete crd grpcroutes.gateway.networking.k8s.io
kubectl delete crd tlsroutes.gateway.networking.k8s.io
```

## Troubleshooting

### CRDs not installed

Check the job status:

```bash
kubectl get jobs -n <namespace>
kubectl logs job/<release-name>-gateway-crd-installer -n <namespace>
```

### Permission issues

Ensure the service account has cluster-admin permissions to create CRDs.

### Version conflicts

If you have existing Gateway API CRDs with a different version, you may need to:

1. Check current CRD versions: `kubectl get crd -o custom-columns=NAME:.metadata.name,VERSION:.spec.versions[0].name | grep gateway`
2. Plan your upgrade strategy
3. Consider using `kubectl replace` instead of `kubectl apply` for major version changes

## Compatibility

| Gateway API Version | Kubernetes Version | Chart Version |
| ------------------- | ------------------ | ------------- |
| v1.2.0              | 1.19+              | 0.1.0         |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the chart
5. Submit a pull request

## Resources

- [Gateway API Documentation](https://gateway-api.sigs.k8s.io/)
- [Gateway API GitHub Repository](https://github.com/kubernetes-sigs/gateway-api)
- [Gateway API Installation Guide](https://gateway-api.sigs.k8s.io/guides/#installing-gateway-api)
- [Gateway API Versioning](https://gateway-api.sigs.k8s.io/concepts/versioning/)

## License

This chart is licensed under the same license as the Gateway API project.
