# Manager

The manager is stateless and implements the core logic for processing OCPP messages. It uses an MQTT v5
shared subscription to read messages published to the incoming MQTT topic for any charge station and routes
them to an appropriate handler. This means that different messages from a charge station may be distributed
to different instances of the manager. The MQTT subscription is managed by the 
[Handler](../manager/mqtt/handler.go) type. 

The messages on the MQTT topics correspond to the [Message](../manager/mqtt/messages.go) type. This is an
amalgamation of the OCPP call, call result and call error. Responses to CSMS initiated calls will include
the original call details for context.

Messages are routed to an appropriate OCPP handler by the [Router](../manager/mqtt/router.go). The
`Router` validates the message using the OCPP JSON schemas and (if all is well) uses the OCPP version (from the topic),
OCPP message type and OCPP action (from the message) to determine which handler to invoke. The router is responsible
for creating the handler instances and providing them with their dependencies.

![Diagram showing MQTT handler subscribing to incoming messages from the MQTT broker and processing them via either an OCPP 1.6 or OCPP 2.0.1 router which distributes them to a handler. A handler may initiate a call with a Call Maker which publishes a message to the MQTT broker](manager.png)

All integration points are decoupled from the OCPP handlers and provided to the handler as dependencies. At present 
the following integration points are used:
* [CertificateSignerService](../manager/services/certificate_signer.go): provides support for signing the TLS
certificates used by charge stations (for either the V2G or CPO interfaces). Current implementation uses Hubject EST
and assumes the same certificate will be used for both interfaces.
* [CertificateValidationService](../manager/services/certificate_validation.go): provides support for verifying a
contract certificate chain based on either receiving the chain as a PEM encoded string or via certificate element
hashes. Current implementation uses an OCSP server to verify that the certificates have not been revoked.
* [EVCertificateProvider](../manager/services/ev_certificate_provider.go): provides supports for retrieving a contract
certificate given an ISO-15118-2 get EV certificate request. Current implementation uses the Hubject contract
certificate pool.
* [TariffService](../manager/services/tariff.go): provides support for calculating the cost of a transaction. Current
implementation applies a fixed charge per kWh energy consumed.
* [TokenStore](../manager/services/token.go): provides support for looking up token details. Current implementation
uses an in-memory map.
* [TransactionStore](../manager/services/transaction.go): provides support for storing details of an ongoing
transaction. Current implementations are either in-memory (N.B. this is not stateless) or use Redis as a store.

There are separate interfaces that a handler implements for handling a call
(a [CallHandler](../manager/handlers/types.go)) and a call result (a [CallResultHandler](../manager/handlers/types.go)).
A handler that wants to initiate a call to a charge station can use a [CallMaker](../manager/handlers/types.go)
to send an OCPP message to a charge station.