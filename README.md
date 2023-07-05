[![Manager](https://github.com/twlabs/maeve-csms/workflows/Manager/badge.svg)](https://github.com/twlabs/maeve-csms/actions/workflows/manager.yml)
[![Gateway](https://github.com/twlabs/maeve-csms/workflows/Gateway/badge.svg)](https://github.com/twlabs/maeve-csms/actions/workflows/gateway.yml)

# MaEVe

MaEVe is an EV charge station management system (CSMS). It began life as a simple proof of concept for
implementing ISO-15118-2 Plug and Charge (PnC) functionality and remains a work in progress. It is hoped that over
time it will become more complete, but already provides a useful basis for experimentation.

The system currently integrates with [Hubject](https://hubject.stoplight.io/) for PnC functionality. 

## Table of Contents
- [Documentation](#documentation)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Documentation
MaEVe is implemented in Go 1.20. Learn more about MaEVe and its existing components through this [High-level design document](./docs/design.md).

## Pre-requisites

MaEVe runs in a set of Docker containers. This means you need to have `docker`, `docker-compose` and a docker daemon (e.g. docker desktop, `colima` or `rancher`) installed and running. 

## Getting started

To get the system up and running:

1. Run the [./scripts/generate-tls-cert.sh](./scripts/generate-tls-cert.sh) script which will create a server
certificate for the CSMS
2. Run the [./scripts/get-ca-cert.sh](./scripts/get-ca-cert.sh) script with a token retrieved from 
the [Hubject test environment](https://hubject.stoplight.io/docs/open-plugncharge/6bb8b3bc79c2e-authorization-token) (minus the "`Bearer `" part) 
to retrieve the V2G root certificate and CPO Sub CA certificates
3. Run the [./scripts/run.sh](./scripts/run.sh) script with the same token to run all the required components

Charge stations can connect to the CSMS using:
* `ws://localhost/ws/<cs-id>`
* `wss://localhost/ws/<cs-id>`

Charge stations can use either OCPP 1.6j or OCPP 2.0.1.

For TLS, the charge station should use a certificate provisioned using the 
[Hubject CPO EST service](https://hubject.stoplight.io/docs/open-plugncharge/486f0b8b3ded4-simple-enroll-iso-15118-2-and-iso-15118-20).

_At present, the set of charge stations that are accepted and their credentials are hard-coded
in [gateway/cmd/serve.go](gateway/cmd/serve.go)._

_The set of tokens that are accepted are also currently hard-coded in [manager/cmd/serve.go](manager/cmd/serve.go)._

## Configuration

All configuration in the system is currently through command-line flags. The available flags for each
component can be viewed using the `-h` flag. The configuration is mostly limited to connection details for the
various components and their dependencies. As mentioned in [Getting started](#getting-started) the allowed charge
stations and tokens are currently hard-coded in the server start up commands.

## Contributing

Learn more about how to contribute on this project through [Contributing](./CONTRIBUTING.md)

## License
MaEVe is [Apache licensed](./LICENSE).
