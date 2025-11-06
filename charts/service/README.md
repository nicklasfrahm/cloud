# Service 🚀

This chart is a general purpose chart to deploy a service to Kubernetes. It can be used to deploy webservices, queue workers and other workloads that are built to be cloud-native.

## Features ✨

- **Automatic port injection via `PORT` environment variable 🔌**  
  By default, the chart will set the `PORT` environment variable to set your applications internal port. If your application does not support this, you can configure the port via the `service.containerPort` value.

- **TLS via Gateway API 🔒**  
  Automatically configures TLS termination using Gateway API when `expose.enabled` is set to `true`. You can specify the full host name via `expose.host`, let the chart generate the host name automatically based on the cluster domain and release name, following the pattern `<release-name>.<cluster>.nicklasfrahm.dev` or only customize the subdomain via `expose.hostPrefix`.

- **Horizontal Pod Autoscaling 📈**  
  By default, a Horizontal Pod Autoscaler (HPA) is created to automatically scale your service based on 90% memory utilization and 100% CPU utilization. You can customize the HPA settings or disable it entirely via the `autoscaling.enabled` values.

- **Resource requests and limits ⚡**  
  By default, the chart will set CPU requests and memory requests that match the memory limits to prevent throttling and memory overcommit. You can customize these settings via the `resources` value.

- **Liveness, Readiness and Startup probes 🩺**  
  The chart expects your application to expose `/ready`, `/live` and `/startup` endpoints for health checks. You can customize the probe settings or disable them entirely via the `probes` value.

- **Pod Disruption Budgets 💥**  
  The chart will create a Pod Disruption Budget (PDB) to ensure that a minimum number of pods are always available during voluntary disruptions. You can customize the PDB settings via the `reliability.maxUnavailable` value, which defaults to `50%`.
