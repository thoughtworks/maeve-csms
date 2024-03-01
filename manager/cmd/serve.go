// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"github.com/thoughtworks/maeve-csms/manager/server"
	"github.com/thoughtworks/maeve-csms/manager/sync"
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
		settings.MsgHandler.Connect(errCh)

		if settings.OcpiApi != nil {
			ocpiServer := server.New("ocpi", cfg.Ocpi.Addr, nil, server.NewOcpiHandler(settings.Storage, clock.RealClock{}, settings.OcpiApi, settings.MsgEmitter))
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
