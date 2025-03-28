variable "region" {
  description = "The configuration of the region."
  type = object({
    metadata = object({
      name = string
    })
    spec = object({
      provider = string
      baremetal = object({
        controlplanes = list(object({
          name = string
        }))
      })
    })
  })
}
