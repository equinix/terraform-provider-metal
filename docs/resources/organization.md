---
page_title: "Equinix Metal: metal_organization"
subcategory: ""
description: |-
  Provides an Equinix Metal Organization resource.
---

# metal\_organization

Provides a resource to manage organization resource in Equinix Metal.

## Example Usage

```hcl
# Create a new Project
resource "metal_organization" "tf_organization_1" {
  name        = "foobar"
  description = "quux"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Organization
* `address` - (Required) An object that has the address information. See [Address](#address)
below for more details.
* `description` - (Optional) Description string
* `website` - (Optional) Website link
* `twitter` - (Optional) Twitter handle
* `logo` - (Optional) Logo URL

### Address

The `address` block contains:

* `address` - (Required) Postal address.
* `city` - (Required) City name.
* `country` - (Required) Two letter country code (ISO 3166-1 alpha-2), e.g. US.
* `zip_code` - (Required) Zip Code.
* `state` - (Optional) State name.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the organization
* `name` - The name of the Organization
* `description` - Description string
* `website` - Website link
* `twitter` - Twitter handle
* `logo` - Logo URL

## Import

This resource can be imported using an existing organization ID:

```sh
terraform import metal_organization {existing_organization_id}
```
