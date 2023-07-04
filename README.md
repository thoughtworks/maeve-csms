[![Manager](https://github.com/twlabs/ocpp2-broker-core/workflows/Manager/badge.svg)](https://github.com/twlabs/ocpp2-broker-core/actions/workflows/manager.yml)
[![Gateway](https://github.com/twlabs/ocpp2-broker-core/workflows/Gateway/badge.svg)](https://github.com/twlabs/ocpp2-broker-core/actions/workflows/gateway.yml)

# MaEVe

MaEVe is an EV charge station management system (CSMS). It began life as a simple proof of concept for
implementing ISO-15118-2 Plug and Charge (PnC) functionality and remains a work in progress. It is hoped that over
time it will become more complete, but already provides a useful basis for experimentation.

The system currently integrates with [Hubject](https://hubject.stoplight.io/) for PnC functionality. 

## Getting started

To get the system up and running:

1. Run the [./scripts/generate-tls-cert.sh](./scripts/generate-tls-cert.sh) script which will create a server
certificate for the CSMS
2. Run the [./scripts/get-ca-cert.sh](./scripts/get-ca-cert.sh) script with a token retrieved from 
the [Hubject test environment](https://hubject.stoplight.io/docs/open-plugncharge/6bb8b3bc79c2e-authorization-token)
to retrieve the V2G root certificate and CPO Sub CA certificates
3. Run the [./scripts/run.sh](./scripts/run.sh) script with the same token to run all the required components

Charge stations can connect to the CSMS using:
* `ws://localhost/ws/<cs-id>`
* `wss://localhost/ws/<cs-id>`

Charge stations can use either OCPP 1.6j or OCPP 2.0.1.

For TLS, the charge station should use a certificate provisioned using the 
[Hubject CPO EST service](https://hubject.stoplight.io/docs/open-plugncharge/486f0b8b3ded4-simple-enroll-iso-15118-2-and-iso-15118-20).

At present, the set of charge stations that are accepted and their credentials are hard-coded
in [gateway/cmd/serve.go](gateway/cmd/serve.go).

The set of tokens that are accepted are also currently hard-coded in [manager/cmd/serve.go](manager/cmd/serve.go).

## Developing

MaEVe is implemented in go 1.20. MaEVe currently has two components:
1. The [gateway](gateway) is a stateful service that accepts OCPP websocket connections from charge stations.
OCPP messages from charge stations are unpacked and published on an MQTT topic. OCPP messages from the
CSMS are read from another MQTT topic and sent to the charge station over the established websocket.
2. The [manager](manager) is a stateless service that reads from and writes to the relevant MQTT topics to
generate responses to OCPP calls that come from the charge station, make OCPP calls to the charge station and handle the 
associated responses.

To build the components, switch to the relevant directory and run:

```shell
$ go build ./...
$ go test ./...
```

To run the system an MQTT v5 broker, such as [Mosquitto](https://mosquitto.org/) is required. This is used to
decouple the stateful connections required by OCPP from the stateless message handlers. The system currently
connects to the MQTT broker anonymously.

The system also currently uses [Redis](https://redis.io/) as a storage engine for transaction details. This
is also connected to anonymously.

## Configuration

All configuration in the system is currently through command-line flags. The available flags for each
component can be viewed using the `-h` flag. The configuration is mostly limited to connection details for the
various components and their dependencies. As mentioned in [Getting started](#getting-started) the allowed charge
stations and tokens are currently hard-coded in the server start up commands.

