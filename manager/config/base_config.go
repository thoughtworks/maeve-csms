// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bufio"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
)

type BaseConfig struct {
	Api                       ApiSettingsConfig               `mapstructure:"api" toml:"api"`
	Mqtt                      MqttSettingsConfig              `mapstructure:"mqtt" toml:"mqtt"`
	Observability             ObservabilitySettingsConfig     `mapstructure:"observability" toml:"observability"`
	Storage                   StorageConfig                   `mapstructure:"storage" toml:"storage"`
	ContractCertValidator     ContractCertValidatorConfig     `mapstructure:"contract_cert_validator" toml:"contract_cert_validator"`
	ContractCertProvider      ContractCertProviderConfig      `mapstructure:"contract_cert_provider" toml:"contract_cert_provider"`
	ChargeStationCertProvider ChargeStationCertProviderConfig `mapstructure:"charge_station_cert_provider" toml:"charge_station_cert_provider"`
	TariffService             TariffServiceConfig             `mapstructure:"tariff_service" toml:"tariff_service"`
}

var DefaultConfig = BaseConfig{
	Api: ApiSettingsConfig{
		Addr: "localhost:9410",
	},
	Mqtt: MqttSettingsConfig{
		Urls:              []string{"mqtt://localhost:1883"},
		Prefix:            "cs",
		Group:             "manager",
		ConnectTimeout:    "10s",
		ConnectRetryDelay: "1s",
		KeepAliveInterval: "10s",
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

func (c *BaseConfig) Load(reader io.Reader) error {
	decoder := toml.NewDecoder(reader)
	err := decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *BaseConfig) LoadFromFile(configFile string) error {
	//#nosec G304 - only files specified by the person running the application will be loaded
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	return c.Load(bufio.NewReader(f))
}
