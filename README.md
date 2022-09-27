# Terraform Provider Honeybadger

Run the following command to build the provider

```shell
go build -o terraform-provider-honeybadger
```


## How do I use this?

You can download and run make install to compile and install locally.  This is great for testing, check out the examples/main.tf file for an example of testing stuff out.

You can also just specify it like this in your project

```
terraform {
  required_providers {
    honeybadger = {
      version = "~> 1.0.0"
      source  = "sequra.com/providers/honeybadger"
    }
  }
}
provider "honeybadger" {
  api_key = "<INTRODUCE_API_KEY>"
}
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, spin up the mock server

```shell
cd docker_compose
docker-compose up -d
```


Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples
terraform init && terraform apply
```
