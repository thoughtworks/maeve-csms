// SPDX-License-Identifier: Apache-2.0

package pipe_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/maeve-csms/gateway/ocpp"
	"github.com/thoughtworks/maeve-csms/gateway/pipe"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func TestChargeStationCallAndResponse(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "1234",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "1234",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, callMessage, msg)
				p.CSMSRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Equal(t, callResponseMessage, msg)
				doneCh <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	// make a call from the CS
	p.ChargeStationRx <- callMessage

	select {
	case <-doneCh:
		// complete
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestCSMSCallAndResponse(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "4321",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "4321",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		Action:          "CSMSCall",
		MessageId:       "4321",
		RequestPayload:  json.RawMessage(`{"call":true}`),
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Equal(t, callMessage, msg)
				p.ChargeStationRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, csmsCallResponseMessage, msg)
				doneCh <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CSMS
	p.CSMSRx <- callMessage

	select {
	case <-doneCh:
		// complete
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestChargeStationCallMessageIdReused(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "4321",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}

	go func() {
		// incoming CSMS messages
		var found bool
		for {
			select {
			case <-p.CSMSTx:
				if found {
					time.Sleep(100 * time.Millisecond)
					assert.Fail(t, "unexpected message received at CSMS")
				}
				found = true
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Fail(t, fmt.Sprintf("unexpected message sent to charge station: %v", msg))
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CS
	p.ChargeStationRx <- callMessage
	p.ChargeStationRx <- callMessage

	<-ctx.Done()
}

func TestLateChargeStationResponseToCSMSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "9999",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "9999",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		Action:          "CSMSCall",
		MessageId:       "9999",
		RequestPayload:  json.RawMessage(`{"call":true}`),
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	go func() {
		// incoming CS messages
		for {
			select {
			case <-p.ChargeStationTx:
				time.Sleep(100 * time.Millisecond)
				p.ChargeStationRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, msg, csmsCallResponseMessage)
				doneCh <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CSMS
	p.CSMSRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestVeryLateChargeStationResponseToCSMSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "0001",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callMessage2 := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "0002",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "0001",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	callResponseMessage2 := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "0002",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		Action:          "CSMSCall",
		MessageId:       "0001",
		RequestPayload:  json.RawMessage(`{"call":true}`),
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallResponseMessage2 := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		Action:          "CSMSCall",
		MessageId:       "0002",
		RequestPayload:  json.RawMessage(`{"call":true}`),
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				time.Sleep(200 * time.Millisecond)
				if msg == callMessage {
					p.ChargeStationRx <- callResponseMessage
				} else {
					p.ChargeStationRx <- callResponseMessage2
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		// incoming CSMS messages
		var second bool

		for {
			select {
			case msg := <-p.CSMSTx:
				if !second {
					assert.Equal(t, msg, csmsCallResponseMessage)
					second = true
				} else {
					assert.Equal(t, msg, csmsCallResponseMessage2)
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CSMS
	p.CSMSRx <- callMessage
	p.CSMSRx <- callMessage2

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestChargeStationReplyToUnknownCSMSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		// incoming CSMS messages
		for {
			select {
			case <-p.CSMSTx:
				assert.Fail(t, "unexpected message received at CSMS")
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CS messages
		for {
			select {
			case <-p.ChargeStationTx:
				assert.Fail(t, "unexpected message sent to charge station")
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call response from CS
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "7654",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	p.ChargeStationRx <- callResponseMessage

	<-ctx.Done()
}

func TestCSMSCallWhilstHandlingCSCalls(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "3579",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "3579",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "6666",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}

	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, callMessage, msg)
				p.CSMSRx <- csmsCallMessage
				p.CSMSRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		var second bool
		for {
			select {
			case msg := <-p.ChargeStationTx:
				if !second {
					assert.Equal(t, callResponseMessage, msg)
					second = true
				} else {
					assert.Equal(t, csmsCallMessage, msg)
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CS
	p.ChargeStationRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestCSMSCallsWithInsufficientBufferWhilstHandlingCSCalls(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithCSMSCallQueueLen(1))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "3579",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "3579",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "6666",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	csmsCallMessage2 := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "7777",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}

	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, callMessage, msg)
				p.CSMSRx <- csmsCallMessage
				p.CSMSRx <- csmsCallMessage2
				p.CSMSRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		var second bool
		for {
			select {
			case msg := <-p.ChargeStationTx:
				if !second {
					assert.Equal(t, callResponseMessage, msg)
					second = true
				} else {
					assert.Equal(t, csmsCallMessage, msg)
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CS
	p.ChargeStationRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestCSCallWhilstHandlingCSMSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "2233",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "2233",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "4680",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		var second bool
		for {
			select {
			case msg := <-p.ChargeStationTx:
				if !second {
					assert.Equal(t, csmsCallMessage, msg)
					p.ChargeStationRx <- callMessage
					second = true
				} else {
					assert.Equal(t, callResponseMessage, msg)
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, callMessage, msg)
				p.CSMSRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	// make call from CSMS
	p.CSMSRx <- csmsCallMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestLateCSMSResponseToCSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "6543",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "6543",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Equal(t, callResponseMessage, msg)
				doneCh <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, callMessage, msg)
				time.Sleep(100 * time.Millisecond)
				p.CSMSRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	p.ChargeStationRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestVeryLateCSMSResponseToCSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "4865",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "4865",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	callResponseMessage2 := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "6789",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Equal(t, callResponseMessage, msg)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		for {
			select {
			case <-p.CSMSTx:
				time.Sleep(200 * time.Millisecond)
				p.CSMSRx <- callResponseMessage
				p.CSMSRx <- callResponseMessage2
			case <-ctx.Done():
				return
			}
		}
	}()

	p.ChargeStationRx <- callMessage

	<-ctx.Done()
}

func TestCSMSCallResponseToWrongCSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "4865",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "6789",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	callResponseMessage2 := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "4865",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Equal(t, callResponseMessage2, msg)
				doneCh <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		for {
			select {
			case <-p.CSMSTx:
				p.CSMSRx <- callResponseMessage
				p.CSMSRx <- callResponseMessage2
			case <-ctx.Done():
				return
			}
		}
	}()

	p.ChargeStationRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestCSMSDuplicateCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "5599",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "5599",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		for {
			select {
			case msg := <-p.ChargeStationTx:
				assert.Equal(t, callMessage, msg)
				p.ChargeStationRx <- callResponseMessage
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		var second bool
		for {
			select {
			case msg := <-p.CSMSTx:
				assert.Equal(t, callResponseMessage, msg)
				if !second {
					second = true
				} else {
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	p.CSMSRx <- callMessage
	p.CSMSRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestCSMSDuplicateBufferedCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe()
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSCall",
		MessageId:      "2876",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "2876",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	csmsCallMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "2000",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	csmsCallResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "2000",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		var count int
		for {
			select {
			case msg := <-p.ChargeStationTx:
				switch count {
				case 0:
					assert.Equal(t, callResponseMessage, msg)
				case 1:
					assert.Equal(t, csmsCallMessage, msg)
					p.ChargeStationRx <- csmsCallResponseMessage
				case 2:
					assert.Equal(t, csmsCallMessage, msg)
				}
				count++
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		for {
			select {
			case msg := <-p.CSMSTx:
				if msg == callMessage {
					p.CSMSRx <- csmsCallMessage
					p.CSMSRx <- csmsCallMessage
					p.CSMSRx <- callResponseMessage
				} else {
					assert.Equal(t, csmsCallResponseMessage, msg)
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	p.ChargeStationRx <- callMessage

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestChargeStationResponseForAPreviousCSMSCall(t *testing.T) {
	defer goleak.VerifyNone(t)

	p := pipe.NewPipe(pipe.WithResponseTimeout(50 * time.Millisecond))
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callMessage := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "7766",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}
	callResponseMessage := &pipe.GatewayMessage{
		MessageType:     ocpp.MessageTypeCallResult,
		MessageId:       "7766",
		ResponsePayload: json.RawMessage(`{"call":false}`),
	}
	callMessage2 := &pipe.GatewayMessage{
		MessageType:    ocpp.MessageTypeCall,
		Action:         "CSMSCall",
		MessageId:      "2555",
		RequestPayload: json.RawMessage(`{"call":true}`),
	}

	doneCh := make(chan struct{})
	go func() {
		// incoming CS messages
		var second bool
		for {
			select {
			case msg := <-p.ChargeStationTx:
				if !second {
					assert.Equal(t, callMessage, msg)
					p.ChargeStationRx <- callResponseMessage
					second = true
				} else {
					p.ChargeStationRx <- callResponseMessage
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		// incoming CSMS messages
		var second bool
		for {
			select {
			case msg := <-p.CSMSTx:
				if !second {
					assert.Equal(t, callResponseMessage, msg)
					second = true
				} else {
					assert.Equal(t, callResponseMessage, msg)
					doneCh <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	p.CSMSRx <- callMessage
	p.CSMSRx <- callMessage2

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}
