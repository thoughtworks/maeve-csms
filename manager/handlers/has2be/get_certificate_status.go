package has2be

import (
	"context"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	types201 "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

type GetCertificateStatusHandler struct {
	Handler201 handlers201.GetCertificateStatusHandler
}

func (g GetCertificateStatusHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*typesHasToBe.GetCertificateStatusRequestJson)

	ocpp201RequestData := types201.OCSPRequestDataType{
		HashAlgorithm:  types201.HashAlgorithmEnumType(req.OcspRequestData.HashAlgorithm),
		IssuerKeyHash:  req.OcspRequestData.IssuerKeyHash,
		IssuerNameHash: req.OcspRequestData.IssuerNameHash,
		SerialNumber:   req.OcspRequestData.SerialNumber,
	}

	if req.OcspRequestData.ResponderURL != nil {
		ocpp201RequestData.ResponderURL = *req.OcspRequestData.ResponderURL
	}

	req201 := &types201.GetCertificateStatusRequestJson{
		OcspRequestData: ocpp201RequestData,
	}

	res, err := g.Handler201.HandleCall(ctx, chargeStationId, req201)
	if err != nil {
		return nil, err
	}
	res201 := res.(*types201.GetCertificateStatusResponseJson)

	return &typesHasToBe.GetCertificateStatusResponseJson{
		Status:     typesHasToBe.GetCertificateStatusEnumType(res201.Status),
		OcspResult: res201.OcspResult,
	}, nil
}
