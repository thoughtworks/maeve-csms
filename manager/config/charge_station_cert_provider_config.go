// SPDX-License-Identifier: Apache-2.0

package config

type OpcpChargeStationCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth"`
}

type ChargeStationCertProviderConfig struct {
	Type string                               `mapstructure:"type" toml:"type"`
	Opcp *OpcpChargeStationCertProviderConfig `mapstructure:"opcp,omitempty" toml:"opcp,omitempty"`
}
