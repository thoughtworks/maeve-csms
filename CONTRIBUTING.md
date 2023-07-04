# How to contribute to MaEVe
Thank you for taking the time to read this guide. We are happy you wish to contribute to this project :partying_face: !

## Requirements
Please take some time to familiarise with the components in MaEVe. These are in the [README's documentation section](./README.md/#documentation).

You will need to have Go 1.20 installed in your development machine. Instructions are available on the [Go website](https://go.dev/doc/install)

MaEVe components can be run through a [docker-compose file](./docker-compose.yml), hence you will require Docker on your machine. 
Additionally, you will require the [compose plugin](https://docs.docker.com/compose/install) to run the containers listed in the compose file.

### What dependencies are provisioned through docker compose?
Currently, the [docker-compose file](./docker-compose.yml) provisions the following 
in addition to the [Manager](./docs/manager.md) and [Gateway](./docs/gateway.md) components:

- #### MQTT v5 broker
    This is used to decouple the stateful connections required by OCPP from the stateless message handlers. The system currently 
connects to the MQTT broker anonymously.

- #### Redis
    The system uses [Redis](https://redis.io/) as a storage engine for transaction details. The manager connects to the store anonymously.

- #### Load balancer
    A load balancer implemented through [envoyproxy](https://www.envoyproxy.io/) to load balance calls coming from the charge stations and allowing the system to
perform better under high demand.

## Build projects locally
Whilst the components can be run through Docker, if you wish to run locally without containerising them, you can run the following command:
```shell
$ go build ./...
```
_You will need to run this from the correct directory (e.g. ./gateway or ./manager)_

## Testing
Each Go project contains a set of tests that are run upon each push reaches the remote repository. 
These tests also run as part of git hooks before committing to the local repository.

If you wish to run the tests manually, you can run the following:
```shell
$ go test ./...
```

## How to report a bug and track issues
_Come back to this later, we are currently working on it._

## How to submit changes
_Come back to this later, we are currently working on it._

## Coding style guideline
_Come back to this later, we are currently working on it._

## Contributors
_Come back to this later, we are currently working on it._