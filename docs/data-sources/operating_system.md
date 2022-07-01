---
page_title: "Equinix Metal: operating_system"
subcategory: ""
description: |-
  Get an Equinix Metal operating system image
---

# metal_operating_system (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_operating_system`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_operating_system) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_operating_system`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this data source to get Equinix Metal Operating System image.

## Example Usage

```hcl
data "metal_operating_system" "example" {
  distro           = "ubuntu"
  version          = "20.04"
  provisionable_on = "c3.medium.x86"
}

resource "metal_device" "server" {
  hostname         = "tf.ubuntu"
  plan             = "c3.medium.x86"
  facilities       = ["ny5"]
  operating_system = data.metal_operating_system.example.id
  billing_cycle    = "hourly"
  project_id       = local.project_id
}
```

## Argument Reference

* `distro` - (Optional) Name of the OS distribution.
* `name` - (Optional) Name or part of the name of the distribution. Case insensitive.
* `provisionable_on` - (Optional) Plan name.
* `version` - (Optional) Version of the distribution

## Attributes Reference

* `id` - Operating system slug
* `slug` - Operating system slug (same as `id`)
