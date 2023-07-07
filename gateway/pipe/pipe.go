// SPDX-License-Identifier: Apache-2.0

package pipe

import (
	"container/ring"
	"github.com/thoughtworks/maeve-csms/gateway/ocpp"
	"log"
	"time"
)

// Pipe provides a bidirectional RPC pipe between a ChargeStation and the CSMS
type Pipe struct {
	// ChargeStationRx is an incoming channel from the Charge Station
	ChargeStationRx chan *GatewayMessage
	// ChargeStationTx is on outgoing channel to the Charge Station
	ChargeStationTx chan *GatewayMessage
	// CSMSRx is an incoming channel from the CSMS
	CSMSRx chan *GatewayMessage
	// CSMSTx is an outgoing channel to the CSMS
	CSMSTx chan *GatewayMessage
	// halt in a signal channel used to stop the pipe
	halt chan struct{}
	// csmsRxCallBuf is a channel that buffers CSMS calls whilst waiting for a CallResult
	csmsRxCallBuf chan *GatewayMessage
	// responseTimeout is the duration that the pipe will wait for a response to a call
	responseTimeout time.Duration
	// messageIdBufferLen is the number of previously used message ids that will be stored
	messageIdBufferLen int
	// csmsMessageQueueLen is the maximum number of CSMS messages that will be queued for delivery
	csmsMessageQueueLen int
	// csmsCallQueueLen is the maximum number of CSMS call messages that will be queued looking for a CallResult
	csmsCallQueueLen int
	// csmsCallResponseBufferLen is the maximum number of CSMS messages that will be cached waiting for a response
	csmsCallResponseBufferLen int
}

// NewPipe creates a pipe for connecting a charge station to a CSMS
func NewPipe(opts ...Opt) *Pipe {
	pipe := new(Pipe)

	ensureDefaults(pipe)

	for _, opt := range opts {
		opt(pipe)
	}

	pipe.ChargeStationRx = make(chan *GatewayMessage, 1)
	pipe.ChargeStationTx = make(chan *GatewayMessage, 1)
	pipe.CSMSRx = make(chan *GatewayMessage, pipe.csmsMessageQueueLen)
	pipe.CSMSTx = make(chan *GatewayMessage, 1)
	pipe.halt = make(chan struct{}, 1)
	pipe.csmsRxCallBuf = make(chan *GatewayMessage, pipe.csmsCallQueueLen)

	return pipe
}

func ensureDefaults(conn *Pipe) {
	conn.responseTimeout = 10 * time.Second
	conn.messageIdBufferLen = 10
	conn.csmsMessageQueueLen = 5
	conn.csmsCallQueueLen = 5
	conn.csmsCallResponseBufferLen = 5
}

type Opt func(*Pipe)

// WithResponseTimeout is a pipe option that sets the response timeout
func WithResponseTimeout(timeout time.Duration) Opt {
	return func(p *Pipe) {
		p.responseTimeout = timeout
	}
}

// WithMessageIdBufferLen is a pipe option that sets the number of previously used message ids to store
func WithMessageIdBufferLen(bufferLen int) Opt {
	return func(p *Pipe) {
		p.messageIdBufferLen = bufferLen
	}
}

// WithCSMSMessageQueueLen is a pipe option that sets the maximum number of CSMS messages to queue
func WithCSMSMessageQueueLen(queueLen int) Opt {
	return func(p *Pipe) {
		p.csmsMessageQueueLen = queueLen
	}
}

// WithCSMSCallQueueLen is a pipe option that sets the maximum number of CSMS Call messages to queue when looking for
// a CallResult to send to the charge station
func WithCSMSCallQueueLen(queueLen int) Opt {
	return func(p *Pipe) {
		p.csmsCallQueueLen = queueLen
	}
}

// WithCSMSCallResponseBufferLen is a pipe option that sets the maximum number of CSMS Call messages to buffer waiting
// for CallResult responses from the charge station
func WithCSMSCallResponseBufferLen(queueLen int) Opt {
	return func(p *Pipe) {
		p.csmsCallResponseBufferLen = queueLen
	}
}

type Status int

const (
	StatusWaiting Status = iota
	StatusChargeStationCall
	StatusCSMSCall
)

