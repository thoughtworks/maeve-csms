package transport

import "context"

type MessageHandler interface {
	// Handle a Message produced by the broker. The message will be for a specific OcppVersion
	// and come from the charge station identified by the chargeStationId.
	Handle(ctx context.Context, chargeStationId string, message *Message)
}

type MessageHandlerFunc func(ctx context.Context, chargeStationId string, message *Message)

func (h MessageHandlerFunc) Handle(ctx context.Context, chargeStationId string, message *Message) {
	h(ctx, chargeStationId, message)
}

type Listener interface {
	// Connect establishes a connection to the broker and subscribes to receive messages for
	// either a specific charge station or all charge stations (for a specific OcppVersion). The messages
	// are delivered to the provided MessageHandler.
	//
	// Returns either a Connection on success or an error.
	Connect(ctx context.Context, ocppVersion OcppVersion, chargeStationId *string, handler MessageHandler) (Connection, error)
}

type Connection interface {
	// Disconnect drops the previously established connection to the broker.
	Disconnect(ctx context.Context) error
}
