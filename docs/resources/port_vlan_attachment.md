---
page_title: "Equinix Metal: metal_port_vlan_attachment"
subcategory: ""
description: |-
  Provides a Resource for Attaching VLANs to Device Ports
---

# metal_port_vlan_attachment (Resource)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_port_vlan_attachment`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_metal_port_vlan_attachment) resource from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_port_vlan_attachment`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Provides a resource to attach device ports to VLANs.

Device and VLAN must be in the same facility.

If you need this resource to add the port back to bond on removal, set `force_bond = true`.

To learn more about Layer 2 networking in Equinix Metal, refer to

* <https://metal.equinix.com/developers/docs/networking/layer2/>
* <https://metal.equinix.com/developers/docs/networking/layer2-configs/>

## Example Usage

### Hybrid network type

```hcl
resource "metal_vlan" "test" {
  description = "VLAN in New Jersey"
  facility    = "ny5"
  project_id  = local.project_id
}

resource "metal_device" "test" {
  hostname         = "test"
  plan             = "c3.small.x86"
  facilities       = ["ny5"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type      = "hybrid"
}

resource "metal_port_vlan_attachment" "test" {
  device_id = metal_device_network_type.test.id
  port_name = "eth1"
  vlan_vnid = metal_vlan.test.vxlan
}

```

### Layer 2 network

```hcl
resource "metal_device" "test" {
  hostname         = "test"
  plan             = "c3.small.x86"
  facilities       = ["ny5"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type      = "layer2-individual"
}

resource "metal_vlan" "test1" {
  description = "VLAN in New Jersey"
  facility    = "ny5"
  project_id  = local.project_id
}

resource "metal_vlan" "test2" {
  description = "VLAN in New Jersey"
  facility    = "ny5"
  project_id  = local.project_id
}

resource "metal_port_vlan_attachment" "test1" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = metal_vlan.test1.vxlan
  port_name = "eth1"
}

resource "metal_port_vlan_attachment" "test2" {
  device_id  = metal_device_network_type.test.id
  vlan_vnid  = metal_vlan.test2.vxlan
  port_name  = "eth1"
  native     = true
  depends_on = ["metal_port_vlan_attachment.test1"]
}
```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) ID of device to be assigned to the VLAN
* `port_name` - (Required) Name of network port to be assigned to the VLAN
* `force_bond` - Add port back to the bond when this resource is removed. Default is false.
* `vlan_vnid` - VXLAN Network Identifier, integer
* `native` - (Optional) Mark this VLAN a native VLAN on the port. This can be used only if this assignment assigns second or further VLAN to the port. To ensure that this attachment is not first on a port, you can use `depends_on` pointing to another metal_port_vlan_attachment, just like in the layer2-individual example above.

## Attribute Referece

* `id` - UUID of device port used in the assignment
* `vlan_id` - UUID of VLAN API resource
* `port_id` - UUID of device port
