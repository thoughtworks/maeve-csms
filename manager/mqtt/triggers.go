package mqtt

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
	"time"
)

func SyncTriggers(ctx context.Context,
	tracer trace.Tracer,
	engine store.Engine,
	clock clock.PassiveClock,
	v16CallMaker,
	dataTransferCallMaker,
	v201CallMaker handlers.CallMaker,
	runEvery,
	retryAfter time.Duration) {
	var previousChargeStationId string
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down sync triggers")
			return
		case <-time.After(runEvery):
			func() {
				ctx, span := tracer.Start(ctx, "sync triggers", trace.WithSpanKind(trace.SpanKindInternal),
					trace.WithAttributes(attribute.String("sync.trigger.previous", previousChargeStationId)))
				defer span.End()
				triggerMessages, err := engine.ListChargeStationTriggerMessages(ctx, 50, previousChargeStationId)
				if err != nil {
					span.RecordError(err)
					return
				}
				if len(triggerMessages) > 0 {
					previousChargeStationId = triggerMessages[len(triggerMessages)-1].ChargeStationId
				} else {
					previousChargeStationId = ""
				}
				span.SetAttributes(attribute.Int("sync.trigger.count", len(triggerMessages)))
				for _, pendingTriggerMessage := range triggerMessages {
					func() {
						ctx, span := tracer.Start(ctx, "sync trigger", trace.WithSpanKind(trace.SpanKindInternal),
							trace.WithAttributes(
								attribute.String("chargeStationId", pendingTriggerMessage.ChargeStationId),
								attribute.String("sync.trigger.status", string(pendingTriggerMessage.TriggerStatus)),
								attribute.String("sync.trigger.message", string(pendingTriggerMessage.TriggerMessage)),
								attribute.String("sync.trigger.after", pendingTriggerMessage.SendAfter.Format(time.RFC3339)),
							))
						defer span.End()
						details, err := engine.LookupChargeStationRuntimeDetails(ctx, pendingTriggerMessage.ChargeStationId)
						if err != nil {
							span.RecordError(err)
							return
						}
						if details == nil {
							span.RecordError(fmt.Errorf("no runtime details for charge station"))
							return
						}

						csId := pendingTriggerMessage.ChargeStationId
						if clock.Now().After(pendingTriggerMessage.SendAfter) {
							span.SetAttributes(attribute.String("sync.trigger.ocpp_version", details.OcppVersion))
							err = engine.SetChargeStationTriggerMessage(ctx, csId, &store.ChargeStationTriggerMessage{
								TriggerMessage: pendingTriggerMessage.TriggerMessage,
								TriggerStatus:  store.TriggerStatusPending,
								SendAfter:      clock.Now().Add(retryAfter),
							})
							if err != nil {
								span.RecordError(err)
								return
							}

							if details.OcppVersion == "1.6" {
								if pendingTriggerMessage.TriggerMessage == store.TriggerMessageBootNotification ||
									pendingTriggerMessage.TriggerMessage == store.TriggerMessageDiagnosticStatusNotification ||
									pendingTriggerMessage.TriggerMessage == store.TriggerMessageFirmwareStatusNotification ||
									pendingTriggerMessage.TriggerMessage == store.TriggerMessageHeartbeat ||
									pendingTriggerMessage.TriggerMessage == store.TriggerMessageMeterValues ||
									pendingTriggerMessage.TriggerMessage == store.TriggerMessageStatusNotification {
									err = v16CallMaker.Send(ctx, csId, &ocpp16.TriggerMessageJson{
										RequestedMessage: ocpp16.TriggerMessageJsonRequestedMessage(pendingTriggerMessage.TriggerMessage),
									})
								} else {
									err = dataTransferCallMaker.Send(ctx, csId, &ocpp201.TriggerMessageRequestJson{
										RequestedMessage: ocpp201.MessageTriggerEnumType(pendingTriggerMessage.TriggerMessage),
									})
								}
							} else {
								err = v201CallMaker.Send(ctx, csId, &ocpp201.TriggerMessageRequestJson{
									RequestedMessage: ocpp201.MessageTriggerEnumType(pendingTriggerMessage.TriggerMessage),
								})
							}

							if err != nil {
								span.RecordError(err)
							}
						}
					}()
				}
			}()
		}
	}
}
