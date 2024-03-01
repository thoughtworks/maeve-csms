# Configuration

Configuration is defined in TOML.

## Table of Contents

* [General settings](#general-settings)
* [Service settings](#service-settings)
* [Transport](#transport)
* [Storage](#storage)
* [Contract certificate validator](#contract-certificate-validator)
* [Contract certificate provider](#contract-certificate-provider)
* [Charge station certificate provider](#charge-station-certificate-provider)
* [Tariff service](#tariff-service)
* [Root certificate provider](#root-certificate-provider)
* [Http auth service](#http-auth-service)
* [Example configuration](#example-configuration)

## General settings

| Section       | Key                 | Type   | Description                                                          |
|---------------|---------------------|--------|----------------------------------------------------------------------|
| api           | addr                | string | Address that API server will listen on, e.g. localhost:9410          |
| api           | external_addr       | string | The Externally visible URL that the server is available on           |
| api           | org_name            | string | The organization name to use when issuing client certificates        |
| ocpp          | heartbeat_interval  | string | Frequency to request charge station heartbeat messages at, e.g. "5m" |
| ocpp          | ocpp16_enabled      | bool   | Is OCPP 1.6 support enabled, e.g. "true"?                            |
| ocpp          | ocpp201_enabled     | bool   | Is OCPP 2.0.1 support enabled, e.g. "true"?                          |
| observability | log_format          | string | Either "json" or "text"                                              |
| observability | otel_collector_addr | string | Address of the OpenTelemetry collector, e.g. "localhost:4317"        |
| observability | tls_keylog_file     | string | File where TLS session keys will be written for use with Wireshark   |

## Transport settings

This section consists of a `type` parameter and a set of parameters specific to that type prefixed by the type name.

e.g.

```toml
[transport]
type = "mqtt"
mqtt.urls = ["mqtt://mqtt-server.example.com:1883"]
```

### MQTT

Configures the MQTT transport.

| Section | Key                 | Type             | Description                                            |
|---------|---------------------|------------------|--------------------------------------------------------|
| mqtt    | urls                | array of strings | List of MQTT broker URLs, e.g. [mqtt://localhost:1883] |
| mqtt    | prefix              | string           | MQTT topic prefix, e.g. "cs"                           |
| mqtt    | group               | string           | MQTT subscriber group name, e.g. "manager"             |
| mqtt    | connect_timeout     | string           | MQTT connection timeout, e.g. "10s"                    |
| mqtt    | connect_retry_delay | string           | MQTT connection retry delay, e.g. "1s"                 |
| mqtt    | keep_alive_interval | string           | MQTT keep alive interval, e.g. "10s"                   |

## Service settings

The following types of service can be configured, each service has its own section:
* [`storage`](#storage) - configures the backing store will be used
* [`contract_cert_validator`](#contract-certificate-validator) - configures how contract certificates will be validated
* [`contract_cert_provider`](#contract-certificate-provider) - configures how contract certificates are provided
* [`charge_station_cert_provider`](#charge-station-certificate-provider) - configures how charge station certificates are provided
* [`tariff_service`](#tariff-service) - configures how tariffs are calculated

Each section consists of a `type` parameter and a set of parameters specific to that type prefixed by the type name.

e.g.

```toml
[storage]
type = "firestore"
firestore.project_id = "my-google-project"
```

### Storage

There are two storage implementations:
* [`firestore`](#firestore) - Google Firestore
* [`in_memory`](#in-memory) - in-memory storage for testing

#### Firestore

| Key        | Type   | Description             |
|------------|--------|-------------------------|
| project_id | string | Google Cloud project ID |

#### In-memory

There is no additional configuration for in-memory storage.

### Contract certificate validator

There is just one contract certificate validator implementation:
* [`ocsp`](#ocsp-contract-certificate-validator) - checks the certificate chain and validates the OCSP status of each provided certificate 

#### OCSP contract certificate validator

| Key          | Type                                           | Description                                                          |
|--------------|------------------------------------------------|----------------------------------------------------------------------|
| root_certs   | [RootCertProvider](#root-certificate-provider) | Configures how to retrieve the trusted root certificates             |
| max_attempts | int                                            | Maximum number of attempts to check the OCSP status of a certificate |

### Contract certificate provider

There are two contract certificate provider implementations:
* [`opcp`](#opcp-contract-certificate-provider) - contract certificates are retrieved from a contract certificate pool using the Open Plug&Charge Protocol (OPCP)
* [`default`](#default-contract-certificate-provider) - returns an error for all requests

#### OPCP contract certificate provider

| Key        | Type                                  | Description                                                           |
|------------|---------------------------------------|-----------------------------------------------------------------------|
| url        | string                                | Base URL for OPCP service that provides the contract certificate pool |
| auth       | [HttpAuthService](#http-auth-service) | Configures how to authenticate with the OPCP service                  |

#### Default contract certificate provider

There is no additional configuration for the default contract certificate provider.

### Charge station certificate provider

There are four charge station certificate provider implementations:
* [`opcp`](#opcp-charge-station-certificate-provider) - charge station certificates are issued using the EST service from the Open Plug&Charge Protocol (OPCP)
* [`local`](#local-charge-station-certificate-provider) - charge station certificates are issued using a CA implemented by the CSMS
* [`delegating`](#delegating-charge-station-certificate-provider) - supports different charge station certificate providers for issuing V2G and CSO certificates
* [`default`](#default-charge-station-certificate-provider) - returns an error for all requests

#### OPCP charge station certificate provider

| Key        | Type                                  | Description                                             |
|------------|---------------------------------------|---------------------------------------------------------|
| url        | string                                | Base URL for OPCP service that provides the EST service |
| auth       | [HttpAuthService](#http-auth-service) | Configures how to authenticate with the OPCP service    |

#### Local charge station certificate provider

| Key  | Type                         | Description                                                                         |
|------|------------------------------|-------------------------------------------------------------------------------------|
| cert | [LocalSource](#local-source) | The source that provides the signing certificate, must be a PEM encoded certificate |
| key  | [LocalSource](#local-source) | The source that provides the signing key, must be a PEM encoded private key         |

#### Delegating charge station certificate provider

| Key | Type                                                                     | Description                                                  |
|-----|--------------------------------------------------------------------------|--------------------------------------------------------------|
| v2g | [ChargeStationCertificateProvider](#charge-station-certificate-provider) | The charge station certificate provider for V2G certificates |
| cso | [ChargeStationCertificateProvider](#charge-station-certificate-provider) | The charge station certificate provider for CSO certificates |

#### Default charge station certificate provider

There is no additional configuration for the default charge station certificate provider.

### Tariff service

There is a single tariff service implementation:
* [`kwh`](#kwh-tariff-service) - calculates the tariff based on the energy consumed

#### kWh tariff service

There is no additional configuration for the kWh tariff service.

### Root certificate provider

There are several implementations of RootCertProvider:
* [`opcp`](#opcp-root-certificate-provider) - root certificates are retrieved from a root certificate pool using the Open Plug&Charge Protocol (OPCP)
* [`file`](#file-root-certificate-provider) - root certificates are retrieved from a file

#### OPCP root certificate provider

| Key  | Type                                  | Description                                                       |
|------|---------------------------------------|-------------------------------------------------------------------|
| url  | string                                | Base URL for OPCP service that provides the root certificate pool |
| ttl  | string                                | Time before cached values are discarded, e.g. "1h"                |
| auth | [HttpAuthService](#http-auth-service) | Configures how to authenticate with the OPCP service              |

#### File root certificate provider

| Key   | Type             | Description                                |
|-------|------------------|--------------------------------------------|
| files | array of strings | List of files containing root certificates |

### Http auth service

There are several implementation of HttpAuthService:
* [`env_token`](#environment-token-auth-service) - token is read from an environment variable
* [`fixed_token`](#fixed-token-auth-service) - token is read from the configuration
* [`oauth2_token`](#oauth2-token-auth-service) - token is retrieved using OAuth2 client credentials grant
* [`hubject_test_token`](#hubject-test-token-auth-service) - token is scraped from the Hubject test environment authorization page

#### Environment token auth service

| Key      | Type   | Description                                          |
|----------|--------|------------------------------------------------------|
| variable | string | Name of the environment variable                     |

#### Fixed token auth service

| Key   | Type   | Description     |
|-------|--------|-----------------|
| token | string | The token value |

#### OAuth2 token auth service

| Key                   | Type   | Description                                                                                    |
|-----------------------|--------|------------------------------------------------------------------------------------------------|
| url                   | string | The URL of the OAuth2 Authorization Server token endpoint                                      |
| client_id             | string | The client id to use in the client credentials grant                                           |
| client_secret         | string | The client secret to use in the client credentials grant                                       |
| client_secret_env_var | string | The environment variable to read the client secret from to use in the client credentials grant |

#### Hubject test token auth service

| Key | Type   | Description                                            |
|-----|--------|--------------------------------------------------------|
| url | string | URL of the Hubject test environment authorization page |
| ttl | string | The duration for which the token is valid, e.g. "1h"   |

### Local source

There are two different local source implementations:
* [`file`](#file-local-source) - data is read from a file
* [`google_cloud_secret`](#google-cloud-secret-local-source) - data is read from a google cloud secret

#### File local source

The value is a file to be read from the operating system.

#### Google cloud secret local source

The value is the name of a secret to be read from google cloud secrets. The name must be of the form:

```
projects/<project-number>/secrets/<secret-name>/[latest|<version>]
```

## Example configuration

```toml
[api]
addr = ":9410"

[mqtt]
urls = ["mqtt://mqtt:1883"]

[observability]
otel_collector_addr = "otel-collector:4317"

[storage]
type = "firestore"
firestore.project_id = "*detect-project-id*"

[contract_cert_validator]
type = "ocsp"

[contract_cert_validator.ocsp.root_certs]
type = "opcp"
opcp.url = "https://open.plugncharge-test.hubject.com/mo/cacerts/ISO15118-2"
opcp.auth.type = "env_token"
opcp.auth.env_token.variable = "HUBJECT_TOKEN"

[contract_cert_provider]
type = "opcp"
opcp.url = "https://open.plugncharge-test.hubject.com"
opcp.auth.type = "env_token"
opcp.auth.env_token.variable = "HUBJECT_TOKEN"

[charge_station_cert_provider]
type = "opcp"
opcp.url = "https://open.plugncharge-test.hubject.com"
opcp.auth.type = "env_token"
opcp.auth.env_token.variable = "HUBJECT_TOKEN"
```
