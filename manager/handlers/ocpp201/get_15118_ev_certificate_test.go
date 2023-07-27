// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"testing"
)

var calledTimes int

type dummyEvCertificateProvider struct{}

func (d dummyEvCertificateProvider) ProvideCertificate(_ context.Context, exiRequest string) (services.EvCertificate15118Response, error) {
	calledTimes++
	if exiRequest == "success" {
		return services.EvCertificate15118Response{
			Status:                     types.Iso15118EVCertificateStatusEnumTypeAccepted,
			CertificateInstallationRes: "dummy exi",
		}, nil
	} else {
		return services.EvCertificate15118Response{}, errors.New("failure, try again")
	}

}

func TestGet15118EvCertificate(t *testing.T) {
	req := &types.Get15118EVCertificateRequestJson{
		Action:                types.CertificateActionEnumTypeInstall,
		Iso15118SchemaVersion: "urn:iso:15118:2:2013:MsgDef",
		ExiRequest:            "success",
	}

	h := handlers.Get15118EvCertificateHandler{
		EvCertificateProvider: dummyEvCertificateProvider{},
	}

	got, err := h.HandleCall(context.Background(), "cs001", req)
	want := &types.Get15118EVCertificateResponseJson{
		Status:      types.Iso15118EVCertificateStatusEnumTypeAccepted,
		ExiResponse: "dummy exi",
	}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
