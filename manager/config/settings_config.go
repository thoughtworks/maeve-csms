// SPDX-License-Identifier: Apache-2.0

package config

type ApiSettingsConfig struct {
	Addr         string `mapstructure:"addr" toml:"addr" validate:"required"`
	ExternalAddr string `mapstructure:"external_addr,omitempty" toml:"external_addr,omitempty"`
	OrgName      string `mapstructure:"org_name,omitempty" toml:"org_name,omitempty"`
}

type OcppSettingsConfig struct {
	HeartbeatInterval string `mapstructure:"heartbeat_interval" toml:"heartbeat_interval" validate:"required"`
	Ocpp16Enabled     bool   `mapstructure:"ocpp16_enabled" toml:"ocpp16_enabled" validate:"required_without=Ocpp201Enabled"`
	Ocpp201Enabled    bool   `mapstructure:"ocpp201_enabled" toml:"ocpp201_enabled" validate:"required_without=Ocpp16Enabled"`
}

type ObservabilitySettingsConfig struct {
	LogFormat         string `mapstructure:"log_format" toml:"log_format" validate:"required"`
	OtelCollectorAddr string `mapstructure:"otel_collector_addr" toml:"otel_collector_addr"`
	TlsKeylogFile     string `mapstructure:"tls_keylog_file" toml:"tls_keylog_file"`
}
