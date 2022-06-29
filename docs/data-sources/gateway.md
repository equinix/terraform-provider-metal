---
page_title: "Equinix Metal: Metal Gateway"
subcategory: ""
description: |-
  Retrieve Equinix Metal Gateways
---

# metal_gateway (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_gateway`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_gateway) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_gateway`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this datasource to retrieve Metal Gateway resources in Equinix Metal.

~> VRF features are not generally available. The interfaces related to VRF resources may change ahead of general availability.

## Example Usage

```hcl
# Create Metal Gateway for a VLAN with a private IPv4 block with 8 IP addresses

resource "metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

data "metal_gateway" "test" {
  gateway_id = local.gateway_id
}
```

## Argument Reference

* `gateway_id` - (Required) UUID of the metal gateway resource to retrieve

## Attributes Reference

* `project_id` - UUID of the project where the gateway is scoped to
* `vlan_id` - UUID of the VLAN where the gateway is scoped to
* `ip_reservation_id` - UUID of IP reservation block bound to the gateway
* `private_ipv4_subnet_size` - Size of the private IPv4 subnet bound to this metal gateway, one of (8, 16, 32, 64, 128)`
* `state` - Status of the gateway resource
* `vrf_id` - UUID of the VRF associated with the IP Reservation.
