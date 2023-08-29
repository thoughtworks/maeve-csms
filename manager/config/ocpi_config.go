package config

type OcpiConfig struct {
	Addr        string `mapstructure:"addr" toml:"addr"`
	ExternalURL string `mapstructure:"external_url" toml:"external_url"`
	CountryCode string `mapstructure:"country_code" toml:"country_code"`
	PartyId     string `mapstructure:"party_id" toml:"party_id"`
}
