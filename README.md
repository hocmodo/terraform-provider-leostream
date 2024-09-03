# Terraform Provider Leostream

 This Leostream provider acts as a bridge between Terraform and the Leostream REST API.
## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-leostream
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then add the directory you checked out to the dev_overrides in the terraform configuration file.

Then, navigate to the `examples` directory.

```shell
$ cd examples/pick-one
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform plan/apply -var-file="secret.tfvars"
```
