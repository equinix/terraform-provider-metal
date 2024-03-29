---
page_title: "Equinix Metal: metal_hardware_reservation"
subcategory: ""
description: |-
  Retrieve Equinix Metal Hardware Reservation
---

# metal_hardware_reservation (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_hardware_reservation`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_hardware_reservation) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_hardware_reservation`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Use this data source to retrieve a [hardware reservation resource from Equinix Metal](https://metal.equinix.com/developers/docs/deploy/reserved/).

You can look up hardware reservation by its ID or by ID of device which occupies it.

## Example Usage

```hcl
// lookup by ID
data "hardware_reservation" "example" {
  id     = "4347e805-eb46-4699-9eb9-5c116e6a0172"
}

// lookup by device ID
data "hardware_reservation" "example_by_device_id" {
  device_id     = "ff85aa58-c106-4624-8f1c-7c64554047ea"
}
```

## Argument Reference

* `id` - ID of the hardware reservation
* `device_id` - UUID of device occupying the reservation

## Attributes Reference

* `id` - ID of the hardware reservation to look up
* `short_id` - Reservation short ID
* `project_id` - UUID of project this reservation is scoped to
* `device_id` - UUID of device occupying the reservation
* `plan` - Plan type for the reservation
* `facility` - Plan type for the reservation
* `provisionable` - Flag indicating whether the reserved server is provisionable or not. Spare devices can't be provisioned unless they are activated first
* `spare` -  Flag indicating whether the Hardware Reservation is a spare. Spare Hardware Reservations are used when a Hardware Reservations requires service from Metal Equinix
* `switch_uuid` - Switch short ID, can be used to determine if two devices are connected to the same switch
