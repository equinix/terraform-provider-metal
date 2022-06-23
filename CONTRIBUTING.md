# Hello Contributors

Thx for your interest! We're so glad you're here.

Before continuing, please note that this project is deprecated and all new work should be started in the [Equinix Terraform Provider](https://github.com/equinix/terraform-provider-equinix). See the [README.md](README.md) for more details.

## Contributing
### Important Resources

- bugs: [https://github.com/equinix/terraform-provider-metal/issues](https://github.com/equinix/terraform-provider-metal/issues)
- features: [https://github.com/equinix/terraform-provider-metal/issues](https://github.com/equinix/terraform-provider-metal/issues)

### Code of Conduct

Available via [https://github.com/equinix/terraform-provider-metal/blob/main/.github/CODE_OF_CONDUCT.md](https://github.com/equinix/terraform-provider-metal/blob/main/.github/CODE_OF_CONDUCT.md)

### Environment Details

[https://github.com/equinix/terraform-provider-metal/blob/main/GNUmakefile](https://github.com/equinix/terraform-provider-metal/blob/main/GNUmakefile)

### How to Submit Change Requests

Please submit change requests and / or features via [Issues](https://github.com/equinix/terraform-provider-metal/issues). There's no guarantee it'll be changed, but you never know until you try. We'll try to add comments as soon as possible, though.

### How to Report a Bug

Bugs are problems in code, in the functionality of an application or in its UI design; you can submit them through [Issues/(https://github.com/equinix/terraform-provider-metal/issues) as well.

## Development
### Requirements

- [Terraform 0.12+](https://www.terraform.io/downloads.html) (for v1.0.0 of this provider and newer)
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

### Building the provider

Clone the repository, enter the provider directory, and build the provider.

```sh
git clone git@github.com:equinix/terraform-provider-metal
cd terraform-provider-metal
make build
```

### Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-metal
...
```

### Testing provider code

We have mostly acceptance tests in the provider. There's no point for you to run them all, but you should run the one covering the functionality which you change. The acceptance test run will cost you some money, so feel free to abstain. The acceptance test suite will be run for your PR during the review process.

To run an acceptance test, find the relevant test function in `*_test.go` (for example TestAccMetalDevice_Basic), and run it as

```sh
TF_ACC=1 go test -v -timeout=20m ./... -run=TestAccMetalDevice_Basic
```

If you want to see HTTP traffic, set `TF_LOG=DEBUG`, i.e.

```sh
TF_LOG=DEBUG TF_ACC=1 go test -v -timeout=20m ./... -run=TestAccMetalDevice_Basic
```

### Testing the provider with Terraform

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider) above), it can be incorporated within your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will then use the plugin from your development environment.

```sh
terraform init -plugin-dir $GOPATH/bin
```
