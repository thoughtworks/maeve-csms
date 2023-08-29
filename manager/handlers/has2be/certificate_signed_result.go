// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
)

type CertificateSignedResultHandler struct{}

func (c CertificateSignedResultHandler) HandleCallResult(ctx context.Context, _ string, _ ocpp.Request, response ocpp.Response, _ any) error {
	span := trace.SpanFromContext(ctx)

	resp := response.(*has2be.CertificateSignedResponseJson)

	span.SetAttributes(attribute.String("response.status", string(resp.Status)))

	return nil
}
