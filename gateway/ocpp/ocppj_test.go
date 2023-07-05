package ocpp_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/twlabs/maeve-csms/gateway/ocpp"
	"testing"
)

func TestMessageCanBeMarshaledToJSON(t *testing.T) {
	msg := &ocpp.Message{
		MessageTypeId: ocpp.MessageTypeCall,
		MessageId:     "1",
		Data:          []json.RawMessage{json.RawMessage("\"ActionName\""), json.RawMessage("\"Payload\"")},
	}

	got, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshaling json: %v", err)
	}

	want := `[2,"1","ActionName","Payload"]`

	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("marshalled json: want (-) got (+)\n%s", diff)
	}
}

func TestMessageCanBeUnmarshaledFromJSON(t *testing.T) {
	in := []byte(`[2,"1","ActionName","Payload"]`)

	var got ocpp.Message
	err := json.Unmarshal(in, &got)
	if err != nil {
		t.Fatalf("unmarshaling json: %v", err)
	}

	want := ocpp.Message{
		MessageTypeId: ocpp.MessageTypeCall,
		MessageId:     "1",
		Data:          []json.RawMessage{json.RawMessage("\"ActionName\""), json.RawMessage("\"Payload\"")},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unmarshaled json: want (-) got (+)\n%s", diff)
	}
}
