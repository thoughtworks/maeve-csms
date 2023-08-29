// SPDX-License-Identifier: Apache-2.0

package has2be_test

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	handlersHasToBe "github.com/thoughtworks/maeve-csms/manager/handlers/has2be"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	types201 "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"testing"
)

type spy201SignCertificateHandler struct {
	RecordedCsr *string
}

func (d *spy201SignCertificateHandler) HandleCall(_ context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	d.recordReceivedCsr(request.(*types201.SignCertificateRequestJson).Csr)
	return &types201.SignCertificateResponseJson{
		Status: types201.GenericStatusEnumTypeAccepted,
	}, nil
}

func (d *spy201SignCertificateHandler) recordReceivedCsr(csr string) {
	*d.RecordedCsr = csr
}

var csrBytes = []byte("some-csr")
var pemEncodedCsr = string(pem.EncodeToMemory(&pem.Block{
	Type:  "CERTIFICATE REQUEST",
	Bytes: csrBytes,
}))

func TestPassesPEMEncodedCsrOnAsIs(t *testing.T) {
	// it would be preferable to mock the dependencies of the 201 handler instead of the handler itself
	// however, the 201 handler calls its dependencies from a go routine, but the things I tried in making the test wait
	// for the go routine to finish (like passing a waitgroup into the mocked dependencies and calling wg.Done from
	// its methods) didn't work
	csrRecorder := ""
	spy201Handler := spy201SignCertificateHandler{
		RecordedCsr: &csrRecorder,
	}
	handler := handlersHasToBe.SignCertificateHandler{
		Handler201: &spy201Handler,
	}
	req := &typesHasToBe.SignCertificateRequestJson{
		Csr: pemEncodedCsr,
	}

	_, err := handler.HandleCall(context.Background(), "cs001", req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	want := pemEncodedCsr
	got := *spy201Handler.RecordedCsr
	assert.Equal(t, want, got)
}

func TestDecodesBase64EncodedDERAndReencodesAsPEM(t *testing.T) {
	csrRecorder := ""
	spy201Handler := spy201SignCertificateHandler{
		RecordedCsr: &csrRecorder,
	}
	handler := handlersHasToBe.SignCertificateHandler{
		Handler201: &spy201Handler,
	}
	req := &typesHasToBe.SignCertificateRequestJson{
		Csr: base64.StdEncoding.EncodeToString(csrBytes),
	}

	_, err := handler.HandleCall(context.Background(), "cs001", req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	want := pemEncodedCsr
	got := *spy201Handler.RecordedCsr
	assert.Equal(t, want, got)
}
