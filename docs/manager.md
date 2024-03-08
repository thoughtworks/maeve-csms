# Manager

The manager is stateless and implements the core logic for processing OCPP messages. 

Messages are received from the gateway using a transport (currently only MQTT 5) and are 
then routed to a handler that will process that message. If the incoming message was an 
OCPP call then an OCPP call result will be emitted.

The CSMS may also emit messages in order to manage the charge stations.

The manager is configured using a TOML file. This configuration is defined in the
[config](../manager/config) package which also documents the available options.

There is an administration API that allows the CSMS to be configured. This is defined in 
the [api](../manager/api) package with [API documentation](../manager/api/API.md).

Support for OCPI is provided by the [ocpi](../manager/ocpi) package.

The structure of the manager source code is:
```
manager/
├─ adminui/       Administration UI
├─ api/           Administration API
├─ cmd/           Executable commands
├─ config/        Configuration management and dependency injection 
├─ handlers/      Common implementations for handling OCPP messages
│  ├─ has2be/     Handlers for the Has2Be OCPP 1.6 extension messages 
│  ├─ ocpp16/     Handlers for OCPP 1.6 messages
│  ├─ ocpp201/    Handlers for OCPP 2.0.1 messages
├─ ocpi/          OCPI API
├─ ocpp/          Common types for OCPP messages
│  ├─ has2be/     Types representing the Has2Be OCPP 1.6 extension messages
│  ├─ ocpp16/     Types representing OCPP 1.6 messages
│  ├─ ocpp201/    Types representing OCPP 2.0.1 messages
├─ schemas/       Support for schema validation
│  ├─ has2be/     JSON schema files for the Has2Be OCPP 1.6 extension messages
│  ├─ ocpp16/     JSON schema files for the OCPP 1.6 messages
│  ├─ ocpp201/    JSON schema files for the OCPP 2.0.1 messages
├─ server/        Support for providing HTTP-based endpoints
├─ services/      Pluggable implementations used by handlers
├─ store/         Interface for interacting with the persistent store
│  ├─ firestore/  Persistent store implementation using Google Firestore
│  ├─ inmemory/   In-memory implementation of the persistent store (for testing) 
├─ sync/          Synchronize configuration to charge stations
├─ transport/     Interface for sending/receiving messages
│  ├─ mqtt/       Transport interface implemented using MQTT
```

The incoming message flow (for MQTT) is:
```
├─ [receiver](../manager/transport/mqtt/receiver.go)  Receives the call message from the gateway
│  ├─ [router](../manager/handlers/router.go)         Routes the message to the appropriate handler
│  │  ├─ handler                                      Provides the message specific behaviour
│  ├─ [emitter](../manager/transport/mqtt/emitter.go) Emits the response to the gateway
```

The outgoing message flow (for MQTT) is:
```
├─ [call maker](../manager/handlers/call_maker.go)    Encodes the message
│  ├─ [emitter](../manager/transport/mqtt/emitter.go) Emits the message to the gateway
```
