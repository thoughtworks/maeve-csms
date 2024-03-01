// SPDX-License-Identifier: Apache-2.0

package config

type OcspContractCertValidatorConfig struct {
	RootCertProvider RootCertProviderConfig `mapstructure:"root_certs" toml:"root_certs" validate:"required"`
	MaxAttempts      int                    `mapstructure:"max_attempts" toml:"max_attempts" validate:"required"`
}

type ContractCertValidatorConfig struct {
	Type string                           `mapstructure:"type" toml:"type" validate:"required,oneof=ocsp"`
	Ocsp *OcspContractCertValidatorConfig `mapstructure:"ocsp,omitempty" toml:"ocsp,omitempty" validate:"required_if=Type ocsp"`
}
