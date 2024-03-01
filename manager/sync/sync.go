// SPDX-License-Identifier: Apache-2.0

package sync

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/utils/clock"
	"time"
)

func Sync(storageEngine store.Engine, clock clock.PassiveClock, tracer trace.Tracer, emitter transport.Emitter) {
	v16SyncCallMaker := ocpp16.NewCallMaker(emitter)
	dataTransferCallMaker := ocpp16.NewDataTransferCallMaker(emitter)
	v201SyncCallMaker := ocpp201.NewCallMaker(emitter)

	go SyncSettings(context.Background(),
		storageEngine,
		clock,
		v16SyncCallMaker,
		v201SyncCallMaker,
		1*time.Minute,
		2*time.Minute)
	go SyncCertificates(context.Background(),
		storageEngine,
		clock,
		dataTransferCallMaker,
		v201SyncCallMaker,
		1*time.Minute,
		2*time.Minute)
	go SyncTriggers(context.Background(),
		tracer,
		storageEngine,
		clock,
		v16SyncCallMaker,
		dataTransferCallMaker,
		v201SyncCallMaker,
		1*time.Minute,
		2*time.Minute)
}
