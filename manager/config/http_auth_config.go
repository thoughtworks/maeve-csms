// SPDX-License-Identifier: Apache-2.0

package config

type EnvTokenHttpAuthConfig struct {
	EnvVar string `mapstructure:"variable" toml:"variable"`
}

type HubjectTestTokenHttpAuthConfig struct {
	Url string `mapstructure:"url" toml:"url"`
}

type HttpAuthConfig struct {
	Type             string                          `mapstructure:"type" toml:"type"`
	EnvToken         *EnvTokenHttpAuthConfig         `mapstructure:"env_token,omitempty" toml:"env_token,omitempty"`
	HubjectTestToken *HubjectTestTokenHttpAuthConfig `mapstructure:"hubject_test_token,omitempty" toml:"hubject_test_token,omitempty"`
}
