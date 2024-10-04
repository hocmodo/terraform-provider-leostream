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

## More enhanced way for importing resources

Use https://gitlab.hocmodo.nl/community/leostream-admin-cli to pull the data from the Leostream API and get the id's

## AWS Pools

```shell
for pool_id in `leostream-admin-cli pool list --json | jq '.[].id'`
do
  terraform import leostream_aws_pool $pool_id
done
```

## Basic pools

```shell
for pool_id in `leostream-admin-cli pool list --json | jq '.[].id'`
do
    terraform import leostream_basic_pool $pool_id
done
```


## Centers

```shell
for center_id in `leostream-admin-cli center list --json | jq '.[].id'`
do
  terraform import leostream_center $center_id
done
```

## Gateways

```shell
for gateway_id in `leostream-admin-cli gateway list --json | jq '.[].id'`
do
  terraform import leostream_gateway $gateway_id
done
```
