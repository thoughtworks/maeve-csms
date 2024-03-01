// SPDX-License-Identifier: Apache-2.0

package transport

// Receiver is the interface implemented in order to receive messages
// from the gateway (and thus from the charge station). The Receiver
// must establish a connection to the gateway and then use a Router
// to process any messages that are received.
type Receiver interface {
	Connect(errCh chan error)
}
