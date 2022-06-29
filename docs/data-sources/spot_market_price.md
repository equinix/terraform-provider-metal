---
page_title: "Equinix Metal: spot_market_price"
subcategory: ""
description: |-
  Get an Equinix Metal Spot Market Price
---

# metal_operating_system (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_operating_system`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_operating_system) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_operating_system`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this data source to get Equinix Metal Spot Market Price for a plan.

## Example Usage

Lookup by facility:

```hcl
data "metal_spot_market_price" "example" {
  facility = "ny5"
  plan     = "c3.small.x86"
}
```

Lookup by metro:

```hcl
data "metal_spot_market_price" "example" {
  metro    = "sv"
  plan     = "c3.small.x86"
}
```

## Argument Reference

* `plan` - (Required) Name of the plan.
* `facility` - (Optional) Name of the facility.
* `metro` - (Optional) Name of the metro.

## Attributes Reference

* `price` - Current spot market price for given plan in given facility.
