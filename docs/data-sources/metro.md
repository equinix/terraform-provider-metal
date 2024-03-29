---
page_title: "Equinix Metal: metal_metro"
subcategory: ""
description: |-
  Provides an Equinix Metal metro datasource. This can be used to read metros.
---

# metal_metro (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_metro`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_metro) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_metro`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Provides an Equinix Metal metro datasource.

## Example Usage

```hcl
# Fetch a metro by code and show its ID

data "metal_metro" "sv" {
  code = "sv"
}

output "id" {
  value = data.metal_metro.sv.id
}
```


```hcl
# Verify that metro "sv" has capacity for provisioning 2 c3.small.x86 
  devices and 1 c3.medium.x86 device

data "metal_facility" "test" {
  code = "dc13"
  capacity {
    plan = "c3.small.x86"
    quantity = 2
  }
  capacity {
    plan = "c3.medium.x86"
    quantity = 1
  }
}

  ```

## Argument Reference

The following arguments are supported:

* `code` - The metro code

Metros can be looked up by `code`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the metro
* `code` - The code of the metro
* `country` - The country of the metro
* `name` - The name of the metro
* `capacity` - (Optional) Ensure that queried metro has capacity for specified number of given plans
  - `plan` - device plan to check
  - `quantity` - number of device to check
