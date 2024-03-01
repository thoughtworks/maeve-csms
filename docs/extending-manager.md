# Extending the Manager

## Adding a new call handler

1. Compile the JSON request and response schemas using [json-to-go.sh](../scripts/json-to-go.sh)
2. Remove all the custom serialization from the generated classes in the appropriate [ocpp subdirectory](../manager/ocpp)
3. Implement the ocpp.Request or ocpp.Response stub method on the Request or Response type
4. Implement the ocpp.CallHandler interface in the appropriate [handlers subdirectory](../manager/handlers)
5. Add the new handler to the appropriate router (either [ocpp16](../manager/handlers/ocpp16/routing.go) or [ocpp201](../manager/handlers/ocpp201/routing.go))

The process for DataTransfer messages is identical. The only difference is the router configuration where there is a
generic data-transfer handler that uses the data transfer vendor and message ids to route the message.

## Adding a new call response handler

As above, but implement the ocpp.CallResultHandler interface.

## Adding a new CSMS initiated message

Compile up the schemas (as above) and add the new handler to the appropriate call maker (either [ocpp16](../manager/handlers/ocpp16/routing.go) or [ocpp201](../manager/handlers/ocpp201/routing.go))

## Adding a new API method

1. Add the method to the [OpenAPI](../manager/api/api-spec.yaml) schema.
2. Use go generate to compile the associated types (see [../manager/api/api.go](../manager/api/api.go)).
3. Update the [server](../manager/api/server.go) so it implements the [api.ServerInterface](../manager/api/api.gen.go)