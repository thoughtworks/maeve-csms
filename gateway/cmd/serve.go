// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"github.com/thoughtworks/maeve-csms/gateway/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/url"
	"os"
	"time"
)

var (
	mqttAddr          string
	wsAddr            string
	wssAddr           string
	statusAddr        string
	tlsServerCert     string
	tlsServerKey      string
	tlsTrustCert      []string
	orgNames          []string
	managerApiAddr    string
	trustProxyHeaders bool
	otelCollectorAddr string
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(collectorAddr string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName("maeve-csms-gateway"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorAddr,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := trace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the gateway server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if otelCollectorAddr != "" {
			shutdown, err := initProvider(otelCollectorAddr)
			if err != nil {
				return err
			}
			defer func() {
				err := shutdown(context.Background())
				if err != nil {
					slog.Error("shutting down OTLP exporter", "error", err)
				}
			}()
		}

		tracer := otel.Tracer("gateway")

		brokerUrl, err := url.Parse(mqttAddr)
		if err != nil {
			return fmt.Errorf("parsing mqtt broker url: %v", err)
		}

		remoteRegistry := registry.RemoteRegistry{
			ManagerApiAddr: managerApiAddr,
		}
		statusServer := server.New("status", statusAddr, nil, server.NewStatusHandler())
		websocketHandler := server.NewWebsocketHandler(
			server.WithMqttBrokerUrl(brokerUrl),
			server.WithMqttTopicPrefix("cs"),
			server.WithDeviceRegistry(remoteRegistry),
			server.WithOrgNames(orgNames),
			server.WithTrustProxyHeaders(trustProxyHeaders),
			server.WithOtelTracer(tracer))
		wsServer := server.New("ws", wsAddr, nil, websocketHandler)
		var wssServer *server.Server

		if wssAddr != "" {
			if tlsServerCert == "" {
				return fmt.Errorf("no tls server cert specified for wss connection")
			}
			if tlsServerKey == "" {
				return fmt.Errorf("no tls server key specified for wss connection")
			}

			//#nosec G304 - only files specified by the person running the application will be loaded
			cb, err := os.ReadFile(tlsServerCert)
			if err != nil {
				return fmt.Errorf("reading tls cert from %s: %v", tlsServerCert, err)
			}
			//#nosec G304 - only files specified by the person running the application will be loaded
			kb, err := os.ReadFile(tlsServerKey)
			if err != nil {
				return fmt.Errorf("reading tls key from %s: %v", tlsServerKey, err)
			}
			tlsCert, err := tls.X509KeyPair(cb, kb)
			if err != nil {
				return fmt.Errorf("processing tls key pair: %v", err)
			}
			trustedCerts := x509.NewCertPool()
			for _, tc := range tlsTrustCert {
				//#nosec G304 - only files specified by the person running the application will be loaded
				tcb, err := os.ReadFile(tc)
				if err != nil {
					return fmt.Errorf("reading trusted certs from %s: %v", tc, err)
				}
				if ok := trustedCerts.AppendCertsFromPEM(tcb); !ok {
					return fmt.Errorf("processing trusted certs from %s: no certificate found", tc)
				}
			}
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{tlsCert},
				ClientCAs:    trustedCerts,
				ClientAuth:   tls.VerifyClientCertIfGiven,
				MinVersion:   tls.VersionTLS12,
			}

			wssServer = server.New("wss", wssAddr, tlsConfig, websocketHandler)
		}

		errCh := make(chan error, 1)
		wsServer.Start(errCh)
		if wssServer != nil {
			wssServer.Start(errCh)
		}
		statusServer.Start(errCh)

		err = <-errCh
		return err
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&mqttAddr, "mqtt-addr", "m", "mqtt://127.0.0.1:1883",
		"The address of the MQTT broker, e.g. mqtt://127.0.0.1:1883")
	serveCmd.Flags().StringVarP(&wsAddr, "ws-addr", "a", "127.0.0.1:9310",
		"The address that the insecure websocket server will listen on for connections, e.g. 127.0.0.1:9310")
	serveCmd.Flags().StringVarP(&wssAddr, "wss-addr", "w", "",
		"The address that the secure websocket server will listen on for connections, e.g. 127.0.0.1:9311")
	serveCmd.Flags().StringVarP(&statusAddr, "status-addr", "s", "127.0.0.1:9312",
		"The address that the status server will listen on for connections, e.g. 127.0.0.1:9312")
	serveCmd.Flags().StringVarP(&tlsServerCert, "tls-server-cert", "c", "",
		"A file that contains a PEM encoded certificate to use as the TLS server cert")
	serveCmd.Flags().StringVarP(&tlsServerKey, "tls-server-key", "k", "",
		"A file that contains a PEM encoded private key to use as the TLS server key")
	serveCmd.Flags().StringArrayVarP(&tlsTrustCert, "tls-trust-cert", "t", []string{},
		"A file that contains a PEM encoded certificate to add to the TLS trust store")
	serveCmd.Flags().StringSliceVarP(&orgNames, "org-name", "o", []string{"Thoughtworks"},
		"A comma-separated list of organisation names that are valid in client certificates")
	serveCmd.Flags().StringVarP(&managerApiAddr, "manager-api-addr", "r", "http://127.0.0.1:9410",
		"The address of the CSMS manager API, e.g. http://127.0.0.1:9410")
	serveCmd.Flags().BoolVar(&trustProxyHeaders, "trust-proxy", false,
		"Trust proxy headers when determining the client's TLS status")
	serveCmd.Flags().StringVar(&otelCollectorAddr, "otel-collector-addr", "",
		"The address of the open telemetry collector that will receive traces, e.g. localhost:4317")
}
