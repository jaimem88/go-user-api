# Go User API

[![Build status](https://badge.buildkite.com/49a408b9fd43b6628a55e541f5f7a73e4e548c37c1e7fb0491.svg)](https://buildkite.com/assembly-payments/banking-be-hiring-go-user-api)

## *DISCLAIMER*

This is repository is used by Assembly's banking Backend team for hiring and evaluation purposes. Any functionality described in this repository does not entirely reflect or show how the actual implementation of some services is done.

## Service definition

This is a gRPC service that provides functionality in relation to users.
The service definition is in the protocol buffer definition inside [./pkg/proto/user.proto](./pkg/proto/user.proto)

## Requirements

1. [Docker](https://docs.docker.com/install/)
2. Alternatively install [Go 1.13.4](https://golang.org/dl/) locally

## How to build, test and run

A [Makefile](Makefile) is available with the following recipes defined

### Build

```sh
    make build
```

### Test

Unit tests

```sh
    make test
```

End-to-end tests

```sh
    make test-e2e
```

Both unit and e2e tests

```sh
    make test-all
```

### Lint

```sh
    make lint
```

### Run

```sh
    make run
```

## GitHub Actions
