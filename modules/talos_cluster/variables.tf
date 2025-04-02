# Passing the global configuration to the module allows for dependency injection.
variable "global_config" {
  description = "The global configuration."
  type = object({
    talos = object({
      version = string
    })
    kubernetes = object({
      version = string,
      oidc = object({
        issuer_url = string
        client_id  = string
      })
    })
    dns = object({
      zone = string
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

variable "config" {
  description = "The configuration of the cluster."
  type = object({
    metadata = object({
      name = string
    })
    spec = object({
      infrastructure = object({
        loadBalancer = optional(
          object({
            host = optional(string, "")
            port = optional(number, 6443)
          }),
          {
            host = ""
            port = 6443
          }
        )
        controlplanes = list(
          object({
            name = string
          })
        )
        workers = optional(
          list(
            object({
              name = string
            })
          ),
          []
        )
      })
    })
  })
}
