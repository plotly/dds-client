# dds-client

A golang client for the dds GraphQL API

## Requirements

- Go 1.13+

## Building

```shell
go build -o dds-client
```

## Usage

> You can retrieve the `DASH_API_KEY` from the Dash Deployment Server Management UI.

```shell
export DASH_ENTERPRISE_API_KEY="SOME_KEY"
export DASH_ENTERPRISE_URL="https://dash.local"
export DASH_ENTERPRISE_USERNAME="your-username"

dds-client -h
```


### Commands

#### `apps:list`

List all apps

```shell
dds-client apps:list
```

#### `apps:exists`

Check if an app exists

```shell
dds-client apps:exists --name dopsa
```

#### `apps:create`

Create an app

```shell
dds-client apps:create --name dopsa
```

#### `apps:delete`

Delete an app

```shell
dds-client apps:delete --name dopsa
```
