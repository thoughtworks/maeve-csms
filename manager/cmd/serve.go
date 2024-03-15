// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"github.com/thoughtworks/maeve-csms/manager/server"
	"github.com/thoughtworks/maeve-csms/manager/sync"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
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

		apiServer := server.New("api", cfg.Api.Addr, nil,
			server.NewApiHandler(settings.Api, settings.Storage, settings.OcpiApi, settings.ChargeStationCertProviderService))

		sync.Sync(settings.Storage, clock.RealClock{}, settings.Tracer, settings.MsgEmitter)

		errCh := make(chan error, 1)
		apiServer.Start(errCh)
		var ocpp16Connection transport.Connection
		if settings.Ocpp16Handler != nil {
			ocpp16Connection, err = settings.MsgListener.Connect(context.Background(), transport.OcppVersion16, nil, settings.Ocpp16Handler)
			if err != nil {
				errCh <- err
			}
		}

		var ocpp201Connection transport.Connection
		if settings.Ocpp201Handler != nil {
			ocpp201Connection, err = settings.MsgListener.Connect(context.Background(), transport.OcppVersion201, nil, settings.Ocpp201Handler)
			if err != nil {
				errCh <- err
			}
		}

		if settings.OcpiApi != nil {
			ocpiServer := server.New("ocpi", cfg.Ocpi.Addr, nil, server.NewOcpiHandler(settings.Storage, clock.RealClock{}, settings.OcpiApi, settings.MsgEmitter))
			ocpiServer.Start(errCh)
		}

		err = <-errCh

		if ocpp16Connection != nil {
			err := ocpp16Connection.Disconnect(context.Background())
			if err != nil {
				slog.Warn("disconnecting from broker", "err", err)
			}
		}
		if ocpp201Connection != nil {
			err := ocpp201Connection.Disconnect(context.Background())
			if err != nil {
				slog.Warn("disconnecting from broker", "err", err)
			}
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&configFile, "config-file", "c", "/config/config.toml",
		"The config file to use")
}
