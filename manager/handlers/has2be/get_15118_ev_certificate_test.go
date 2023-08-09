// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"testing"
)

var calledTimes int

type dummyContractCertificateProvider struct{}

func (d dummyContractCertificateProvider) ProvideCertificate(_ context.Context, exiRequest string) (services.EvCertificate15118Response, error) {
	calledTimes++
	if exiRequest == "success" {
		return services.EvCertificate15118Response{
			Status:                     ocpp201.Iso15118EVCertificateStatusEnumTypeAccepted,
			CertificateInstallationRes: "dummy exi",
		}, nil
	} else {
		return services.EvCertificate15118Response{}, errors.New("failure, try again")
	}

}

func TestGet15118EvCertificate(t *testing.T) {
	schemaVersion := "urn:iso:15118:2:2013:MsgDef"
	req := &typesHasToBe.Get15118EVCertificateRequestJson{
		A15118SchemaVersion: &schemaVersion,
		ExiRequest:          "success",
	}

	h := Get15118EvCertificateHandler{
		Handler201: handlers201.Get15118EvCertificateHandler{
			ContractCertificateProvider: dummyContractCertificateProvider{},
		},
	}

	got, err := h.HandleCall(context.Background(), "cs001", req)
	want := &typesHasToBe.Get15118EVCertificateResponseJson{
		Status:      typesHasToBe.Iso15118EVCertificateStatusEnumTypeAccepted,
		ExiResponse: "dummy exi",
	}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
