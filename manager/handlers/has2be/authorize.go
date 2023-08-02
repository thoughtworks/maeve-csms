package has2be

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthorizeHandler struct {
	TokenStore                   store.TokenStore
	CertificateValidationService services.CertificateValidationService
}

func (a AuthorizeHandler) HandleCall(ctx context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.AuthorizeRequestJson)

	span.SetAttributes(
		attribute.String("authorize.token", req.IdToken.IdToken),
		attribute.String("authorize.token_type", string(req.IdToken.Type)))

	if req.ISO15118CertificateHashData != nil {
		span.SetAttributes(attribute.String("authorize.certificate", "hash"))
	} else {
		span.SetAttributes(attribute.String("authorize.certificate", "none"))
	}

	status := types.AuthorizationStatusEnumTypeUnknown
	tok, err := a.TokenStore.LookupToken(ctx, req.IdToken.IdToken)
	if err != nil {
		return nil, err
	}
	if tok != nil {
		status = types.AuthorizationStatusEnumTypeAccepted
	}

	certificateStatus := types.AuthorizeCertificateStatusEnumTypeAccepted
	if status == types.AuthorizationStatusEnumTypeAccepted {
		if req.ISO15118CertificateHashData != nil {
			var ocpp201CertificateHashData []ocpp201.OCSPRequestDataType

			for _, hashData := range req.ISO15118CertificateHashData {
				ocspRequestData := ocpp201.OCSPRequestDataType{
					HashAlgorithm:  ocpp201.HashAlgorithmEnumType(hashData.HashAlgorithm),
					IssuerKeyHash:  hashData.IssuerKeyHash,
					IssuerNameHash: hashData.IssuerNameHash,
					SerialNumber:   hashData.SerialNumber,
				}
				if hashData.ResponderURL != nil {
					ocspRequestData.ResponderURL = *hashData.ResponderURL
				}

				ocpp201CertificateHashData = append(ocpp201CertificateHashData, ocspRequestData)
			}

			_, err := a.CertificateValidationService.ValidateHashedCertificateChain(ctx, ocpp201CertificateHashData)

			if err != nil {
				status = types.AuthorizationStatusEnumTypeBlocked
				certificateStatus = types.AuthorizeCertificateStatusEnumTypeCertificateRevoked
				span.SetAttributes(attribute.String("authorize.cert_error", err.Error()))
			}
		}
	}

	if status != types.AuthorizationStatusEnumTypeAccepted {
		span.SetAttributes(attribute.String("authorize.cert_status", string(certificateStatus)))
	}

	span.SetAttributes(attribute.String("request.status", string(status)))

	return &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: status,
		},
		CertificateStatus: certificateStatus,
	}, nil
}
