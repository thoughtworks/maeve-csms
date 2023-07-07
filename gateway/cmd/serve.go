// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"github.com/thoughtworks/maeve-csms/gateway/server"
	"net/url"
	"os"
)

var mqttAddr string
var wsAddr string
var wssAddr string
var statusAddr string
var tlsServerCert string
var tlsServerKey string
var tlsTrustCert []string
var orgNames []string
var csPassword string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the gateway server",
	RunE: func(cmd *cobra.Command, args []string) error {
		brokerUrl, err := url.Parse(mqttAddr)
		if err != nil {
			return fmt.Errorf("parsing mqtt broker url: %v", err)
		}

		sha256pw := sha256.Sum256([]byte(csPassword))
		mockRegistry := registry.NewMockRegistry()
		mockRegistry.ChargeStations["cp001"] = &registry.ChargeStation{
			ClientId:             "cp001",
			SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
			Base64SHA256Password: base64.StdEncoding.EncodeToString(sha256pw[:]),
		}
		mockRegistry.ChargeStations["cs001"] = &registry.ChargeStation{
			ClientId:        "cs001",
			SecurityProfile: registry.TLSWithClientSideCertificates,
		}
		mockRegistry.ChargeStations["cs002"] = &registry.ChargeStation{
			ClientId:        "cs002",
			SecurityProfile: registry.TLSWithClientSideCertificates,
		}
		mockRegistry.ChargeStations["cs003"] = &registry.ChargeStation{
			ClientId:        "cs003",
			SecurityProfile: registry.TLSWithClientSideCertificates,
		}

		statusServer := server.New("status", statusAddr, nil, server.NewStatusHandler())
		websocketHandler := server.NewWebsocketHandler(
			server.WithMqttBrokerUrl(brokerUrl),
			server.WithMqttTopicPrefix("cs"),
			server.WithDeviceRegistry(mockRegistry),
			server.WithOrgNames(orgNames))
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
	serveCmd.Flags().StringVar(&csPassword, "cs-password", "",
		"The password to use for the charge station")
}
