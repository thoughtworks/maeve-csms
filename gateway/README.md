# OCPP Websocket to MQTT Gateway

The gateway decouples the processing of OCPP messages from the underlying
websocket transport. OCPP messages are read from the websocket connected to the
charge station and published on a per charge station MQTT topic for the CSMS 
logic to consume. Responses and CSMS initiated calls are also sent by publishing 
an OCPP message to a per charge station MQTT topic which the gateway reads and
sends to charge station over the websocket.

Messages from the charge station are published on a topic named:

`<prefix>/in/<protocol>/<charge-station-id>`

Messages for the charge station are published on a topic named:

`<prefix>/out/<protocol>/<charge-station-id>`

Where `<prefix>` is a configured prefix; `<protocol>` is the OCPP protocol
version (either `ocpp1.6` or `ocpp2.0.1`) and `<charge-station-id>` is the
unique charge station id read from the websocket connection URL.

OCPP uses a bidirectional request/response model - but all request/response 
pairs must be serialized (i.e. you cannot send a CSMS request when the charge
station is expecting a response to a request it initiated). The gateway does
its best to ensure that this constraint is complied with (although there is
always a chance of messages crossing during transmission).

The gateway works exclusively at the OCPP/J layer. Errors will be sent if a 
message is syntactically invalid OCPP/J, but the gateway is agnostic to the
payload. The gateway will endeavour to always send a response to every call:
if a valid response is not received from the corresponding party then the
gateway will respond with an error.

The gateway is essentially stateless and can be scaled horizontally behind 
a load-balancer. Messages to the charge station will always be read by the
instance that currently has the open websocket. 

## Building

The project depends on the CSMS broker code. This is in a private Github 
repository. To allow go to read this repository:

```shell
git config --global --add url."ssh://git@github.com/twlabs/".insteadOf "https://github.com/twlabs/"
```

And export the following variables in your shell:

```shell
export GOPRIVATE=github.com/twlabs
```

For local development, you can also add a `go.work` file to make things easier.
In the root directory of the ocpp2-broker-core project:

```shell
$ go work init
$ go work use csms serve
```

To build the dockerfile, ensure you have an ssh-agent running and that your
github private keys are available, e.g. via ssh-add, then run:

```shell
docker build --ssh default -t gw .
```