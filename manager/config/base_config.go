// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bufio"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
)

// BaseConfig provides the data structures that represent the configuration
// and provides the ability to load the configuration from a TOML file.
type BaseConfig struct {
	Api                       ApiSettingsConfig               `mapstructure:"api" toml:"api" validate:"required"`
	Transport                 TransportConfig                 `mapstructure:"transport" toml:"transport" validate:"required"`
	Ocpp                      OcppSettingsConfig              `mapstructure:"ocpp" toml:"ocpp" validate:"required"`
	Observability             ObservabilitySettingsConfig     `mapstructure:"observability" toml:"observability" validate:"required"`
	Storage                   StorageConfig                   `mapstructure:"storage" toml:"storage" validate:"required"`
	ContractCertValidator     ContractCertValidatorConfig     `mapstructure:"contract_cert_validator" toml:"contract_cert_validator" validate:"required"`
	ContractCertProvider      ContractCertProviderConfig      `mapstructure:"contract_cert_provider" toml:"contract_cert_provider" validate:"required"`
	ChargeStationCertProvider ChargeStationCertProviderConfig `mapstructure:"charge_station_cert_provider" toml:"charge_station_cert_provider" validate:"required"`
	TariffService             TariffServiceConfig             `mapstructure:"tariff_service" toml:"tariff_service" validate:"required"`
	Ocpi                      *OcpiConfig                     `mapstructure:"ocpi,omitempty" toml:"ocpi,omitempty"`
}

// DefaultConfig provides the default configuration. The configuration
// read from the TOML file will overlay this configuration.
var DefaultConfig = BaseConfig{
	Api: ApiSettingsConfig{
		Addr:    "localhost:9410",
		Host:    "localhost",
		WsPort:  80,
		WssPort: 443,
		OrgName: "Thoughtworks",
	},
	Transport: TransportConfig{
		Type: "mqtt",
		Mqtt: &MqttSettingsConfig{
			Urls:              []string{"mqtt://localhost:1883"},
			Prefix:            "cs",
			Group:             "manager",
			ConnectTimeout:    "10s",
			ConnectRetryDelay: "1s",
			KeepAliveInterval: "10s",
		},
	},
	Ocpp: OcppSettingsConfig{
		HeartbeatInterval: "5m",
		Ocpp16Enabled:     true,
		Ocpp201Enabled:    true,
	},
	Observability: ObservabilitySettingsConfig{
		LogFormat: "text",
	},
	Storage: StorageConfig{
		Type: "in_memory",
	},
	ContractCertValidator: ContractCertValidatorConfig{
		Type: "ocsp",
		Ocsp: &OcspContractCertValidatorConfig{
			RootCertProvider: RootCertProviderConfig{
				Type: "file",
				File: &FileRootCertProviderConfig{
					FileNames: []string{"root_ca.pem"},
				},
			},
			MaxAttempts: 1,
		},
	},
	ContractCertProvider: ContractCertProviderConfig{
		Type: "default",
	},
	ChargeStationCertProvider: ChargeStationCertProviderConfig{
		Type: "default",
	},
	TariffService: TariffServiceConfig{
		Type: "kwh",
	},
}

// Load reads TOML configuration from a reader.
func (c *BaseConfig) Load(reader io.Reader) error {
	decoder := toml.NewDecoder(reader)
	err := decoder.Decode(c)
	if err != nil {
		return err
	}
	c.replaceDefaults()
	return nil
}

// LoadFromFile reads TOML configuration from a file.
func (c *BaseConfig) LoadFromFile(configFile string) error {
	//#nosec G304 - only files specified by the person running the application will be loaded
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	err = c.Load(bufio.NewReader(f))
	c.replaceDefaults()
	return err
}

// replaceDefaults removes any default configuration (defined in DefaultConfig)
// that has been replaced by configuration loaded from the file: this happens
// when the configured `type` is different from the default
func (c *BaseConfig) replaceDefaults() {
	switch c.ContractCertValidator.Type {
	case "ocsp":
		if c.ContractCertValidator.Ocsp != nil {
			if c.ContractCertValidator.Ocsp.RootCertProvider.Type != "file" {
				c.ContractCertValidator.Ocsp.RootCertProvider.File = nil
			}
		}
	}
}

// Validate ensures that the configuration is structurally valid.
func (c *BaseConfig) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
