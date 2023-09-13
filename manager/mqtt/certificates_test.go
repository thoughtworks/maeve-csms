// SPDX-License-Identifier: Apache-2.0

package mqtt_test

import (
	"context"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"testing"
	"time"
)

type updateChargeStation struct {
}

func (*updateChargeStation) update(ctx context.Context, engine store.Engine, chargeStationId string, request ocpp.Request) error {
	switch r := request.(type) {
	case *ocpp201.CertificateSignedRequestJson:
		var typ store.CertificateType
		switch *r.CertificateType {
		case ocpp201.CertificateSigningUseEnumTypeV2GCertificate:
			typ = store.CertificateTypeEVCC
		case ocpp201.CertificateSigningUseEnumTypeChargingStationCertificate:
			typ = store.CertificateTypeChargeStation
		}

		certificateId, err := handlers201.GetCertificateId(r.CertificateChain)
		if err != nil {
			return err
		}

		return engine.UpdateChargeStationInstallCertificates(ctx, chargeStationId, &store.ChargeStationInstallCertificates{
			Certificates: []*store.ChargeStationInstallCertificate{
				{
					CertificateType:               typ,
					CertificateId:                 certificateId,
					CertificateData:               r.CertificateChain,
					CertificateInstallationStatus: store.CertificateInstallationAccepted,
				},
			},
		})
	case *ocpp201.InstallCertificateRequestJson:
		var typ store.CertificateType
		switch r.CertificateType {
		case ocpp201.InstallCertificateUseEnumTypeV2GRootCertificate:
			typ = store.CertificateTypeV2G
		case ocpp201.InstallCertificateUseEnumTypeCSMSRootCertificate:
			typ = store.CertificateTypeCSMS
		case ocpp201.InstallCertificateUseEnumTypeMORootCertificate:
			typ = store.CertificateTypeMO
		case ocpp201.InstallCertificateUseEnumTypeManufacturerRootCertificate:
			typ = store.CertificateTypeMF
		}

		certificateId, err := handlers201.GetCertificateId(r.Certificate)
		if err != nil {
			return err
		}

		return engine.UpdateChargeStationInstallCertificates(ctx, chargeStationId, &store.ChargeStationInstallCertificates{
			Certificates: []*store.ChargeStationInstallCertificate{
				{
					CertificateType:               typ,
					CertificateId:                 certificateId,
					CertificateData:               r.Certificate,
					CertificateInstallationStatus: store.CertificateInstallationAccepted,
				},
			},
		})
	}

	return nil
}

func TestSyncCertificates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.SetChargeStationRuntimeDetails(ctx, "cs001", &store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	})
	require.NoError(t, err)
	err = engine.SetChargeStationRuntimeDetails(ctx, "cs002", &store.ChargeStationRuntimeDetails{
		OcppVersion: "2.0.1",
	})
	require.NoError(t, err)
	err = engine.SetChargeStationRuntimeDetails(ctx, "cs003", &store.ChargeStationRuntimeDetails{
		OcppVersion: "2.0.1",
	})
	require.NoError(t, err)

	evccPemBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("evcc"),
	}
	evccPemBytes := pem.EncodeToMemory(&evccPemBlock)
	evccCertId, err := handlers201.GetCertificateId(string(evccPemBytes))
	require.NoError(t, err)

	v2gPemBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("v2g"),
	}
	v2gPemBytes := pem.EncodeToMemory(&v2gPemBlock)
	v2gCertId, err := handlers201.GetCertificateId(string(v2gPemBytes))
	require.NoError(t, err)

	err = engine.UpdateChargeStationInstallCertificates(ctx, "cs001", &store.ChargeStationInstallCertificates{
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateId:                 evccCertId,
				CertificateInstallationStatus: store.CertificateInstallationPending,
				CertificateType:               store.CertificateTypeEVCC,
				CertificateData:               string(evccPemBytes),
			},
			{
				CertificateId:                 v2gCertId,
				CertificateInstallationStatus: store.CertificateInstallationPending,
				CertificateType:               store.CertificateTypeV2G,
				CertificateData:               string(v2gPemBytes),
			},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationInstallCertificates(ctx, "cs002", &store.ChargeStationInstallCertificates{
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateId:                 evccCertId,
				CertificateInstallationStatus: store.CertificateInstallationPending,
				CertificateType:               store.CertificateTypeEVCC,
				CertificateData:               string(evccPemBytes),
			},
			{
				CertificateId:                 v2gCertId,
				CertificateInstallationStatus: store.CertificateInstallationPending,
				CertificateType:               store.CertificateTypeV2G,
				CertificateData:               string(v2gPemBytes),
			},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationInstallCertificates(ctx, "cs003", &store.ChargeStationInstallCertificates{
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateId:                 v2gCertId,
				CertificateInstallationStatus: store.CertificateInstallationRejected,
				CertificateType:               store.CertificateTypeV2G,
				CertificateData:               string(v2gPemBytes),
			},
		},
	})
	require.NoError(t, err)

	updater := &updateChargeStation{}
	v16CallMaker := &mockCallMaker{
		engine:   engine,
		updateFn: updater.update,
	}
	v201CallMaker := &mockCallMaker{
		engine:   engine,
		updateFn: updater.update,
	}

	mqtt.SyncCertificates(ctx, engine, v16CallMaker, v201CallMaker, 100*time.Millisecond, 100*time.Millisecond)

	require.Len(t, v16CallMaker.callEvents, 2)
	assert.Equal(t, v16CallMaker.callEvents[0].chargeStationId, "cs001")
	assert.IsType(t, &ocpp201.CertificateSignedRequestJson{}, v16CallMaker.callEvents[0].request)
	assert.Equal(t, v16CallMaker.callEvents[1].chargeStationId, "cs001")
	assert.IsType(t, &ocpp201.InstallCertificateRequestJson{}, v16CallMaker.callEvents[1].request)

	require.Len(t, v201CallMaker.callEvents, 3)
	assert.Equal(t, v201CallMaker.callEvents[0].chargeStationId, "cs002")
	assert.IsType(t, &ocpp201.CertificateSignedRequestJson{}, v201CallMaker.callEvents[0].request)
	assert.Equal(t, v201CallMaker.callEvents[1].chargeStationId, "cs002")
	assert.IsType(t, &ocpp201.InstallCertificateRequestJson{}, v201CallMaker.callEvents[1].request)
	assert.Equal(t, v201CallMaker.callEvents[2].chargeStationId, "cs003")
	assert.IsType(t, &ocpp201.InstallCertificateRequestJson{}, v201CallMaker.callEvents[2].request)
}
