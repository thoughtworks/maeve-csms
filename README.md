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
Scripts that fetch various tokens use `jq`. Make sure you have it installed.

## Getting started

To get the system up and running:

1. `(cd config/certificates && make)`
2. Run the [./scripts/run.sh](./scripts/run.sh) script

Charge stations can connect to the CSMS using:
* `ws://localhost/ws/<cs-id>`
* `wss://localhost/ws/<cs-id>`

If the charge station is also running in a Docker container then the charge
station docker container can connect to the `maeve-csms` network and the
charge station can connect to the CSMS using:
* `ws://gateway:9310/ws/<cs-id>`
* `wss://gateway:9311/ws/<cs-id>`

Charge stations can use either OCPP 1.6j or OCPP 2.0.1.

For TLS, the charge station should use a certificate provisioned using the
[Hubject CPO EST service](https://hubject.stoplight.io/docs/open-plugncharge/486f0b8b3ded4-simple-enroll-iso-15118-2-and-iso-15118-20).

A charge station must first be registered with the CSMS before it can be used. This can be done using the
[manager API](./manager/api/API.md). e.g. for TLS with client certificate, use:

```shell
$ curl http://localhost:9410/api/v0/cs/<cs-id> -H 'content-type: application/json' -d '{"securityProfile":2}'
```

Tokens, which identify a payment method for a non-contract charge, must also be registered with the CSMS before they can be used. This can also be done using the
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

## Troubleshooting

Docker compose doesn't always rebuild the docker images which can cause all kinds of errors. If in doubt, force a rebuild by `docker-compose build` before launching containers.

`java.io.IOException: keystore password was incorrect`
This error results from incompatibility between java version and openssl; try upgrading your java version.

## Configuration

The gateway is configured through command-line flags. The available flags can be viewed using the `-h` flag. 

The manager is configured through a TOML configuration file. An example configuration file can be found in 
[./config/manager/config.toml](./config/manager/config.toml). Details of the available configuration options
can be found in [./manager/config/README.md](./manager/config/README.md).

## Contributing

Learn more about how to contribute on this project through [Contributing](./CONTRIBUTING.md)

## License
MaEVe is [Apache licensed](./LICENSE).


