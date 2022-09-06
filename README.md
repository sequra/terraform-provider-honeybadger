# Terraform Provider Honeybadger

Run the following command to build the provider

```shell
go build -o terraform-provider-honeybadger
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
