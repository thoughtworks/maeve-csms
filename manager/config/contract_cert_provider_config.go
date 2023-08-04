// SPDX-License-Identifier: Apache-2.0

package config

type OpcpContractCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth"`
}

type ContractCertProviderConfig struct {
	Type string                          `mapstructure:"type" toml:"type"`
	Opcp *OpcpContractCertProviderConfig `mapstructure:"opcp,omitempty" toml:"opcp,omitempty"`
}
