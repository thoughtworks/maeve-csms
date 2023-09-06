// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/server"
	"go.opentelemetry.io/otel"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
	"time"
)

var (
	configFile string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long: `Starts the server which will subscribe to messages from
the gateway and send appropriate responses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig
		if configFile != "" {
			err := cfg.LoadFromFile(configFile)
			if err != nil {
				return err
			}
		}

		settings, err := config.Configure(context.Background(), &cfg)
		if err != nil {
			return err
		}
		defer func() {
			err := settings.TracerProvider.Shutdown(context.Background())
			if err != nil {
				slog.Warn("shutting down tracer provider", "error", err)
			}
		}()

		tracer := otel.Tracer("manager")

		apiServer := server.New("api", cfg.Api.Addr, nil, server.NewApiHandler(settings.Storage, settings.OcpiApi))

		mqttHandler := mqtt.NewHandler(
			mqtt.WithMqttBrokerUrls(settings.Mqtt.Urls),
			mqtt.WithMqttPrefix(settings.Mqtt.Prefix),
			mqtt.WithMqttGroup(settings.Mqtt.Group),
			mqtt.WithMqttConnectSettings(settings.Mqtt.ConnectTimeout, settings.Mqtt.ConnectRetryDelay, settings.Mqtt.KeepAliveInterval),
			mqtt.WithStorageEngine(settings.Storage),
			mqtt.WithCertValidationService(settings.ContractCertValidationService),
			mqtt.WithContractCertificateProvider(settings.ContractCertProviderService),
			mqtt.WithChargeStationCertificateProvider(settings.ChargeStationCertProviderService),
			mqtt.WithTariffService(settings.TariffService),
			mqtt.WithOtelTracer(tracer))

		errCh := make(chan error, 1)
		apiServer.Start(errCh)
		mqttHandler.Connect(errCh)

		if settings.OcpiApi != nil {
			mqttSender := mqtt.NewSender(settings.Mqtt.Urls,
				settings.Mqtt.Prefix,
				"manager",
				settings.Mqtt.ConnectTimeout,
				settings.Mqtt.ConnectRetryDelay,
				uint16(settings.Mqtt.KeepAliveInterval.Round(time.Second).Seconds()),
				tracer,
			)
			mqttSender.Connect(errCh)
			ocpiServer := server.New("ocpi", cfg.Ocpi.Addr, nil, server.NewOcpiHandler(settings.Storage, clock.RealClock{}, settings.OcpiApi, mqttSender.V16CallMaker))
			ocpiServer.Start(errCh)
		}

		err = <-errCh
		return err
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&configFile, "config-file", "c", "/config/config.toml",
		"The config file to use")
}
