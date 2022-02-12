# go-coffee-service

[![CI](https://github.com/valantonini/go-coffee-service/actions/workflows/makefile.yml/badge.svg)](https://github.com/valantonini/go-coffee-service/actions/workflows/makefile.yml)

Playing with patterns and libraries for distributed systems in the Go ecosystem.

Based on the [Hashicorp consul tutorial](https://learn.hashicorp.com/tutorials/consul/kubernetes-extract-microservice?in=consul/microservices) and [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

## caveats

- i am new to go and many things aren't in their final shape
- striving to use as few dependencies as possible to understand the problems libraries and frameworks are addressing 
- secret management will not be production ready
- many variables are inline for the moment and need to be extracted out to env
- still to do:
  - durable messaging
  - message retry
  - outbox pattern
  - implement a db

## env

- go 1.17
- docker 20.10.11
- make

```bash
docker network create -d bridge coffee-service-network
```

## run

```bash
# run the coffee-services
docker compose up

# run unit and integration tests
make test_all
```


## optional tooling

nats realtime monitoring

```bash
go get github.com/nats-io/nats-top

nats-top
```

nats command line tool (osx)
```bash
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```
