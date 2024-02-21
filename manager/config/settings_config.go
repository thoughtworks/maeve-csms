// SPDX-License-Identifier: Apache-2.0

package config

type ApiSettingsConfig struct {
	Addr         string `mapstructure:"addr" toml:"addr"`
	ExternalAddr string `mapstructure:"external_addr,omitempty" toml:"external_addr,omitempty"`
	OrgName      string `mapstructure:"org_name,omitempty" toml:"org_name,omitempty"`
}

type MqttSettingsConfig struct {
	Urls              []string `mapstructure:"urls" toml:"urls"`
	Prefix            string   `mapstructure:"prefix" toml:"prefix"`
	Group             string   `mapstructure:"group" toml:"group"`
	ConnectTimeout    string   `mapstructure:"connect_timeout" toml:"connect_timeout"`
	ConnectRetryDelay string   `mapstructure:"connect_retry_delay" toml:"connect_retry_delay"`
	KeepAliveInterval string   `mapstructure:"keep_alive_interval" toml:"keep_alive_interval"`
}

type ObservabilitySettingsConfig struct {
	LogFormat         string `mapstructure:"log_format" toml:"log_format"`
	OtelCollectorAddr string `mapstructure:"otel_collector_addr" toml:"otel_collector_addr"`
	TlsKeylogFile     string `mapstructure:"tls_keylog_file" toml:"tls_keylog_file"`
}
