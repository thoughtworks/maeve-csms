// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/server"
	"go.opentelemetry.io/otel"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
	"math/rand"
	"reflect"
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var v16Emitter mqtt.Emitter
			v16Emitter = &mqtt.ProxyEmitter{}
			readyCh := make(chan struct{})
			var mqttConn *autopaho.ConnectionManager
			mqttConn, err = autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
				BrokerUrls:        settings.Mqtt.Urls,
				KeepAlive:         uint16(settings.Mqtt.KeepAliveInterval.Round(time.Second).Seconds()),
				ConnectRetryDelay: settings.Mqtt.ConnectRetryDelay,
				OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
					v16Emitter = mqtt.NewMqttEmitter(mqttConn, settings.Mqtt.Prefix, "ocpp1.6", tracer)
					readyCh <- struct{}{}
				},
				ClientConfig: paho.ClientConfig{
					ClientID: fmt.Sprintf("%s-%s", "manager", randSeq(5)),
					Router:   paho.NewStandardRouter(),
				},
			})
			if err != nil {
				slog.Error("error setting up mqttConn", "err", err)
				return err
			}

			select {
			case <-ctx.Done():
				slog.Error("timed out waiting for mqtt connection setup")
				return err
			case <-readyCh:
				// do nothing
			}

			slog.Info("setup call maker")
			v16CallMaker := mqtt.BasicCallMaker{
				E: v16Emitter,
				Actions: map[reflect.Type]string{
					reflect.TypeOf(&ocpp16.RemoteStartTransactionJson{}): "RemoteStartTransaction",
				},
			}
			ocpiServer := server.New("ocpi", cfg.Ocpi.Addr, nil, server.NewOcpiHandler(settings.Storage, clock.RealClock{}, settings.OcpiApi, v16CallMaker))
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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		//#nosec G404 - client suffix does not require secure random number generator
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
