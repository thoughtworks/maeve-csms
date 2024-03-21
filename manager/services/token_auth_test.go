// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	fakeclock "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func TestOcppTokenAuthServiceAcceptsNoAuthorization(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeNoAuthorization,
			IdToken: "",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeAccepted,
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "NoAuthorization",
		"token_auth.id":     "",
		"token_auth.status": "Accepted",
	})
}

func TestOcppTokenAuthServiceAcceptsCentral(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeCentral,
			IdToken: "SomeToken",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeAccepted,
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "Central",
		"token_auth.id":     "SomeToken",
		"token_auth.status": "Accepted",
	})
}

func TestOcppTokenAuthServiceRejectsLocalAuth(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeLocal,
			IdToken: "some-local-id",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeInvalid,
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "Local",
		"token_auth.id":     "some-local-id",
		"token_auth.status": "Invalid",
	})
}

func TestOcppTokenAuthServiceReturnsUnknownIfNoTokenRegistered(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeISO14443,
			IdToken: "DEADBEEF",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeUnknown,
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "ISO14443",
		"token_auth.id":     "DEADBEEF",
		"token_auth.status": "Unknown",
	})
}

func TestOcppTokenAuthServiceReturnsInvalidIfTokenIsInvalid(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	err := tokenStore.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "TWKABC1234",
		Issuer:      "Thoughtworks",
		GroupId:     nil,
		Valid:       false,
		CacheMode:   "ALWAYS",
	})
	require.NoError(t, err)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeISO14443,
			IdToken: "DEADBEEF",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeInvalid,
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "ISO14443",
		"token_auth.id":     "DEADBEEF",
		"token_auth.status": "Invalid",
	})
}

func TestOcppTokenAuthServiceReturnsAcceptedIfTokenIsValid(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	err := tokenStore.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "TWKABC1234",
		Issuer:      "Thoughtworks",
		GroupId:     nil,
		Valid:       true,
		CacheMode:   "ALWAYS",
	})
	require.NoError(t, err)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeISO14443,
			IdToken: "DEADBEEF",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeAccepted,
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "ISO14443",
		"token_auth.id":     "DEADBEEF",
		"token_auth.status": "Accepted",
	})
}

func TestOcppTokenAuthServiceIncludesGroupIdIfConfigured(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	err := tokenStore.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "TWKABC1234",
		Issuer:      "Thoughtworks",
		GroupId:     makePtr("SOMEGROUP"),
		Valid:       true,
		CacheMode:   "ALWAYS",
	})
	require.NoError(t, err)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeISO14443,
			IdToken: "DEADBEEF",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status: ocpp201.AuthorizationStatusEnumTypeAccepted,
			GroupIdToken: &ocpp201.IdTokenType{
				Type:    ocpp201.IdTokenEnumTypeCentral,
				IdToken: "SOMEGROUP",
			},
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":     "ISO14443",
		"token_auth.id":       "DEADBEEF",
		"token_auth.status":   "Accepted",
		"token_auth.group_id": "SOMEGROUP",
	})
}

func TestOcppTokenAuthServiceSetsExpiryTimeToNowIfCacheModeIsNever(t *testing.T) {
	now := time.Now()
	clock := fakeclock.NewFakePassiveClock(now)
	tokenStore := inmemory.NewStore(clock)

	err := tokenStore.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "TWKABC1234",
		Issuer:      "Thoughtworks",
		GroupId:     nil,
		Valid:       true,
		CacheMode:   "NEVER",
	})
	require.NoError(t, err)

	tokenAuthService := services.OcppTokenAuthService{
		TokenStore: tokenStore,
		Clock:      clock,
	}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		tokenInfo := tokenAuthService.Authorize(ctx, ocpp201.IdTokenType{
			Type:    ocpp201.IdTokenEnumTypeISO14443,
			IdToken: "DEADBEEF",
		})

		assert.Equal(t, ocpp201.IdTokenInfoType{
			Status:              ocpp201.AuthorizationStatusEnumTypeAccepted,
			CacheExpiryDateTime: makePtr(now.Format(time.RFC3339)),
		}, tokenInfo)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"token_auth.type":   "ISO14443",
		"token_auth.id":     "DEADBEEF",
		"token_auth.status": "Accepted",
	})
}
