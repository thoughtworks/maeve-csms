// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
)

type AuthorizeHandler struct {
	TokenAuthService             services.TokenAuthService
	CertificateValidationService services.CertificateValidationService
}

func (a AuthorizeHandler) HandleCall(ctx context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.AuthorizeRequestJson)

	if req.Certificate != nil {
		span.SetAttributes(attribute.String("authorize.certificate", "chain"))
	} else if req.Iso15118CertificateHashData != nil {
		span.SetAttributes(attribute.String("authorize.certificate", "hash"))
	} else {
		span.SetAttributes(attribute.String("authorize.certificate", "none"))
	}

	idTokenInfo := a.TokenAuthService.Authorize(ctx, req.IdToken)

	var certificateStatus *types.AuthorizeCertificateStatusEnumType
	if idTokenInfo.Status == types.AuthorizationStatusEnumTypeAccepted {
		if req.Certificate != nil {
			_, err := a.CertificateValidationService.ValidatePEMCertificateChain(ctx, []byte(*req.Certificate), req.IdToken.IdToken)
			idTokenInfo.Status, certificateStatus = handleCertificateValidationError(err)
			if err != nil {
				span.SetAttributes(attribute.String("authorize.cert_error", err.Error()))
			}
		}

		if req.Iso15118CertificateHashData != nil {
			_, err := a.CertificateValidationService.ValidateHashedCertificateChain(ctx, *req.Iso15118CertificateHashData)
			idTokenInfo.Status, certificateStatus = handleCertificateValidationError(err)
			if err != nil {
				span.SetAttributes(attribute.String("authorize.cert_error", err.Error()))
			}
		}
	}

	if idTokenInfo.Status != types.AuthorizationStatusEnumTypeAccepted {
		var certStatus types.AuthorizeCertificateStatusEnumType
		if certificateStatus != nil {
			certStatus = *certificateStatus
		} else {
			certStatus = types.AuthorizeCertificateStatusEnumTypeNoCertificateAvailable
		}

		span.SetAttributes(attribute.String("authorize.cert_status", string(certStatus)))
	}

	span.SetAttributes(attribute.String("request.status", string(idTokenInfo.Status)))

	return &types.AuthorizeResponseJson{
		IdTokenInfo:       idTokenInfo,
		CertificateStatus: certificateStatus,
	}, nil
}

func handleCertificateValidationError(err error) (types.AuthorizationStatusEnumType, *types.AuthorizeCertificateStatusEnumType) {
	status := types.AuthorizationStatusEnumTypeAccepted
	certStatus := types.AuthorizeCertificateStatusEnumTypeAccepted
	var validationErr services.ValidationError
	if errors.As(err, &validationErr) {
		status = types.AuthorizationStatusEnumTypeBlocked
		switch validationErr {
		case services.ValidationErrorCertExpired:
			certStatus = types.AuthorizeCertificateStatusEnumTypeCertificateExpired
		case services.ValidationErrorCertRevoked:
			certStatus = types.AuthorizeCertificateStatusEnumTypeCertificateRevoked
		default:
			certStatus = types.AuthorizeCertificateStatusEnumTypeCertChainError
		}
	} else if err != nil {
		status = types.AuthorizationStatusEnumTypeBlocked
		certStatus = types.AuthorizeCertificateStatusEnumTypeSignatureError
	}

	return status, &certStatus
}
