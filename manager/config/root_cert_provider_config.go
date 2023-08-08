// SPDX-License-Identifier: Apache-2.0

package config

type OpcpRootCertProviderConfig struct {
	Url      string         `mapstructure:"url" toml:"url"`
	Ttl      string         `mapstructure:"ttl" toml:"ttl"`
	HttpAuth HttpAuthConfig `mapstructure:"auth" toml:"auth"`
}

type FileRootCertProviderConfig struct {
	FileNames []string `mapstructure:"files" toml:"files"`
}

type RootCertProviderConfig struct {
	Type string                      `mapstructure:"type" toml:"type"`
	Opcp *OpcpRootCertProviderConfig `mapstructure:"opcp,omitempty" toml:"opcp,omitempty"`
	File *FileRootCertProviderConfig `mapstructure:"file,omitempty" toml:"file,omitempty"`
}
