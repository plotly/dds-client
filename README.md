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
export DASH_API_KEY="SOME_KEY"
export DASH_API_USER="your-username"
export DASH_API_URL="https://dash.local/Manager/graphql"

dds-client -h
```


### Commands

#### `apps:list`

List all apps

```shell
dds-client apps:list
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
