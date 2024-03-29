---
page_title: "Equinix Metal: precreated_port"
subcategory: ""
description: |-
  Fetch device ports
---

# metal_port (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_port`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_port) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_port`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this data source to read ports of existing devices. You can read port by either its UUID, or by a device UUID and port name.

## Example Usage

Create a device and read it's eth0 port to the datasource.

```hcl
locals {
  project_id = "<UUID_of_your_project>"
}

resource "metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  facilities       = ["sv15"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

data "metal_port" "test" {
    device_id = metal_device.test.id
    name      = "eth0"
}
```

## Argument Reference

* `id` - (Required) ID of the port to read, conflicts with device_id.
* `device_id` - (Required) 
* `name` - (Required) Whether to look for public or private block.

## Attributes Reference

* `network_type` - One of layer2-bonded, layer2-individual, layer3, hybrid, hybrid-bonded
* `type` - Type is either "NetworkBondPort" for bond ports or "NetworkPort" for bondable ethernet ports
* `mac` - MAC address of the port
* `bond_id` - UUID of the bond port"
* `bond_name` - Name of the bond port
* `bonded` - Flag indicating whether the port is bonded
* `disbond_supported` - Flag indicating whether the port can be removed from a bond
* `native_vlan_id` - UUID of native VLAN of the port
* `vlan_ids` - UUIDs of attached VLANs
* `vxlan_ids` - VXLAN ids of attached VLANs

