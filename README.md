# Messages

[![Docker](https://github.com/skamenetskiy/messages/actions/workflows/docker.yml/badge.svg)](https://github.com/skamenetskiy/messages/actions/workflows/docker.yml)

Messages service.

## Features

- 🚀 Extremely simple API.
- 🚀 Included dynamic sharding.
- 🚀 Included mentions.
- 🚀 Both GRPC and HTTP protocols.

## Server configuration

The server can be configured with the following environment variables

| Variable name | Default | Description         |
|---------------|---------|---------------------|
| GRPC_HOST     | 0.0.0.0 | GRPC host to listen |
| GRPC_PORT     | 50051   | GRPC port to listen |
| HTTP_HOST     | 0.0.0.0 | HTTP host to listen |
| HTTP_PORT     | 8080    | HTTP port to listen |

## Database configuration

To configure the database simply add `db.yml` to the container (see [db.sample.yml](db.sample.yml)).

### Shard configuration

| Property | Type      | Description                                   |
|----------|-----------|-----------------------------------------------|
| id       | `uint16`  | Unique shard ID                               |
| writable | `boolean` | Defines if the shard is writable or read-only |
| dsn      | `string`  | Database conneciton string                    |