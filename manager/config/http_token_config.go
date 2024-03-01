// SPDX-License-Identifier: Apache-2.0

package config

type EnvHttpTokenConfig struct {
	EnvVar string `mapstructure:"variable" toml:"variable" validate:"required"`
}

type FixedHttpTokenConfig struct {
	Token string `mapstructure:"token" toml:"token" validate:"required"`
}

type OAuth2HttpTokenConfig struct {
	Url                string  `mapstructure:"url" toml:"url" validate:"required"`
	ClientId           string  `mapstructure:"client_id" toml:"client_id" validate:"required"`
	ClientSecret       *string `mapstructure:"client_secret,omitempty" toml:"client_secret,omitempty" validate:"required_without=ClientSecretEnvVar"`
	ClientSecretEnvVar *string `mapstructure:"client_secret_env_var,omitempty" toml:"client_secret_env_var,omitempty" validate:"required_without=ClientSecret"`
}

type HubjectTestHttpTokenConfig struct {
	Url string `mapstructure:"url" toml:"url"`
	Ttl string `mapstructure:"ttl" toml:"ttl"`
}

type HttpAuthConfig struct {
	Type             string                      `mapstructure:"type" toml:"type" validate:"required,oneof=env_token fixed_token oauth2_token hubject_test_token"`
	EnvToken         *EnvHttpTokenConfig         `mapstructure:"env_token,omitempty" toml:"env_token,omitempty" validate:"required_if=Type env_token"`
	FixedToken       *FixedHttpTokenConfig       `mapstructure:"fixed_token,omitempty" toml:"fixed_token,omitempty" validate:"required_if=Type fixed_token"`
	OAuth2Token      *OAuth2HttpTokenConfig      `mapstructure:"oauth2_token,omitempty" toml:"oauth2_token,omitempty" validate:"required_if=Type oauth2_token"`
	HubjectTestToken *HubjectTestHttpTokenConfig `mapstructure:"hubject_test_token,omitempty" toml:"hubject_test_token,omitempty" validate:"required_if=Type hubject_test_token"`
}
