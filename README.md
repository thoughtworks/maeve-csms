[![Manager](https://github.com/thoughtworks/Maeve-csms/workflows/Manager/badge.svg)](https://github.com/thoughtworks/Maeve-csms/actions/workflows/manager.yml)
[![Gateway](https://github.com/thoughtworks/Maeve-csms/workflows/Gateway/badge.svg)](https://github.com/thoughtworks/Maeve-csms/actions/workflows/gateway.yml)

# MaEVe

Maeve is an EV Charge Station Management System (CSMS). It began life as a simple proof of concept for
implementing ISO-15118-2 Plug and Charge (PnC) functionality and remains a work in progress. Over
time it will become more complete. It provides a useful basis for experimentation with [Everest](https://github.com/EVerest).

Maeve integrates with [Hubject](https://hubject.stoplight.io/) for Plug and Charge (PnC) functionality.

## Documentation

Learn more about Maeve and its existing components through this [High-level design document](./docs/design.md).

## Pre-requisites

Maeve runs in a set of Docker containers. This means you need to have `docker`, `docker-compose` and a docker daemon (e.g. [colima](https://github.com/abiosoft/colima)) installed and running.

## Getting started

To get the system up and running:

1. `(cd config/certificates && make)`
2. `./scripts/run.sh`

Charge stations connect to Maeve using Web sockets:

* `ws://localhost/ws/<cs-id>`
* `wss://localhost/ws/<cs-id>`

Charge stations can use either OCPP 1.6j or OCPP 2.0.1.

For TLS, the charge station should use a certificate provisioned using the
[Hubject CPO EST service](https://hubject.stoplight.io/docs/open-plugncharge/486f0b8b3ded4-simple-enroll-iso-15118-2-and-iso-15118-20).

A charge station must register with Maeve before use, via the
[manager API](./manager/api/API.md). e.g. for unsecured transport with basic auth use:

```shell
$ cd manager
$ ENC_PASSWORD=$(go run main.go auth encode-password <password> | cur -d' ' -f2)
$ curl http://localhost:9410/api/v0/cs/<cs-id> -H 'content-type: application/json' -d '{"securityProfile":0,"base64SHA256Password":"'$ENC_PASSWORD'"}'
```

Tokens are an authentication mechanism that needs to be set before use:

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

Note Maeve doesn't persist any tokens between restarts.