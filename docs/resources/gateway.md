---
page_title: "Equinix Metal: Metal Gateway"
subcategory: ""
description: |-
  Create Equinix Metal Gateways
---

# metal_gateway (Resource)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_gateway`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_metal_gateway) resource from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_gateway`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this resource to create Metal Gateway resources in Equinix Metal.

~> VRF features are not generally available. The interfaces related to VRF resources may change ahead of general availability.

## Example Usage

```hcl
# Create Metal Gateway for a VLAN with a private IPv4 block with 8 IP addresses

resource "metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "metal_gateway" "test" {
  project_id               = local.project_id
  vlan_id                  = metal_vlan.test.id
  private_ipv4_subnet_size = 8
}
```

```hcl
# Create Metal Gateway for a VLAN and reserved IP address block

resource "metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "metal_reserved_ip_block" "test" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 2
}

resource "metal_gateway" "test" {
  project_id        = local.project_id
  vlan_id           = metal_vlan.test.id
  ip_reservation_id = metal_reserved_ip_block.test.id
}
```

## Argument Reference

* `project_id` - (Required) UUID of the project where the gateway is scoped to.
* `vlan_id` - (Required) UUID of the VLAN where the gateway is scoped to.
* `ip_reservation_id` - (Optional) UUID of Public or VRF IP Reservation to associate with the gateway, the reservation must be in the same metro as the VLAN, conflicts with `private_ipv4_subnet_size`.
* `private_ipv4_subnet_size` - (Optional) Size of the private IPv4 subnet to create for this metal gateway, must be one of (8, 16, 32, 64, 128), conflicts with `ip_reservation_id`.

## Attributes Reference

* `state` - Status of the gateway resource.
* `vrf_id` - UUID of the VRF associated with the IP Reservation.
