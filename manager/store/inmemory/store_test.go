// SPDX-License-Identifier: Apache-2.0

package inmemory_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func TestUpdateChargeStationSettingsWithNewSettings(t *testing.T) {
	now := time.Now()
	engine := inmemory.NewStore(clockTest.NewFakePassiveClock(now))

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending, SendAfter: now.UTC()},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending, SendAfter: now.UTC()},
		},
	}

	err := engine.UpdateChargeStationSettings(context.Background(), "cs001", want)
	require.NoError(t, err)

	got, err := engine.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestUpdateChargeStationSettingsWithExistingSettings(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusAccepted},
		},
	}

	err := engine.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusAccepted},
		},
	})
	require.NoError(t, err)

	got, err := engine.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)

	assert.Equal(t, want.ChargeStationId, got.ChargeStationId)
	assert.Len(t, got.Settings, len(want.Settings))
	assert.Equal(t, store.ChargeStationSettingStatusPending, got.Settings["foo"].Status)
	assert.Equal(t, store.ChargeStationSettingStatusAccepted, got.Settings["baz"].Status)
}

func TestListChargeStationSettingsReturnsDataInPages(t *testing.T) {
	now := time.Now()
	engine := inmemory.NewStore(clockTest.NewFakePassiveClock(now))

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending, SendAfter: now.UTC()},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending, SendAfter: now.UTC()},
		},
	}
	for i := 0; i < 25; i++ {
		csId := fmt.Sprintf("cs%03d", i)
		err := engine.UpdateChargeStationSettings(context.Background(), csId, want)
		require.NoError(t, err)
	}

	csIds := make(map[string]struct{})

	page1, err := engine.ListChargeStationSettings(context.Background(), 10, "")
	require.NoError(t, err)
	require.Len(t, page1, 10)
	for _, got := range page1 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	page2, err := engine.ListChargeStationSettings(context.Background(), 10, page1[len(page1)-1].ChargeStationId)
	require.NoError(t, err)
	require.Len(t, page2, 10)
	for _, got := range page2 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	page3, err := engine.ListChargeStationSettings(context.Background(), 10, page2[len(page2)-1].ChargeStationId)
	require.NoError(t, err)
	require.Len(t, page3, 5)
	for _, got := range page3 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	assert.Len(t, csIds, 25)
}

func TestUpdateChargeStationInstallCertificates(t *testing.T) {
	now := time.Now()
	engine := inmemory.NewStore(clockTest.NewFakePassiveClock(now))

	want := &store.ChargeStationInstallCertificates{
		ChargeStationId: "cs001",
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               store.CertificateTypeV2G,
				CertificateId:                 "v2g001",
				CertificateData:               "v2g-pem-data",
				CertificateInstallationStatus: store.CertificateInstallationPending,
			},
		},
	}

	err := engine.UpdateChargeStationInstallCertificates(context.Background(), "cs001", want)
	require.NoError(t, err)

	got, err := engine.LookupChargeStationInstallCertificates(context.Background(), "cs001")
	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, time.Time{}, got.Certificates[0].SendAfter)
}

func TestUpdateChargeStationCertificateWithExistingCertificate(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.UpdateChargeStationInstallCertificates(context.Background(), "cs001", &store.ChargeStationInstallCertificates{
		ChargeStationId: "cs001",
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               store.CertificateTypeV2G,
				CertificateId:                 "v2g001",
				CertificateData:               "v2g-pem-data",
				CertificateInstallationStatus: store.CertificateInstallationAccepted,
			},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationInstallCertificates(context.Background(), "cs001", &store.ChargeStationInstallCertificates{
		ChargeStationId: "cs001",
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               store.CertificateTypeV2G,
				CertificateId:                 "v2g001",
				CertificateData:               "updated-v2g-pem-data",
				CertificateInstallationStatus: store.CertificateInstallationPending,
			},
		},
	})
	require.NoError(t, err)

	got, err := engine.LookupChargeStationInstallCertificates(context.Background(), "cs001")
	require.NoError(t, err)
	assert.Equal(t, "updated-v2g-pem-data", got.Certificates[0].CertificateData)
	assert.Equal(t, store.CertificateInstallationPending, got.Certificates[0].CertificateInstallationStatus)
}

func TestUpdateChargeStationCertificateWithNewCertificate(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.UpdateChargeStationInstallCertificates(context.Background(), "cs001", &store.ChargeStationInstallCertificates{
		ChargeStationId: "cs001",
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               store.CertificateTypeV2G,
				CertificateId:                 "v2g001",
				CertificateData:               "v2g-pem-data",
				CertificateInstallationStatus: store.CertificateInstallationAccepted,
			},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationInstallCertificates(context.Background(), "cs001", &store.ChargeStationInstallCertificates{
		ChargeStationId: "cs001",
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               store.CertificateTypeEVCC,
				CertificateId:                 "evcc001",
				CertificateData:               "evcc-pem-data",
				CertificateInstallationStatus: store.CertificateInstallationPending,
			},
		},
	})
	require.NoError(t, err)

	got, err := engine.LookupChargeStationInstallCertificates(context.Background(), "cs001")
	require.NoError(t, err)
	assert.Len(t, got.Certificates, 2)
	assert.Equal(t, "v2g-pem-data", got.Certificates[0].CertificateData)
	assert.Equal(t, store.CertificateInstallationAccepted, got.Certificates[0].CertificateInstallationStatus)
	assert.Equal(t, "evcc-pem-data", got.Certificates[1].CertificateData)
	assert.Equal(t, store.CertificateInstallationPending, got.Certificates[1].CertificateInstallationStatus)
}
