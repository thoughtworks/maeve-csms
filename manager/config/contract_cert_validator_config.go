// SPDX-License-Identifier: Apache-2.0

package config

type OcspContractCertValidatorConfig struct {
	RootCertProvider RootCertProviderConfig `mapstructure:"root_certs" toml:"root_certs"`
	MaxAttempts      int                    `mapstructure:"max_attempts" toml:"max_attempts"`
}

type ContractCertValidatorConfig struct {
	Type string                           `mapstructure:"type" toml:"type"`
	Ocsp *OcspContractCertValidatorConfig `mapstructure:"ocsp,omitempty" toml:"ocsp,omitempty"`
}
