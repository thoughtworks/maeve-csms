// SPDX-License-Identifier: Apache-2.0

package config

type EnvHttpTokenConfig struct {
	EnvVar string `mapstructure:"variable" toml:"variable"`
}

type FixedHttpTokenConfig struct {
	Token string `mapstructure:"token" toml:"token"`
}

type OAuth2HttpTokenConfig struct {
	Url                string  `mapstructure:"url" toml:"url"`
	ClientId           string  `mapstructure:"client_id" toml:"client_id"`
	ClientSecret       *string `mapstructure:"client_secret,omitempty" toml:"client_secret,omitempty"`
	ClientSecretEnvVar *string `mapstructure:"client_secret_env_var,omitempty" toml:"client_secret_env_var,omitempty"`
}

type HubjectTestHttpTokenConfig struct {
	Url string `mapstructure:"url" toml:"url"`
	Ttl string `mapstructure:"ttl" toml:"ttl"`
}

type HttpAuthConfig struct {
	Type             string                      `mapstructure:"type" toml:"type"`
	EnvToken         *EnvHttpTokenConfig         `mapstructure:"env_token,omitempty" toml:"env_token,omitempty"`
	FixedToken       *FixedHttpTokenConfig       `mapstructure:"fixed_token,omitempty" toml:"fixed_token,omitempty"`
	OAuth2Token      *OAuth2HttpTokenConfig      `mapstructure:"oauth2_token,omitempty" toml:"oauth2_token,omitempty"`
	HubjectTestToken *HubjectTestHttpTokenConfig `mapstructure:"hubject_test_token,omitempty" toml:"hubject_test_token,omitempty"`
}
