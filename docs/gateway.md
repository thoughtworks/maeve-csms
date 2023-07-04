# Gateway

The gateway component is stateful and is responsible for accepting connections from charge stations and proxying
OCPP messages to/from the stateless backend via the MQTT broker.

The core logic of the gateway is implemented by the [Pipe](../gateway/pipe/pipe.go) type. The pipe has four main 
go channels:
1. `ChargeStationRx` - receives messages from the charge station
2. `ChargeStationTx` - transmits messages to the charge station
3. `CSMSRx` - receives messages from the CSMS manager
4. `CSMSTx` - transmits messages to the CSMS manager

![Diagram showing the gateway reading from a websocket connected to the charge station sending it to the MQTT broker in topic via the Pipe whilst subscribing to the MQTT broker out topic and writing it to the charge station websocket via the pipe](assets/gateway.png)

The Pipe implements a small state machine. The state machine has 3 states:
1. `Waiting` - no OCPP call is in progress, a call will be accepted from the charge station or the CSMS manager.
Calls from the charge station are prioritized (as the charge station will be expecting a response).
2. `ChargeStationCall` - an OCPP call initiated by the charge station is in progress, any calls from the CSMS
manager will be buffered waiting for the call response for the ongoing call.
3. `CSMSCall` - on OCPP call initiated by the CSMS manager is in progress, will wait for a response from the
charge station.

Both the call states have an associated time-out that returns them to the `Waiting` state. However, any late
arriving response from the charge station will be forwarded to the CSMS manager and any late arriving response
from the CSMS manager will be forwarded to the charge station as long as another call has not occurred.

The `Pipe` works in terms of the [GatewayMessage](../gateway/pipe/message.go) type which is an amalgamation of the
OCPP call, call result and call error messages. The Pipe caches CSMS originated calls and merges these with the 
call result or call error received from the charge station before sending them on to the CSMS manager. This is 
needed as the CSMS manager is stateless and has no knowledge of the context of a call result. 

The [WebsocketHandler](../gateway/server/ws.go) establishes the websocket connection and implements the interfaces 
between the various `Pipe` channels and the websocket - converting between OCPP messages and gateway messages - 
on one side and the MQTT broker on the other: subscribing to outgoing (to the charge station) messages and
publishing incoming (from the charge station) messages.

There are individual MQTT topics for each charge station:
* `<prefix>/in/<ocpp-version>/<cs-id>`
* `<prefix>/out>/<ocpp-version>/<cs-id>`

Where `<prefix>` is a configured prefix for all the topics (defaults to `cs`), `<ocpp-version>` is the
version of OCPP being used: either `ocpp16` or `ocpp201` and `<cs-id>` is the charge station identifier.