// Start will begin transferring RPC messages between the ChargeStation and CSMS
func (p Pipe) Start() {
	go func() {
		processedMessageIds := ring.New(p.messageIdBufferLen)
		processedCSMSCalls := ring.New(p.csmsCallResponseBufferLen)

		status := StatusWaiting

		for {
			var currentMsg *GatewayMessage
			if processedCSMSCalls.Value != nil {
				currentMsg = processedCSMSCalls.Value.(*GatewayMessage)
			}

			switch status {
			case StatusWaiting:
				// prioritise pending calls from the ChargeStation
				select {
				case msg := <-p.ChargeStationRx:
					if msg.MessageType == ocpp.MessageTypeCall {
						// call from CS
						if findMessageId(processedMessageIds, msg.MessageId) {
							log.Printf("error: CS message id %s reused", msg.MessageId)
							continue
						}
						processedMessageIds = processedMessageIds.Next()
						processedMessageIds.Value = msg.MessageId
						status = StatusChargeStationCall
						p.CSMSTx <- msg
					} else if currentMsg != nil && msg.MessageId == currentMsg.MessageId {
						// call result / call error for current CSMS call from CS
						log.Printf("warn: CS call response for message %s is late", msg.MessageId)
						msg.Action = currentMsg.Action
						msg.RequestPayload = currentMsg.RequestPayload
						p.CSMSTx <- msg
					} else if csmsCall := findCSMSCall(processedCSMSCalls, msg.MessageId); csmsCall != nil {
						// call result / call error for previous CSMS call from CS
						log.Printf("warn: CS call response for message %s is very late", msg.MessageId)
						msg.Action = csmsCall.Action
						msg.RequestPayload = csmsCall.RequestPayload
						p.CSMSTx <- msg
					} else {
						log.Printf("error: CS call response with message id %s has no corresponding CSMS call", msg.MessageId)
						continue
					}
				case <-p.halt:
					return
				default:
					// if no pending calls from the ChargeStation wait for next call from anywhere
					select {
					case msg := <-p.ChargeStationRx:
						if msg.MessageType == ocpp.MessageTypeCall {
							// call from CS
							if findMessageId(processedMessageIds, msg.MessageId) {
								log.Printf("error: CS message id %s reused", msg.MessageId)
								continue
							}
							processedMessageIds = processedMessageIds.Next()
							processedMessageIds.Value = msg.MessageId
							status = StatusChargeStationCall
							p.CSMSTx <- msg
						} else if currentMsg != nil && msg.MessageId == currentMsg.MessageId {
							// call result / call error for current CSMS call from CS
							log.Printf("warn: CS call response for message %s is late", msg.MessageId)
							msg.Action = currentMsg.Action
							msg.RequestPayload = currentMsg.RequestPayload
							p.CSMSTx <- msg
						} else if csmsCall := findCSMSCall(processedCSMSCalls, msg.MessageId); csmsCall != nil {
							// call result / call error for previous CSMS call from CS
							log.Printf("warn: CS call response for message %s is very late", msg.MessageId)
							msg.Action = csmsCall.Action
							msg.RequestPayload = csmsCall.RequestPayload
							p.CSMSTx <- msg
						} else {
							// call result / call error for unknown CSMS call from CS
							log.Printf("error: CS call response with message id %s has no corresponding CSMS call", msg.MessageId)
						}
					case msg := <-p.CSMSRx:
						if msg.MessageType != ocpp.MessageTypeCall {
							// call result / call error from CSMS
							if processedMessageIds.Value == msg.MessageId {
								log.Printf("warn: CSMS call response with message id %s is late", msg.MessageId)
								p.ChargeStationTx <- msg
							} else {
								log.Printf("error: CSMS message id %s is not a call", msg.MessageId)
							}
						} else {
							// call from CSMS
							if csmsCall := findCSMSCall(processedCSMSCalls, msg.MessageId); csmsCall != nil {
								log.Printf("warn: CSMS call with duplicate message id %s", msg.MessageId)
							}
							processedCSMSCalls = processedCSMSCalls.Next()
							processedCSMSCalls.Value = msg
							status = StatusCSMSCall
							p.ChargeStationTx <- msg
						}
					case msg := <-p.csmsRxCallBuf:
						// buffered call from CSMS
						if csmsCall := findCSMSCall(processedCSMSCalls, msg.MessageId); csmsCall != nil {
							log.Printf("warn: CSMS call with duplicate message id %s", msg.MessageId)
						}
						processedCSMSCalls = processedCSMSCalls.Next()
						processedCSMSCalls.Value = msg
						status = StatusCSMSCall
						p.ChargeStationTx <- msg
					case <-p.halt:
						return
					}
				}
			case StatusChargeStationCall:
				select {
				case msg := <-p.CSMSRx:
					if msg.MessageType == ocpp.MessageTypeCall {
						// call from CSMS
						select {
						case p.csmsRxCallBuf <- msg:
							log.Printf("warn: buffering CSMS call message: %s", msg.MessageId)
							continue
						default:
							log.Printf("warn: CSMS call buffer full - dropping message %s", msg.MessageId)
							continue
						}
					} else {
						// call result / call error from CSMS
						if processedMessageIds.Value == msg.MessageId {
							status = StatusWaiting
							p.ChargeStationTx <- msg
						} else {
							log.Printf("warn: CSMS call response not for current call: %s", msg.MessageId)
							continue
						}
					}
				case <-time.After(p.responseTimeout):
					log.Printf("warn: CSMS did not respond before timeout to message %s", processedMessageIds.Value)
					status = StatusWaiting
				case <-p.halt:
					return
				}
			case StatusCSMSCall:
				select {
				case msg := <-p.ChargeStationRx:
					if msg.MessageType == ocpp.MessageTypeCall {
						// call from CS
						log.Printf("warn: CS made call %s when expecting CS call response to %s", msg.MessageId, currentMsg.MessageId)
						if findMessageId(processedMessageIds, msg.MessageId) {
							log.Printf("error: CS message id %s reused", msg.MessageId)
							continue
						}
						processedMessageIds = processedMessageIds.Next()
						processedMessageIds.Value = msg.MessageId
						status = StatusChargeStationCall
						p.CSMSTx <- msg
					} else if msg.MessageId == currentMsg.MessageId {
						// call result / call error for current CSMS call from CS
						msg.Action = currentMsg.Action
						msg.RequestPayload = currentMsg.RequestPayload
						status = StatusWaiting
						p.CSMSTx <- msg
					} else if csmsCall := findCSMSCall(processedCSMSCalls, msg.MessageId); csmsCall != nil {
						// call result / call error for previous CSMS call from CS
						log.Printf("warn: CS call response to call %s when expecting response to %s", msg.MessageId, currentMsg.MessageId)
						msg.Action = csmsCall.Action
						msg.RequestPayload = csmsCall.RequestPayload
						p.CSMSTx <- msg
					} else {
						// call result / call error for unknown CSMS call from CS
						log.Printf("error: CS call response with message id %s has no corresponding CSMS call", msg.MessageId)
					}
				case <-time.After(p.responseTimeout):
					log.Printf("warn: CS did not respond before timeout to message %s", currentMsg.MessageId)
					status = StatusWaiting
				case <-p.halt:
					return
				}
			}
		}
	}()
}

