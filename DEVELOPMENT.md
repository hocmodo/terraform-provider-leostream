# Development Guide

This guide covers common tasks and techniques useful when developing the Leostream Terraform provider.

## Debugging HTTP requests with mitmproxy

[mitmproxy](https://mitmproxy.org/) is a tool that lets you intercept and inspect HTTP/HTTPS traffic between the provider and the Leostream REST API. It comes with two interfaces:

- `mitmproxy` — interactive terminal UI
- `mitmweb` — browser-based UI

### Install

```shell
brew install mitmproxy
```

### Run

Start the proxy with the terminal UI:

```shell
mitmproxy --listen-port 8080
```

Or use the web UI (accessible at `http://localhost:8081`):

```shell
mitmweb --listen-port 8080
```

### Configure the provider to route traffic through the proxy

Set the standard proxy environment variables before running any `terraform` command:

```shell
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
```

### Trust the mitmproxy CA certificate

On first run, mitmproxy generates a self-signed CA at `~/.mitmproxy/mitmproxy-ca-cert.pem`. The provider's HTTP client needs to trust it, otherwise TLS verification will fail.

Point the Go TLS stack at the certificate:

```shell
export SSL_CERT_FILE=~/.mitmproxy/mitmproxy-ca-cert.pem
```

You can combine all three exports and then run Terraform:

```shell
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
export SSL_CERT_FILE=~/.mitmproxy/mitmproxy-ca-cert.pem

terraform apply
```

All HTTP(S) calls made by the provider will now appear in the mitmproxy UI.

Alternatively, if you want to skip SSL verification entirely during development, you can set:

```shell
export LEOSTREAM_API_SSL_INSECURE=True  # If you want to ignore SSL errors
```
See the [Leostream client README](https://gitlab.com/hocmodo/leostream-client-go)

> **Note:** Only use this in a local development environment. Never disable SSL verification in production.

---

## Running the provider in debug mode

The Terraform Plugin Framework supports attaching a debugger (e.g. [Delve](https://github.com/go-delve/delve)) to a running provider process. This lets you set breakpoints and step through provider code while Terraform is executing.

HashiCorp's official guide is at <https://developer.hashicorp.com/terraform/plugin/debugging>.

### 1. Install Delve

```shell
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 2. Start the provider under Delve

Run the provider binary with the `-debug` flag via `dlv`:

```shell
dlv exec ./terraform-provider-leostream -- -debug
```

Or, if you prefer to let Delve build the binary for you from source:

```shell
dlv debug . -- -debug
```

Delve will pause at startup and print something like:

```
API server listening at: 127.0.0.1:PORT
```

Type `continue` (or `c`) at the Delve prompt to let the provider start. The provider will then print a `TF_REATTACH_PROVIDERS` line to stdout, for example:

```
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable:

	TF_REATTACH_PROVIDERS='{"registry.terraform.io/hocmodo/leostream":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin123"}}}'
```

### 3. Run Terraform with the reattach variable

In a **separate terminal**, copy the exported value and run your Terraform command:

```shell
export TF_REATTACH_PROVIDERS='{"registry.terraform.io/hocmodo/leostream":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin123"}}}'

terraform apply
```

Terraform will connect to the already-running provider process instead of starting a new one, and Delve will hit any breakpoints you have set.

### 4. Without a debugger (log-based debugging)

If you only want the reattach behaviour without Delve, build and run the provider directly with the `-debug` flag:

```shell
go build -o terraform-provider-leostream .
./terraform-provider-leostream -debug
```

Then follow step 3 above. You can increase log verbosity by setting:

```shell
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log
```
