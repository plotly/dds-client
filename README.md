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

#### `postgres:list`

List all postgres services

```shell
dds-client postgres:list
```

#### `postgres:exists`

Check if a postgres service exists

```shell
dds-client postgres:exists --name dopsa
```

#### `postgres:create`

Create a postgres service

```shell
dds-client postgres:create --name dopsa
```

#### `postgres:delete`

Delete a postgres service

```shell
dds-client postgres:delete --name dopsa
```

#### `redis:list`

List all redis services

```shell
dds-client redis:list
```

#### `redis:exists`

Check if a redis service exists

```shell
dds-client redis:exists --name dopsa
```

#### `redis:create`

Create a redis service

```shell
dds-client redis:create --name dopsa
```

#### `redis:delete`

Delete a redis service

```shell
dds-client redis:delete --name dopsa
```
