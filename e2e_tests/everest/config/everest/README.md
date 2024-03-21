# Configuring EVerest

There are two configuration files for EVerest:

For OCPP 2.0.1
* [config-sil-ocpp201.yaml](config-sil-ocpp201.yaml)
* [ocpp/OCPP201/config.json](ocpp/OCPP201/config.json)

For OCPP 1.6
* [config-sil-ocpp.yaml](config-sil-ocpp.yaml)
* [ocpp/OCPP/config.json](ocpp/OCPP/config.json)

In addition, there are various certificates that can be configured.

## Overview

There are three parts to the main configuration file:
* settings
* active_modules
* x-module-layout

There is a [schema file](https://github.com/EVerest/everest-framework/blob/main/schemas/config.yaml)
that describes the syntax for the configuration file.

### Settings

There are a number of settings that we may configure. However, the principal setting is the address
of the MQTT broker which is set via the: `mqtt_broker_host` and `mqtt_broker_port` settings. The host
defaults to `localhost` and the port defaults to `1883`. These settings can also be overridden using
the environment variables `MQTT_SERVER_ADDRESS` and `MQTT_SERVER_PORT`.

### Active Modules

EVerest has a modular structure and these modules are connected via MQTT. The set of modules
is configured in the [config-sil-ocpp201.yaml](config-sil-ocpp201.yaml) file. Each module
takes a set of configuration parameters and sets up links to other modules.

You can see the configuration options, provided and consumed interfaces by looking at the
manifest files associated with each module. These can be found in the
[EVerest core modules directory](https://github.com/EVerest/everest-core/tree/main/modules).

The configuration is complex and requires a deep understanding of EVerest so we are using 
sample configurations from the
[EVerest core config directory](https://github.com/EVerest/everest-core/tree/main/config)
and then tweaking them.

### Module Layout

The `x-module-layout` key contains details of how to layout the modules in the management
UI. It has no other meaning beyond that.

## OCPP

There is a separate file for configuring OCPP. This is consumed by the EVerest dependency
[libocpp](https://github.com/EVerest/libocpp). Many of the sections in this configuration
have an almost 1:1 correspondence with Appendix 3 of the OCPP specification. However, there
are some additional sections.

The `InternalCtrlr` section provides a number of details about the charge station that will
be sent to the central system as part of a boot notification:
* ChargePointId
* ChargeBoxSerialNumber
* ChargePointModel
* ChargePointVendor

There are also details used for connecting to the central system:
* CentralSystemURI
* SupportedCiphers12
* SupportedCiphers13
* WebsocketReconnectInterval

Within the `SecurityCtrlr` section:

`SecurityProfile` specifies the security profile used when connecting:
* `1`: Basic auth over unsecured connection
* `2`: Basic auth over TLS
* `3`: Certificate auth over TLS

Additional examples of configuration files can be found in:
* [OCPP 1.6 configs](https://github.com/EVerest/libocpp/tree/main/config/v16)
* [OCPP 2.0.1 configs](https://github.com/EVerest/libocpp/tree/main/config/v201)

Schemas showing all the configuration options can be found in:
* [OCPP 1.6 schemas](https://github.com/EVerest/libocpp/tree/main/config/v16/profile_schemas)
* [OCPP 2.0.1 schemas](https://github.com/EVerest/libocpp/tree/main/config/v201/component_schemas)

## Certificates

The trusted CSMS CA certificate must be added to the `certs/ca/csms`
directory. If it is a DER encoded certificate then it must have a `.der` extension,
otherwise it is assumed to be a PEM encoded file.

The client certificate details must be added to the `certs/client/csms`
directory. The certificate and key must be PEM encoded and called `CSMS_LEAF.pem` and 
`CSMS_LEAF.key` respectively.


