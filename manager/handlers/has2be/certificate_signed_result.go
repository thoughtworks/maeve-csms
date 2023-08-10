// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"golang.org/x/exp/slog"
)

type CertificateSignedResultHandler struct{}

func (c CertificateSignedResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	resp := response.(*has2be.CertificateSignedResponseJson)

	slog.Info("certificate signed response", slog.Any("status", resp.Status))

	return nil
}
