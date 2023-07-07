// SPDX-License-Identifier: Apache-2.0

// Package pipe provides the core functionality for transferring messages between
// a charge station and a CSMS using OCPP/J.
// A Pipe is a bidirectional connection where OCPP Call and CallResult (or CallError)
// messages are brokered. A Pipe is agnostic to the actual communication channels. It
// is up to the caller to ensure that the charge station and CSMS Rx and Tx channels
// are hooked up to an appropriate source or sink.
package pipe
