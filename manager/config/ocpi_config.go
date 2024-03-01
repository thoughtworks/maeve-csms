// SPDX-License-Identifier: Apache-2.0

package config

type OcpiConfig struct {
	Addr        string `mapstructure:"addr" toml:"addr" validate:"required"`
	ExternalURL string `mapstructure:"external_url" toml:"external_url" validate:"required"`
	CountryCode string `mapstructure:"country_code" toml:"country_code" validate:"required"`
	PartyId     string `mapstructure:"party_id" toml:"party_id" validate:"required"`
}
