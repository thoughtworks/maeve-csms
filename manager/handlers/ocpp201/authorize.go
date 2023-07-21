// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"errors"
	"fmt"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
)

type AuthorizeHandler struct {
	TokenStore                   store.TokenStore
	CertificateValidationService services.CertificateValidationService
}

func (a AuthorizeHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.AuthorizeRequestJson)
	slog.Info("checking", slog.String("chargeStationId", chargeStationId),
		slog.String("idToken", req.IdToken.IdToken),
		slog.String("idTokenType", string(req.IdToken.Type)))
	status := types.AuthorizationStatusEnumTypeUnknown
	tok, err := a.TokenStore.LookupToken(ctx, req.IdToken.IdToken)
	if err != nil {
		return nil, err
	}
	if tok != nil {
		status = types.AuthorizationStatusEnumTypeAccepted
	}

	var certificateStatus *types.AuthorizeCertificateStatusEnumType
	if status == types.AuthorizationStatusEnumTypeAccepted {
		if req.Certificate != nil {
			if err != nil {
				return nil, fmt.Errorf("removing root certificate if present: %w", err)
			}
			_, err = a.CertificateValidationService.ValidatePEMCertificateChain([]byte(*req.Certificate), req.IdToken.IdToken)
			status, certificateStatus = handleCertificateValidationError(err)
		}

		if req.Iso15118CertificateHashData != nil {
			_, err := a.CertificateValidationService.ValidateHashedCertificateChain(*req.Iso15118CertificateHashData)
			status, certificateStatus = handleCertificateValidationError(err)
		}
	}

	if status == types.AuthorizationStatusEnumTypeAccepted {
		slog.Info("charge station authorized", slog.String("chargeStationId", chargeStationId),
			slog.String("idToken", req.IdToken.IdToken),
			slog.String("type", string(req.IdToken.Type)))
	} else {
		var certStatus types.AuthorizeCertificateStatusEnumType
		if certificateStatus != nil {
			certStatus = *certificateStatus
		} else {
			certStatus = types.AuthorizeCertificateStatusEnumTypeNoCertificateAvailable
		}
		slog.Warn("charge station not authorized",
			slog.String("chargeStationId", chargeStationId),
			slog.String("idToken", req.IdToken.IdToken),
			slog.String("type", string(req.IdToken.Type)),
			slog.String("certStatus", string(certStatus)))
	}

	return &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: status,
		},
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
		slog.Error("general validation error", err)
		status = types.AuthorizationStatusEnumTypeBlocked
		certStatus = types.AuthorizeCertificateStatusEnumTypeSignatureError
	}

	return status, &certStatus
}
