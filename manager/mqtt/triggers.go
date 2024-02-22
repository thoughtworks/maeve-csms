package mqtt

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
	"time"
)

func SyncTriggers(ctx context.Context, engine store.Engine, clock clock.PassiveClock, v16CallMaker, v201CallMaker handlers.CallMaker, runEvery, retryAfter time.Duration) {
	var previousChargeStationId string
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down sync triggers")
			return
		case <-time.After(runEvery):
			slog.Info("checking for pending charge station triggers")
			triggerMessages, err := engine.ListChargeStationTriggerMessages(ctx, 50, previousChargeStationId)
			if err != nil {
				slog.Error("list charge station trigger messages", slog.String("err", err.Error()))
				continue
			}
			if len(triggerMessages) > 0 {
				previousChargeStationId = triggerMessages[len(triggerMessages)-1].ChargeStationId
			} else {
				previousChargeStationId = ""
			}
			for _, pendingTriggerMessage := range triggerMessages {
				details, err := engine.LookupChargeStationRuntimeDetails(ctx, pendingTriggerMessage.ChargeStationId)
				if err != nil {
					slog.Error("lookup charge station runtime details", slog.String("err", err.Error()),
						slog.String("chargeStationId", pendingTriggerMessage.ChargeStationId))
				}
				if details != nil {
					var callMaker handlers.CallMaker
					if details.OcppVersion == "1.6" {
						callMaker = v16CallMaker
					} else {
						callMaker = v201CallMaker
					}

					csId := pendingTriggerMessage.ChargeStationId
					if clock.Now().After(pendingTriggerMessage.SendAfter) {
						slog.Info("triggering charge station", slog.String("chargeStationId", csId),
							slog.String("trigger", string(pendingTriggerMessage.TriggerMessage)),
							slog.String("version", details.OcppVersion))
						err = engine.SetChargeStationTriggerMessage(ctx, csId, &store.ChargeStationTriggerMessage{
							TriggerMessage: pendingTriggerMessage.TriggerMessage,
							TriggerStatus:  store.TriggerStatusPending,
							SendAfter:      clock.Now().Add(retryAfter),
						})
						if err != nil {
							slog.Error("update charge station trigger message", slog.String("err", err.Error()))
							continue
						}

						err = callMaker.Send(ctx, csId, &PnCTriggerMessage{})
						if err != nil {
							slog.Error("send trigger message request", slog.String("err", err.Error()),
								slog.String("chargeStationId", csId), slog.String("trigger", string(pendingTriggerMessage.TriggerMessage)))
						}
					}
				}
			}
		}
	}
}

type PnCTriggerMessage struct{}

func (m *PnCTriggerMessage) IsRequest() {}
