---
page_title: "Equinix Metal: connection"
subcategory: ""
description: |-
  Retrieve Equinix Fabric Connection
---

# metal\_connection

Use this data source to retrieve a [connection resource](https://metal.equinix.com/developers/docs/networking/fabric/)

~> Equinix Metal connection with service_token_type `a_side` is not generally available and may not be enabled yet for your organization.

## Example Usage

```hcl
data "metal_connection" "example" {
  connection_id = "4347e805-eb46-4699-9eb9-5c116e6a017d"
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (Required) ID of the connection resource

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Name of the connection resource
* `metro` - Slug of a metro to which the connection belongs
* `facility` - Slug of a facility to which the connection belongs
* `redundancy` - Connection redundancy, reduntant or primary
* `type` - Connection type, dedicated or shared
* `project_id` - ID of project to which the connection belongs
* `speed` - Connection speed, one of 50Mbps, 200Mbps, 500Mbps, 1Gbps, 2Gbps, 5Gbps, 10Gbps
* `description` - Description of the connection resource
* `mode` - Mode for connections in IBX facilities with the dedicated type - standard or tunnel
* `tags` - String list of tags
* `vlans` - Attached VLANs. Only available in shared connection. One vlan for Primary/Single connection and two vlans for Redundant connection
* `service_token_type` - Type of service token, a_side or z_side. One available in shared connection.
* `organization_id` - ID of the organization where the connection is scoped to
* `status` - Status of the connection resource
* `service_tokens` - List of connection service tokens with attributes
  * `id` - UUID of the service token required to configure the connection in the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard).
  * `expires_at` - Expiration date of the service token
  * `max_allowed_speed` - Maximum allowed speed for the service token, string like in the `speed` attribute
  * `type` - Token type, `a_side` or `z_side`
  * `role` - Token role, `primary` or `secondary`
* `ports` - List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)
  * `name` - Port name
  * `id` - Port UUID
  * `role` - Port role - primary or secondary
  * `speed` - Port speed in bits per second
  * `status` - Port status
  * `link_status` - Port link status
  * `virtual_circuit_ids` - List of IDs of virtual cicruits attached to this port
* `token` - (Deprecated) Token to configure the connection in the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)
