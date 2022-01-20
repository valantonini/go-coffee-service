# go-coffee-service

[![CI](https://github.com/valantonini/go-coffee-service/actions/workflows/makefile.yml/badge.svg)](https://github.com/valantonini/go-coffee-service/actions/workflows/makefile.yml)

Playing with patterns and libraries for distributed systems in the Go ecosystem.

Based on the [Hashicorp consul tutorial](https://learn.hashicorp.com/tutorials/consul/kubernetes-extract-microservice?in=consul/microservices)

## caveats

- i am new to go and many things aren't in their final shape
- striving to use as few dependencies as possible to understand the problems libraries and frameworks are addressing 
- secret management will not be production ready
- many variables are inline for the moment and need to be extracted out to env


## env

- go 1.17
- docker 20.10.11
- make

```bash
docker network create -d bridge coffee-service-network
```

### optional env

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

## run

```bash
# run the integration tests
make integration

# run the coffee-service
docker compose up
```
