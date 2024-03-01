// SPDX-License-Identifier: Apache-2.0

package config

type OpcpChargeStationCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url" validate:"required"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth" validate:"required"`
}

type LocalSourceConfig struct {
	Type              string `mapstructure:"type" toml:"type" validate:"required,oneof=file google_cloud_secret"`
	File              string `mapstructure:"file,omitempty" toml:"file,omitempty" validate:"required_if=Type file"`
	GoogleCloudSecret string `mapstructure:"google_cloud_secret,omitempty" toml:"google_cloud_secret,omitempty" validate:"required_if=Type google_cloud_secret"`
}

type LocalChargeStationCertProviderConfig struct {
	CertificateSource *LocalSourceConfig `mapstructure:"cert" toml:"cert" validate:"required"`
	PrivateKeySource  *LocalSourceConfig `mapstructure:"key" toml:"key" validate:"required"`
}

type DelegatingChargeStationCertProviderConfig struct {
	V2G *ChargeStationCertProviderConfig `mapstructure:"v2g" toml:"v2g" validate:"required"`
	CSO *ChargeStationCertProviderConfig `mapstructure:"cso" toml:"cso" validate:"required"`
}

type ChargeStationCertProviderConfig struct {
	Type       string                                     `mapstructure:"type" toml:"type" validate:"required,oneof=default opcp local delegating"`
	Opcp       *OpcpChargeStationCertProviderConfig       `mapstructure:"opcp,omitempty" toml:"opcp,omitempty" validate:"required_if=Type opcp"`
	Local      *LocalChargeStationCertProviderConfig      `mapstructure:"local,omitempty" toml:"local,omitempty" validate:"required_if=Type local"`
	Delegating *DelegatingChargeStationCertProviderConfig `mapstructure:"delegating,omitempty" toml:"delegating,omitempty" validate:"required_if=Type delegating"`
}
