package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"log"
)

type CertificateSignedResultHandler struct{}

func (c CertificateSignedResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	resp := response.(*ocpp201.CertificateSignedResponseJson)

	log.Printf("Certificate signed response: %s", resp.Status)

	return nil
}
