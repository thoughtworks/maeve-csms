// SPDX-License-Identifier: Apache-2.0

package config

type TariffServiceConfig struct {
	Type string `mapstructure:"type" toml:"type" validate:"required,oneof=kwh"`
}
