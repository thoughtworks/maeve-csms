package ocpp201

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/twlabs/ocpp2-broker-core/manager/ocpp"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp201"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
	"log"
	"strings"
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
			certChain := addMissingNewLinesBeforeEndCertificate(*req.Certificate)
			certChain, err := removeRootCertificateIfPresent(certChain)
			if err != nil {
				return nil, fmt.Errorf("removing root certificate if present: %w", err)
			}
			_, err = a.CertificateValidationService.ValidatePEMCertificateChain([]byte(certChain), req.IdToken.IdToken)
			status, certificateStatus = handleCertificateValidationError(err)
		}

		if req.Iso15118CertificateHashData != nil {
			_, err := a.CertificateValidationService.ValidateHashedCertificateChain(*req.Iso15118CertificateHashData)
			status, certificateStatus = handleCertificateValidationError(err)
		}

		if req.IdToken.Type == types.IdTokenEnumTypeEMAID && certificateStatus == nil {
			status = types.AuthorizationStatusEnumTypeInvalid
			certStatus := types.AuthorizeCertificateStatusEnumTypeCertChainError
			certificateStatus = &certStatus
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

const pemPostCertificateBoundary = "-----END CERTIFICATE-----"

// EVerest is missing the newline before the PEM post-encapsulation boundary
// which breaks the PEM parser. The PEM parser doesn't mind if there are
// extra new lines, so we just search for the post-encapsulation boundary
// and always add in an extra line
func addMissingNewLinesBeforeEndCertificate(certificates string) string {
	certElements := strings.Split(certificates, pemPostCertificateBoundary)
	return strings.Join(certElements, "\n\n"+pemPostCertificateBoundary)
}

// EVerest is including the V2G root certificate in the chain of certificates
// which is contrary to OCPP 2.0.1 1.1.1 and breaks the validator as the root
// certificate can't be validated using OCSP so we look for any self-signed
// certificates and remove them.
func removeRootCertificateIfPresent(certificates string) (string, error) {
	certs, err := parseCertificates([]byte(certificates))
	if err != nil {
		return "", err
	}
	if len(certs) > 0 && certs[len(certs)-1].Subject.String() == certs[len(certs)-1].Issuer.String() {
		buf := bytes.NewBuffer([]byte{})
		for i := 0; i < len(certs)-1; i++ {
			err = pem.Encode(buf, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: certs[i].Raw,
			})
			if err != nil {
				return "", err
			}
		}
		return buf.String(), nil
	}
	return certificates, nil
}

func parseCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		cert, rest, err := parseCertificate(pemData)
		if err != nil {
			return nil, err
		}
		if cert == nil {
			break
		}
		certs = append(certs, cert)
		pemData = rest
	}
	return certs, nil
}

func parseCertificate(pemData []byte) (cert *x509.Certificate, rest []byte, err error) {
	block, rest := pem.Decode(pemData)
	if block == nil {
		return
	}
	if block.Type != "CERTIFICATE" {
		return
	}
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		cert = nil
		return
	}
	return
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
