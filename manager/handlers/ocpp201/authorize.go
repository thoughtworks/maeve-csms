package ocpp201

import (
	"context"
	"errors"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"log"
)

type AuthorizeHandler struct {
	TokenStore                   services.TokenStore
	CertificateValidationService services.CertificateValidationService
}

func (a AuthorizeHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.AuthorizeRequestJson)
	log.Printf("Charge station %s authorize token %s(%s)", chargeStationId, req.IdToken.IdToken, req.IdToken.Type)

	status := types.AuthorizationStatusEnumTypeUnknown
	tok, err := a.TokenStore.FindToken(string(req.IdToken.Type), req.IdToken.IdToken)
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
		log.Printf("Charge station %s with token %s(%s) is authorized", chargeStationId, req.IdToken.IdToken, req.IdToken.Type)
	} else {
		var certStatus types.AuthorizeCertificateStatusEnumType
		if certificateStatus != nil {
			certStatus = *certificateStatus
		} else {
			certStatus = types.AuthorizeCertificateStatusEnumTypeNoCertificateAvailable
		}
		log.Printf("Charge station %s with token %s(%s) not authorized - cert status: %s",
			chargeStationId, req.IdToken.IdToken, req.IdToken.Type, certStatus)
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
		log.Printf("general validation error: %v", err)
		status = types.AuthorizationStatusEnumTypeBlocked
		certStatus = types.AuthorizeCertificateStatusEnumTypeSignatureError
	}

	return status, &certStatus
}