func (p Pipe) Close() {
	p.halt <- struct{}{}
}

func findMessageId(in *ring.Ring, messageId string) bool {
	if in.Value == messageId {
		return true
	}
	for p := in.Next(); p != in; p = p.Next() {
		if p.Value == messageId {
			return true
		}
	}
	return false
}

func findCSMSCall(in *ring.Ring, messageId string) *GatewayMessage {
	if in.Value != nil && in.Value.(*GatewayMessage).MessageId == messageId {
		return in.Value.(*GatewayMessage)
	}
	for p := in.Next(); p != in; p = p.Next() {
		if p.Value != nil && p.Value.(*GatewayMessage).MessageId == messageId {
			return p.Value.(*GatewayMessage)
		}
	}
	return nil
}

//
//func printProcessedCSMSCalls(in *ring.Ring) {
//	buf := "["
//	if in.Value == nil {
//		buf += "nil"
//	} else {
//		buf += in.Value.(*GatewayMessage).MessageId
//	}
//	for p := in.Next(); p != in; p = p.Next() {
//		buf += ", "
//		if p.Value == nil {
//			buf += "nil"
//		} else {
//			buf += p.Value.(*GatewayMessage).MessageId
//		}
//	}
//	buf += "]"
//	log.Print(buf)
//}

//func printProcessedMessageIds(in *ring.Ring) {
//	buf := "["
//	if in.Value == nil {
//		buf += "nil"
//	} else {
//		buf += in.Value.(string)
//	}
//	for p := in.Next(); p != in; p = p.Next() {
//		buf += ", "
//		if p.Value == nil {
//			buf += "nil"
//		} else {
//			buf += p.Value.(string)
//		}
//	}
//	buf += "]"
//	log.Print(buf)
//}
