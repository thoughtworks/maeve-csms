// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/exp/slog"
)

type CertificateSignedResultHandler struct{}

func (c CertificateSignedResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	resp := response.(*ocpp201.CertificateSignedResponseJson)

	slog.Info("certificate signed response", slog.Any("status", resp.Status))

	return nil
}
