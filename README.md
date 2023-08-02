[![Manager](https://github.com/thoughtworks/maeve-csms/workflows/Manager/badge.svg)](https://github.com/thoughtworks/maeve-csms/actions/workflows/manager.yml)
[![Gateway](https://github.com/thoughtworks/maeve-csms/workflows/Gateway/badge.svg)](https://github.com/thoughtworks/maeve-csms/actions/workflows/gateway.yml)

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

1. `(cd config/certificates && make)`
2. Run the [./scripts/run.sh](./scripts/run.sh) script with the same token to run all the required components - again, don't forget the quotes around the token

Charge stations can connect to the CSMS using:
* `ws://localhost/ws/<cs-id>`
* `wss://localhost/ws/<cs-id>`

Charge stations can use either OCPP 1.6j or OCPP 2.0.1.

For TLS, the charge station should use a certificate provisioned using the
[Hubject CPO EST service](https://hubject.stoplight.io/docs/open-plugncharge/486f0b8b3ded4-simple-enroll-iso-15118-2-and-iso-15118-20).

A charge station must first be registered with the CSMS before it can be used. This can be done using the
[manager API](./manager/api/API.md). e.g. for unsecured transport with basic auth use:

```shell
$ cd manager
$ ENC_PASSWORD=$(go run main.go auth encode-password <password> | cur -d' ' -f2)
$ curl http://localhost:9410/api/v0/cs/<cs-id> -H 'content-type: application/json' -d '{"securityProfile":0,"base64SHA256Password":"'$ENC_PASSWORD'"}'
```

Tokens must also be registered with the CSMS before they can be used. This can also be done using the
[manager API](./manager/api/API.md). e.g.:

```shell
$ curl http://localhost:9410/api/v0/token -H 'content-type: application/json' -d '{
  "countryCode": "GB",
  "partyId": "TWK",
  "type": "RFID",
  "uid": "DEADBEEF",
  "contractId": "GBTWK012345678V",
  "issuer": "Thoughtworks",
  "valid": true,
  "cacheMode": "ALWAYS"
}'
```

## Configuration

All configuration in the system is currently through command-line flags. The available flags for each
component can be viewed using the `-h` flag. The configuration is mostly limited to connection details for the
various components and their dependencies. As mentioned in [Getting started](#getting-started) the allowed charge
stations and tokens are currently hard-coded in the server start up commands.

## Contributing

Learn more about how to contribute on this project through [Contributing](./CONTRIBUTING.md)

## License
MaEVe is [Apache licensed](./LICENSE).
