---
page_title: "Equinix Metal: metal_plans"
subcategory: ""
description: |-
  Provides an Equinix Metal plans datasource. This can be used to find plans that meet a filter criteria.
---

# metal_plans (Data Source)

!> **PROVIDER DEPRECATED:** Equinix Metal Provider is now Deprecated. Please consider using [`equinix_metal_plans`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/data-sources/equinix_metal_plans) data source from the [Equinix provider](https://registry.terraform.io/providers/equinix/equinix/latest/docs) instead of `metal_plans`. [See the Metal provider section for more details](../index.md#equinix-metal-provider) on the new provider and available migration guides.

Provides an Equinix Metal plans datasource. This can be used to find plans that meet a filter criteria.

## Example Usage

```hcl
# Following example will select device plans which are under 2.5$ per hour, are available in metro 'da' (Dallas)
# OR 'sv' (Sillicon Valley) and sort it by the hourly price ascending.
data "metal_plans" "example" {
    sort {
        attribute = "pricing_hour"
        direction = "asc"
    }
    filter {
        attribute = "pricing_hour"
        values    = [2.5]
        match_by  = "less_than"
    }
    filter {
        attribute = "available_in_metros"
        values    = ["da", "sv"]
    }
}

output "plans" {
    value = data.metal_plans.example.plans
}
```

```hcl
# Following example will select device plans with class containing string 'large', are available in metro 'da' (Dallas)
# AND 'sv' (Sillicon Valley), are elegible for spot_market deployments.
data "metal_plans" "example" {
    filter {
        attribute = "class"
        values    = ["large"]
        match_by  = "substring"
    }
    filter {
        attribute = "deployment_types"
        values    = ["spot_market"]
    }
    filter {
        attribute = "available_in_metros"
        values    = ["da", "sv"]
        all       = true
    }
}

output "plans" {
    value = data.metal_plans.example.plans
}
```

### Ignoring Changes to Plans/Facilities/Metro

Preserve deployed plan, facility and metro when updating a plan. As described in the [`data-resource-behavior`](https://www.terraform.io/language/data-sources#data-resource-behavior), terraform reads data resources during the planning phase in both the terraform plan and terraform apply commands. If the result of the data source has changed since last apply, this will involve changes to other resources where there is a reference to their attributes.

In the case of `metal_plans`, it may happen that a plan is no longer available in a facility/metro because there is no stock at that time, and then returned plans list will be different from last `apply`. Therefore, if a resource such as a `metal_device` uses the output of this data source to select a plan or facility/metro, the plan will report that `metal_device` will be destroyed and recreated.

To prevent that you can take advantage of the [`lifecycle ignore_changes`](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changes) feature like defined in the example below.

```hcl
# Following example will use metal_plans to select the cheapest plan available in metro 'sv' (Sillicon Valley)
data "metal_plans" "example" {
    sort {
        attribute = "pricing_hour"
        direction = "asc"
    }
    filter {
        attribute = "name"
        values    = ["c3.small.x86", "c3.medium.x86", "m3.large.x86"]
    }
    filter {
        attribute = "available_in_metros"
        values    = ["sv"]
    }
}

# This metal_device will use the first returned plan and the list of facilities
resource "metal_device" "example" {
  hostname         = "example"
  plan             = data.metal_plans.example.plans[0].name
  facilities       = data.metal_plans.example.plans[0].available_in
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = var.project_id

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `sort` - (Optional) One or more attribute/direction pairs on which to sort results. If multiple
sorts are provided, they will be applied in order
  - `attribute` - (Required) The attribute used to sort the results. Sort attributes are case-sensitive
  - `direction` - (Optional) Sort results in ascending or descending order. Strings are sorted in alphabetical order. One of: asc, desc
* `filter` - (Optional) One or more attribute/values pairs to filter off of
  - `attribute` - (Required) The attribute used to filter. Filter attributes are case-sensitive
  - `values` - (Required) The filter values. Filter values are case-sensitive. If you specify multiple values for a filter, the values are joined with an OR by default, and the request returns all results that match any of the specified values
  - `match_by` - (Optional) The type of comparison to apply. One of: `in` , `re`, `substring`, `less_than`, `less_than_or_equal`, `greater_than`, `greater_than_or_equal`. Default is `in`.
  - `all` - (Optional) If is set to true, the values are joined with an AND, and the requests returns only the results that match all specified values. Default is `false`.

All fields in the `plans` block defined below can be used as attribute for both `sort` and `filter` blocks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plans` - The ID of the facility
  - `id` - id of the plan
  - `name` - name of the plan
  - `slug`- plan slug
  - `description`- description of the plan
  - `line`- plan line, e.g. baremetal
  - `legacy`- flag showing if it's a legacy plan
  - `class`- plan class
  - `pricing_hour`- plan hourly price
  - `pricing_month`- plan monthly price
  - `deployment_types`- list of deployment types, e.g. on_demand, spot_market
  - `available_in`- list of facilities where the plan is available
  - `available_in_metros`- list of facilities where the plan is available
