# Passing the global configuration to the module allows for dependency injection.
variable "global_config" {
  description = "The global configuration."
  type = object({
    talos = object({
      version = string
    })
  })
}

# Passing in all machines to the module allows for dependency injection.
variable "machines" {
  description = "A list of all available machines."
  type = list(
    object({
      metadata = object({
        name = string
      })
      spec = object({
        // TODO: Add more fields as needed.
      })
    })
  )
}

# The configuration of the cluster.
variable "config" {
  description = "The configuration of the cluster."
  type = object({
    metadata = object({
      name = string
    })
    spec = object({
      infrastructure = object({
        loadBalancer = object({
          host = optional(string, "")
          port = optional(number, 6443)
        })
        controlplanes = list(
          object({
            name = string
          })
        )
        workers = list(
          object({
            name = string
          })
        )
      })
    })
  })
}
