# gRPC gateway

This is a simple gRPC gateway that can be used to expose multiple gRPC services as a GraphQL server.

![arch](./assets//arch.excalidraw.png)

## Pre-requisites

- [Task](https://taskfile.dev/#/installation)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)

## Installation

To install the dependencies, run the following command:

```bash
docker-compose -f "deployment/docker-compose.yml" up -d
```

## Getting started

To start the gateway, run the following command:

```bash
task deps
task gateway
```

## Benchmark

To run the benchmark, run the following command:

```bash
task bench
```

## Pyroscope

Visiting [http://localhost:4040](http://localhost:4040) will show the Pyroscope dashboard.

## Uptrace

Visiting [http://localhost:14318/](http://localhost:14318/) will show the Uptrace dashboard.
