variable "config" {
  description = "The manifest of the Kubernetes cluster."
  type = object({
    metadata = object({
      name = string
    })
    spec = object({
      cluster = object({
        allowSchedulingOnControlPlanes = bool
      })
      infrastructure = object({
        controlPlane = object({
          machinePools = list(object({
            name = string
          }))
        })
      })
    })
  })
}

variable "domain" {
  description = "A fully-qualified domain name."
  type = string
}
