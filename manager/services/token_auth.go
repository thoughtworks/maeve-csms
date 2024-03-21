// SPDX-License-Identifier: Apache-2.0

package services

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/utils/clock"
	"time"
)

type TokenAuthService interface {
	Authorize(ctx context.Context, token ocpp201.IdTokenType) ocpp201.IdTokenInfoType
}

type OcppTokenAuthService struct {
	TokenStore store.TokenStore
	Clock      clock.PassiveClock
}

func (o *OcppTokenAuthService) Authorize(ctx context.Context, token ocpp201.IdTokenType) ocpp201.IdTokenInfoType {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("token_auth.id", token.IdToken),
		attribute.String("token_auth.type", string(token.Type)))

	var tokenInfo *ocpp201.IdTokenInfoType

	switch token.Type {
	case ocpp201.IdTokenEnumTypeNoAuthorization:
		tokenInfo = &ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeAccepted,
		}
	case ocpp201.IdTokenEnumTypeCentral:
		tokenInfo = &ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeAccepted,
		}
	case ocpp201.IdTokenEnumTypeLocal:
		// local auth must be implemented in a different TokenAuthService
		tokenInfo = &ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeInvalid,
		}
	default:
		foundToken, err := o.TokenStore.LookupToken(ctx, token.IdToken)
		if err != nil {
			span.RecordError(err)
			tokenInfo = &ocpp201.IdTokenInfoType{
				Status: ocpp201.AuthorizationStatusEnumTypeUnknown,
			}
		} else if foundToken == nil {
			tokenInfo = &ocpp201.IdTokenInfoType{
				Status: ocpp201.AuthorizationStatusEnumTypeUnknown,
			}
		} else {
			status := ocpp201.AuthorizationStatusEnumTypeInvalid

			if foundToken.Valid {
				status = ocpp201.AuthorizationStatusEnumTypeAccepted
			}

			// if the cache mode is never, prevent the charge station
			// from caching the token by setting its expiry time to now
			var cacheExpiryTime *string
			if foundToken.CacheMode == "NEVER" {
				expiryTime := o.Clock.Now().Format(time.RFC3339)
				cacheExpiryTime = &expiryTime
			}

			var groupIdToken *ocpp201.IdTokenType
			if foundToken.GroupId != nil {
				groupIdToken = &ocpp201.IdTokenType{
					Type:    ocpp201.IdTokenEnumTypeCentral,
					IdToken: *foundToken.GroupId,
				}
				span.SetAttributes(attribute.String("token_auth.group_id", *foundToken.GroupId))
			}

			tokenInfo = &ocpp201.IdTokenInfoType{
				Status:              status,
				GroupIdToken:        groupIdToken,
				CacheExpiryDateTime: cacheExpiryTime,
			}
		}
	}

	span.SetAttributes(
		attribute.String("token_auth.status", string(tokenInfo.Status)))
	return *tokenInfo
}
