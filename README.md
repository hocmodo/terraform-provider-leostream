# Terraform Provider Leostream

**TL;DR** This Leostream provider acts as a bridge between Terraform and the Leostream REST API.

 Leostream is a Remote Desktop Access platform for specialized display protocols. It features a broker component which has a lot of configuration options. These are configurable using the GUI...

 _or_ the REST API 🤩

 More info on the product on https://leostream.com/.

## Develop and build provider

To start developing the provider, you need to have a working Go environment. You can find the installation instructions [here](https://golang.org/doc/install). Also you need to have Terraform installed, you can find the installation instructions [at the hashicorp site.](https://learn.hashicorp.com/tutorials/terraform/install-cli).

You need to read up on the [Terraform Plugin SDK](https://developer.hashicorp.com/terraform/plugin) to understand how to build a provider. The are two versions of the SDK, the old one and the new one. This provider is built using the new one (called protocol version6 or 'Terraform Plugin Framework'). So don't use examples from Terraform Plugin SDKv2 (protocol version5) to build this provider.

There is an excellent [tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider) on how to build a provider.

Also, for handling nested objects I have studied gmichel's [Adguard provider](https://github.com/gmichels/terraform-provider-adguard), because I could not find many other (clear) examples.

If you want to start development work on the provider, you have to make sure you have Terrafom configured to use your code when you want to test it.

First, clone the repository to your `GOPATH`.

```shell
% git clone
```

Navigate to the directory and run the following command to install the dependencies.

```shell
% go mod tidy
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

Run the following command to build the provider, which will create the binary in `~/go/bin`.

```shell
% make install
```

### Create documentation

You can generate documentation from the source by using the accompanied tools:

```shell
cd tools; go generate ./...
```

### Debugging

For detailed instructions on debugging — including running the provider under Delve and intercepting HTTP traffic with mitmproxy — see [DEVELOPMENT.md](DEVELOPMENT.md).

## Run integration tests

There are also integration tests that can be run to verify the provider's functionality. To run the integration tests, you first need a running Leostream instance. Currently there is no way to run the tests without a running Leostream instance.
You can make sure the right data is available by using [the leostream-admin-cli tool](https://gitlab.com/hocmodo/leostream-admin-cli) to add sample data.

### Add a test center for testing the centers datasource

```shell
% leostream-admin-cli center create --access_key your_aws_access_key --name "Test AWS center" --region eu-west-1 --type amazon --secret_key your_aws_secret_key
{
        "stored_data": {
                "id": 52,
                "status": 5,
                "center_definition": {
                        "name": "Test AWS center",
                        "type": "amazon",
                        "type_label": "Amazon Web Services",
                        "wait_inst_status": 1,
                        "wait_sys_status": 1
                },
                "status_label": "Scanning"
        }
}%
```

### Get the center id

```shell
% leostream-admin-cli center list
+-----------+----------------------+--------+--------+---------------------+
| CENTER ID | CENTER NAME          | ONLINE | TYPE   | TYPE LABEL          |
+-----------+----------------------+--------+--------+---------------------+
|        51 | aws-center-us-east-1 |      1 | amazon | Amazon Web Services |
|        52 | Test AWS center      |      1 | amazon | Amazon Web Services |
```

## Test example configuration from /examples

Navigate to the `examples` directory.

```shell
% cd examples/pick-one
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
% terraform init && terraform plan/apply -var-file="secret.tfvars"
```

or skip the terraform init if you have the dev_overrides in the terraformrc file.

## More enhanced way for importing resources

Use the [admin-cli](https://gitlab.com/hocmodo/leostream-admin-cli) to pull the data from the Leostream API and get the id's

## AWS Pools

```shell
% for pool_id in `leostream-admin-cli pool list --json | jq '.[].id'`
do
  terraform import leostream_aws_pool $pool_id
done
```

## Basic pools

```shell
% for pool_id in `leostream-admin-cli pool list --json | jq '.[].id'`
do
    terraform import leostream_basic_pool $pool_id
done
```

## Centers

```shell
% for center_id in `leostream-admin-cli center list --json | jq '.[].id'`
do
  terraform import leostream_center $center_id
done
```

## Gateways

```shell
% for gateway_id in `leostream-admin-cli gateway list --json | jq '.[].id'`
do
  terraform import leostream_gateway $gateway_id
done
```
