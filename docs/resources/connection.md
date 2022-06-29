---
page_title: "Equinix Metal: connection"
subcategory: ""
description: |-
  Request/Create Equinix Fabric Connection
---

# metal_connection (Resource)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_connection`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_metal_connection) resource from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_connection`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this resource to request of create an Interconnection from [Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/)

## Example Usage

```hcl

resource "metal_connection" "test" {
	name               = "My Interconnection"
	project_id         = metal_project.test.id
	type               = "shared"
	redundancy         = "redundant"
	metro              = "sv"
	speed              = "50Mbps"
	service_token_type = "a_side"
}
```

## Argument Reference

* `name` - (Required) Name of the connection resource
* `metro` - (Optional) Metro where the connection will be created
* `facility` - (Optional) Facility where the connection will be created
* `redundancy` - (Required) Connection redundancy - redundant or primary
* `type` - (Required) Connection type - dedicated or shared
* `project_id` - (Required) ID of the project where the connection is scoped to, must be set for shared connection
* `speed` - (Required) Connection speed - one of 50Mbps, 200Mbps, 500Mbps, 1Gbps, 2Gbps, 5Gbps, 10Gbps
* `description` - (Optional) Description for the connection resource
* `mode` - (Optional) Mode for connections in IBX facilities with the dedicated type - standard or tunnel. Default is standard
* `tags` - (Optional) String list of tags
* `vlans` - (Optional) Only used with shared connection. Vlans to attach. Pass one vlan for Primary/Single connection and two vlans for Redundant connection
* `service_token_type` - (Optional) Only used with shared connection. Type of service token to use for the connection, a_side or z_side

## Attributes Reference

* `organization_id` - ID of the organization where the connection is scoped to
* `status` - Status of the connection resource
* `ports` - List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`). Schema of port is described in documentation of the [metal_connection datasource](../data-sources/connection.md).
* `service_tokens` - List of connection service tokens with attributes. Scehma of service_token is described in documentation of the [metal_connection datasource](../data-sources/connection.md).
