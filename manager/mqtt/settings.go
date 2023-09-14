// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
	"regexp"
	"strconv"
	"time"
)

func SyncSettings(ctx context.Context, engine store.Engine, clock clock.PassiveClock, v16CallMaker, v201CallMaker handlers.CallMaker, runEvery time.Duration, retryAfter time.Duration) {
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
				csId := pendingSetting.ChargeStationId
				switch details.OcppVersion {
				case "1.6":
					for name, setting := range pendingSetting.Settings {
						if setting.Status == store.ChargeStationSettingStatusPending && clock.Now().After(setting.SendAfter) {
							slog.Info("updating charge station settings", slog.String("chargeStationId", csId),
								slog.String("key", name),
								slog.String("value", setting.Value),
								slog.String("version", details.OcppVersion))
							err = engine.UpdateChargeStationSettings(ctx, csId, &store.ChargeStationSettings{
								Settings: map[string]*store.ChargeStationSetting{
									name: {Status: setting.Status, Value: setting.Value, SendAfter: clock.Now().Add(retryAfter)},
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
					var variables []ocpp201.SetVariableDataType
					for name, setting := range pendingSetting.Settings {
						slog.Info("updating charge station settings", slog.String("chargeStationId", csId),
							slog.String("key", name),
							slog.String("value", setting.Value),
							slog.String("version", details.OcppVersion))
						if setting.Status == store.ChargeStationSettingStatusPending && clock.Now().After(setting.SendAfter) {
							err = engine.UpdateChargeStationSettings(ctx, csId, &store.ChargeStationSettings{
								Settings: map[string]*store.ChargeStationSetting{
									name: {Status: setting.Status, Value: setting.Value, SendAfter: clock.Now().Add(retryAfter)},
								},
							})
							if err != nil {
								slog.Error("update charge station settings", slog.String("err", err.Error()))
								continue
							}
							var variable ocpp201.SetVariableDataType
							err = parseOcpp201Name(name, &variable)
							if err != nil {
								slog.Error("parse ocpp 2.0.1 name", slog.String("err", err.Error()))
								continue
							}
							variable.AttributeValue = setting.Value
							variables = append(variables, variable)
						}
					}
					if len(variables) > 0 {
						req := &ocpp201.SetVariablesRequestJson{
							SetVariableData: variables,
						}
						err = v201CallMaker.Send(ctx, csId, req)
						if err != nil {
							slog.Error("send set variables request", slog.String("err", err.Error()),
								slog.String("chargeStationId", csId))
						}
					}
				}
			}
		}
	}
}

// ocpp201NamePattern is a regexp that matches the following:
// - Component name - mandatory (first component)
// - Component instance - optional (first component following a ';')
// - EVSE id - optional (second component following a ';')
// - Variable name - mandatory (first component following a '/')
// - Variable instance - optional (first component following a ';')
// - Attribute type - optional (second component following a ';')
var ocpp201NamePattern = regexp.MustCompile(`^([A-Za-z0-9*\-_=:+|@.]+)(?:;([A-Za-z0-9*\-_=:+|@.]+))?(?:;(\d+))?/([A-Za-z0-9*\-_=:+|@.]+)(?:;([A-Za-z0-9*\-_=:+|@.]+))?(?:;(Actual|Target|MinSet|MaxSet))?$`)

func parseOcpp201Name(name string, set *ocpp201.SetVariableDataType) error {
	matches := ocpp201NamePattern.FindStringSubmatch(name)
	if len(matches) != 7 {
		return fmt.Errorf("invalid ocpp 2.0.1 name: %s", name)
	}
	set.Component = ocpp201.ComponentType{
		Name: matches[1],
	}
	if matches[2] != "" {
		set.Component.Instance = &matches[2]
	}
	if matches[3] != "" {
		evseId, err := strconv.Atoi(matches[3])
		if err != nil {
			return fmt.Errorf("invalid ocpp 2.0.1 name (EVSE id is not an integer): %s", name)
		}
		set.Component.Evse = &ocpp201.EVSEType{
			Id: evseId,
		}
	}

	set.Variable = ocpp201.VariableType{
		Name: matches[4],
	}
	if matches[5] != "" {
		set.Variable.Instance = &matches[5]
	}
	if matches[6] != "" {
		set.AttributeType = (*ocpp201.AttributeEnumType)(&matches[6])
	}

	return nil
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
