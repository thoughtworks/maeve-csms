// SPDX-License-Identifier: Apache-2.0

package transport

import (
	"context"
)

// A Router processes messages received by a Receiver and dispatches them to the
// appropriate message handler for processing.
type Router interface {
	// GetOcppVersion returns the OcppVersion that this router processes.
	GetOcppVersion() OcppVersion
	// Route takes a message from a specific charge station - identified by a chargeStationId - and
	// dispatches it to the appropriate message handler for processing.
	Route(ctx context.Context, chargeStationId string, message Message) error
}
