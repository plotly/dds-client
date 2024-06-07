# dds-client

A golang client for the dds GraphQL API

<div align="center">
  <a href="https://dash.plotly.com/project-maintenance">
    <img src="https://dash.plotly.com/assets/images/maintained-by-plotly.png" width="400px" alt="Maintained by Plotly">
  </a>
</div>


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

#### `postgres:link`

Link a postgres service to an app

```shell
dds-client postgres:link --name dopsa --app dopsa
```

#### `postgres:unlink`

Unlink a postgres service from an app

```shell
dds-client postgres:unlink --name dopsa --app dopsa
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

#### `redis:link`

Link a redis service to an app

```shell
dds-client redis:link --name dopsa --app dopsa
```

#### `redis:unlink`

Unlink a redis service from an app

```shell
dds-client redis:unlink --name dopsa --app dopsa
```
