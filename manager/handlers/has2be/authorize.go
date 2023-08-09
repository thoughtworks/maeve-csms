package has2be

import (
	"context"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	types201 "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

type AuthorizeHandler struct {
	Handler201 handlers201.AuthorizeHandler
}

func (a AuthorizeHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*typesHasToBe.AuthorizeRequestJson)

	req201 := &types201.AuthorizeRequestJson{
		IdToken: types201.IdTokenType{
			IdToken: req.IdToken.IdToken,
			Type:    types201.IdTokenEnumType(req.IdToken.Type),
		},
	}

	if req.ISO15118CertificateHashData != nil {
		var ocpp201CertificateHashData []types201.OCSPRequestDataType

		for _, hashData := range req.ISO15118CertificateHashData {
			ocspRequestData := types201.OCSPRequestDataType{
				HashAlgorithm:  types201.HashAlgorithmEnumType(hashData.HashAlgorithm),
				IssuerKeyHash:  hashData.IssuerKeyHash,
				IssuerNameHash: hashData.IssuerNameHash,
				SerialNumber:   hashData.SerialNumber,
			}
			if hashData.ResponderURL != nil {
				ocspRequestData.ResponderURL = *hashData.ResponderURL
			}

			ocpp201CertificateHashData = append(ocpp201CertificateHashData, ocspRequestData)
		}

		req201 = &types201.AuthorizeRequestJson{
			IdToken: types201.IdTokenType{
				IdToken: req.IdToken.IdToken,
				Type:    types201.IdTokenEnumType(req.IdToken.Type),
			},
			Iso15118CertificateHashData: &ocpp201CertificateHashData,
		}
	}

	certificateStatus := typesHasToBe.AuthorizeCertificateStatusEnumTypeAccepted
	res, err := a.Handler201.HandleCall(ctx, chargeStationId, req201)
	if err != nil {
		return nil, err
	}
	res201 := res.(*types201.AuthorizeResponseJson)
	if *res201.CertificateStatus != types201.AuthorizeCertificateStatusEnumTypeAccepted {
		certificateStatus = typesHasToBe.AuthorizeCertificateStatusEnumTypeCertificateRevoked
	}

	return &typesHasToBe.AuthorizeResponseJson{
		IdTokenInfo: typesHasToBe.IdTokenInfoType{
			Status:              typesHasToBe.AuthorizationStatusEnumType(res201.IdTokenInfo.Status),
			CacheExpiryDateTime: res201.IdTokenInfo.CacheExpiryDateTime,
		},
		CertificateStatus: certificateStatus,
	}, nil
}
