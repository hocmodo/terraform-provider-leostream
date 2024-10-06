# Terraform Provider Leostream


**TL;DR** This Leostream provider acts as a bridge between Terraform and the Leostream REST API.

 Leostream is a Remote Desktop Access platform for specialized display protocols. It features a broker component which has a lot of configuration options. These are configurable using the GUI...

 _or_ the REST API 🤩

 More info on the product on https://leostream.com/.

## Develop and build provider

If you want to start development work on the provider, you have to make sure you have Terrafom configured to use your code when you want to test it.

First, clone the repository to your `GOPATH`.

```shell
$ git clone
```

Navigate to the directory and run the following command to install the dependencies.

```shell
$ go mod tidy
```

Add a dev_overrides to the terraform configuration file  (typiclly $HOME/.terraformrc)  to point to the directory you checked out.

```shell
provider_installation {

  dev_overrides {
      "registry.terraform.io/hocmodo/leostream" = "/Path/to/home/dir/go/bin",
      "hashicorp/time" =  "/Path/to/home/dir/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```


Run the following command to build the provider

```shell
$ make install
```

## Test sample configuration


Navigate to the `examples` directory.

```shell
$ cd examples/pick-one
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform plan/apply -var-file="secret.tfvars"
```

or skip the terraform init if you have the dev_overrides in the terraformrc file.



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
