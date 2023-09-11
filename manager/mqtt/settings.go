// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"time"
)

func SyncSettings(ctx context.Context, engine store.Engine, v16CallMaker, v201CallMaker handlers.CallMaker, runEvery time.Duration, retryAfter time.Duration) {
	var previousChargeStationId string
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down sync settings")
			return
		case <-time.After(runEvery):
			slog.Info("checking for pending charge station settings changes")
			settings, err := engine.ListChargeStationSettings(ctx, 50, previousChargeStationId)
			if err != nil {
				slog.Error("list charge station settings", slog.String("err", err.Error()))
				continue
			}
			if len(settings) > 0 {
				previousChargeStationId = settings[len(settings)-1].ChargeStationId
			} else {
				previousChargeStationId = ""
			}
			pendingSettings := filterPendingSettings(settings)
			for _, pendingSetting := range pendingSettings {
				details, err := engine.LookupChargeStationRuntimeDetails(ctx, pendingSetting.ChargeStationId)
				if err != nil {
					slog.Error("lookup charge station runtime details", slog.String("err", err.Error()),
						slog.String("chargeStationId", pendingSetting.ChargeStationId))
				}
				switch details.OcppVersion {
				case "1.6":
					csId := pendingSetting.ChargeStationId
					for name, setting := range pendingSetting.Settings {
						if setting.Status == store.ChargeStationSettingStatusPending && time.Since(setting.LastUpdated) > retryAfter {
							slog.Info("updating charge station settings", slog.String("chargeStationId", csId),
								slog.String("key", name),
								slog.String("value", setting.Value),
								slog.String("version", details.OcppVersion))
							err = engine.UpdateChargeStationSettings(ctx, csId, &store.ChargeStationSettings{
								Settings: map[string]*store.ChargeStationSetting{
									name: {Status: setting.Status, Value: setting.Value},
								},
							})
							if err != nil {
								slog.Error("update charge station settings", slog.String("err", err.Error()))
								continue
							}
							req := &ocpp16.ChangeConfigurationJson{
								Key:   name,
								Value: setting.Value,
							}
							err := v16CallMaker.Send(ctx, csId, req)
							if err != nil {
								slog.Error("send change configuration request", slog.String("err", err.Error()),
									slog.String("chargeStationId", csId), slog.String("key", name), slog.String("value", setting.Value))
							}
						}
					}
				case "2.0.1":
					slog.Warn("2.0.1 not implemented")
				}
			}
		}
	}
}

func filterPendingSettings(settings []*store.ChargeStationSettings) []*store.ChargeStationSettings {
	var pendingSettings []*store.ChargeStationSettings
	for _, setting := range settings {
		var pending bool
		for _, v := range setting.Settings {
			if v.Status == store.ChargeStationSettingStatusPending {
				pending = true
				break
			}
		}
		if pending {
			pendingSettings = append(pendingSettings, setting)
		}
	}
	return pendingSettings
}
