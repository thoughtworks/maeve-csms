// SPDX-License-Identifier: Apache-2.0

package config

type OpcpChargeStationCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth"`
}

type LocalSourceConfig struct {
	Type              string `mapstructure:"type" toml:"type"`
	File              string `mapstructure:"file,omitempty" toml:"file,omitempty"`
	GoogleCloudSecret string `mapstructure:"google_cloud_secret,omitempty" toml:"google_cloud_secret,omitempty"`
}

type LocalChargeStationCertProviderConfig struct {
	CertificateSource *LocalSourceConfig `mapstructure:"cert" toml:"cert"`
	PrivateKeySource  *LocalSourceConfig `mapstructure:"key" toml:"key"`
}

type DelegatingChargeStationCertProviderConfig struct {
	V2G *ChargeStationCertProviderConfig `mapstructure:"v2g" toml:"v2g"`
	CSO *ChargeStationCertProviderConfig `mapstructure:"cso" toml:"csp"`
}

type ChargeStationCertProviderConfig struct {
	Type       string                                     `mapstructure:"type" toml:"type"`
	Opcp       *OpcpChargeStationCertProviderConfig       `mapstructure:"opcp,omitempty" toml:"opcp,omitempty"`
	Local      *LocalChargeStationCertProviderConfig      `mapstructure:"local,omitempty" toml:"local,omitempty"`
	Delegating *DelegatingChargeStationCertProviderConfig `mapstructure:"delegating,omitempty" toml:"delegating,omitempty"`
}
