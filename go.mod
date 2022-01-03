module github.com/equinix/terraform-provider-metal

require (
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.9.0
	github.com/packethost/packngo v0.19.2-0.20210922152159-b073e9ef6568
)

go 1.16

replace github.com/packethost/packngo => github.com/t0mk/packngo v0.0.0-20220103111840-26c097341bbc
