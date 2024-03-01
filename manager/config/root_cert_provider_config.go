// SPDX-License-Identifier: Apache-2.0

package config

type OpcpRootCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url" validate:"required"`
	Ttl      string         `mapstructure:"ttl" toml:"ttl" validate:"required"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth" validate:"required"`
}

type FileRootCertProviderConfig struct {
	FileNames []string `mapstructure:"files" toml:"files" validate:"required,dive,required"`
}

type CompositeRootCertProviderConfig struct {
	Providers []RootCertProviderConfig `mapstructure:"providers" toml:"providers" validate:"required,dive,required"`
}

type RootCertProviderConfig struct {
	Type      string                           `mapstructure:"type" toml:"type" validate:"required,oneof=opcp file composite"`
	Opcp      *OpcpRootCertProviderConfig      `mapstructure:"opcp,omitempty" toml:"opcp,omitempty" validate:"required_if=Type opcp"`
	File      *FileRootCertProviderConfig      `mapstructure:"file,omitempty" toml:"file,omitempty" validate:"required_if=Type file"`
	Composite *CompositeRootCertProviderConfig `mapstructure:"composite,omitempty" toml:"composite,omitempty" validate:"required_if=Type composite"`
}
