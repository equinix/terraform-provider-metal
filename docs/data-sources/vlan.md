---
page_title: "Equinix Metal: metal_vlan"
subcategory: ""
description: |-
  Provides an Equinix Metal Virtual Network datasource. This can be used to read the attributes of existing VLANs.
---

# metal_vlan (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_vlan`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_vlan) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_vlan`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Provides an Equinix Metal Virtual Network datasource. VLANs data sources can be
searched by VLAN UUID, or project UUID and vxlan number.

## Example Usage

Fetch a vlan by ID:

```hcl
resource "metal_vlan" "foovlan" {
        project_id = local.project_id
        metro = "sv"
        vxlan = 5
}

data "metal_vlan" "dsvlan" {
        vlan_id = metal_vlan.foovlan.id
}
```

Fetch a vlan by project ID, vxlan and metro

```hcl
resource "metal_vlan" "foovlan" {
        project_id = local.project_id
        metro = "sv"
        vxlan = 5
}

data "metal_vlan" "dsvlan" {
        project_id = local.project_id
        vxlan      = 5
        metro      = "sv"
}
```

## Argument Reference

The following arguments are supported:

* `vlan_id` - Metal UUID of the VLAN resource to look up
* `project_id` - UUID of parent project of the VLAN. Use together with the vxlan number and metro or facility
* `vxlan` - vxlan number of the VLAN to look up. Use together with the project_id and metro or facility
* `facility` - Facility where the VLAN is deployed
* `metro` - Metro where the VLAN is deployed

## Attributes Reference

The following attributes are exported, in addition to any unspecified arguments.

* `description` - Description text of the VLAN resource
* `assigned_devices_ids` - List of device ID to which this VLAN is assigned
