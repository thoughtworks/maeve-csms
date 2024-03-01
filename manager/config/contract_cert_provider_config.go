// SPDX-License-Identifier: Apache-2.0

package config

type OpcpContractCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url" validate:"required"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth" validate:"required"`
}

type ContractCertProviderConfig struct {
	Type string                          `mapstructure:"type" toml:"type" validate:"required,oneof=default opcp"`
	Opcp *OpcpContractCertProviderConfig `mapstructure:"opcp,omitempty" toml:"opcp,omitempty" validate:"required_if=Type opcp"`
}
