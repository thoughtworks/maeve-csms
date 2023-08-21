// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"
	"encoding/base64"
	"encoding/pem"

	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	types201 "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SignCertificateHandler struct {
	Handler201 handlers.CallHandler
}

func (s SignCertificateHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*typesHasToBe.SignCertificateRequestJson)

	csr, inputFormat, err := normalizeCsrEncoding(req.Csr)

	if err != nil {
		return nil, err
	}
	span.SetAttributes(attribute.String("sign_certificate.input_format", inputFormat))

	req201 := &types201.SignCertificateRequestJson{
		Csr: csr,
	}

	if req.TypeOfCertificate != nil {
		req201 = &types201.SignCertificateRequestJson{
			Csr:             csr,
			CertificateType: (*types201.CertificateSigningUseEnumType)(req.TypeOfCertificate),
		}
	}

	res, err := s.Handler201.HandleCall(ctx, chargeStationId, req201)
	if err != nil {
		return nil, err
	}
	res201 := res.(*types201.SignCertificateResponseJson)

	return &typesHasToBe.SignCertificateResponseJson{
		Status: typesHasToBe.GenericStatusEnumType(res201.Status),
	}, nil
}

func normalizeCsrEncoding(csr string) (string, string, error) {
	if pemDecoded, _ := pem.Decode([]byte(csr)); pemDecoded == nil {
		// not PEM encoded, assume base64-encoded DER
		der, err := base64.StdEncoding.DecodeString(csr)
		if err != nil {
			return "", "", err
		}
		pemBlock := pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: der,
		}
		csr = string(pem.EncodeToMemory(&pemBlock))
		return csr, "base64", nil
	}

	return csr, "pem", nil
}
