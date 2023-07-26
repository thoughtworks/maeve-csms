// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/exp/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"

	"github.com/spf13/cobra"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/server"
	"github.com/thoughtworks/maeve-csms/manager/services"
)

var (
	mqttAddr                  string
	mqttPrefix                string
	mqttGroup                 string
	apiAddr                   string
	moTrustAnchorCertPEMFiles []string
	csoOPCPToken              string
	csoOPCPUrl                string
	moOPCPToken               string
	moOPCPUrl                 string
	storageEngine             string
	gcloudProject             string
	keyLogFile                string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long: `Starts the server which will subscribe to messages from
the gateway and send appropriate responses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		brokerUrl, err := url.Parse(mqttAddr)
		if err != nil {
			return fmt.Errorf("parsing mqtt broker url: %v", err)
		}

		var engine store.Engine
		switch storageEngine {
		case "firestore":
			engine, err = firestore.NewStore(context.Background(), gcloudProject)
			if err != nil {
				return err
			}
		case "inmemory":
			engine = inmemory.NewStore()
		default:
			return fmt.Errorf("unsupported storage engine %s", storageEngine)
		}

		apiServer := server.New("api", apiAddr, nil, server.NewApiHandler(engine))

		var v2gCertificates []*x509.Certificate
		for _, pemFile := range moTrustAnchorCertPEMFiles {
			parsedCerts, err := readCertificatesFromPEMFile(pemFile)
			if err != nil {
				return fmt.Errorf("reading certificates from PEM file: %s: %v", pemFile, err)
			}
			v2gCertificates = append(v2gCertificates, parsedCerts...)
		}

		tariffService := services.BasicKwhTariffService{}
		certValidationService := services.OnlineCertificateValidationService{
			RootCertificates: v2gCertificates,
			MaxOCSPAttempts:  3,
		}

		httpClient := http.DefaultClient

		if keyLogFile != "" {
			slog.Warn("***** TLS key logging enabled *****")

			//#nosec G304 - only files specified by the person running the application will be used
			keyLog, err := os.OpenFile(keyLogFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return fmt.Errorf("opening key log file: %v", err)
			}

			transport := &http.Transport{
				TLSClientConfig: &tls.Config{
					KeyLogWriter: keyLog,
					MinVersion:   tls.VersionTLS12,
				},
			}
			httpClient = &http.Client{Transport: transport}
		}

		var certSignerService services.CertificateSignerService
		var certProviderService services.EvCertificateProvider
		if csoOPCPToken != "" && csoOPCPUrl != "" {
			certSignerService = services.OpcpCpoCertificateSignerService{
				BaseURL:     csoOPCPUrl,
				BearerToken: csoOPCPToken,
				ISOVersion:  services.ISO15118V2,
				HttpClient:  httpClient,
			}
			certProviderService = services.OpcpMoEvCertificateProvider{
				BaseURL:     moOPCPUrl,
				BearerToken: moOPCPToken,
				HttpClient:  httpClient,
			}
		}

		mqttHandler := mqtt.NewHandler(
			mqtt.WithMqttBrokerUrl(brokerUrl),
			mqtt.WithMqttPrefix(mqttPrefix),
			mqtt.WithMqttGroup(mqttGroup),
			mqtt.WithTariffService(tariffService),
			mqtt.WithCertValidationService(certValidationService),
			mqtt.WithCertSignerService(certSignerService),
			mqtt.WithCertificateProviderService(certProviderService),
			mqtt.WithStorageEngine(engine),
		)

		errCh := make(chan error, 1)
		apiServer.Start(errCh)
		mqttHandler.Connect(errCh)

		err = <-errCh
		return err
	},
}

func readCertificatesFromPEMFile(pemFile string) ([]*x509.Certificate, error) {
	//#nosec G304 - only files specified by the person running the application will be loaded
	pemData, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	return parseCertificates(pemData)
}

func parseCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		cert, rest, err := parseCertificate(pemData)
		if err != nil {
			return nil, err
		}
		if cert == nil {
			break
		}
		certs = append(certs, cert)
		pemData = rest
	}
	return certs, nil
}

func parseCertificate(pemData []byte) (cert *x509.Certificate, rest []byte, err error) {
	block, rest := pem.Decode(pemData)
	if block == nil {
		return
	}
	if block.Type != "CERTIFICATE" {
		return
	}
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		cert = nil
		return
	}
	return
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&mqttAddr, "mqtt-addr", "m", "mqtt://127.0.0.1:1883",
		"The address of the MQTT broker, e.g. mqtt://127.0.0.1:1883")
	serveCmd.Flags().StringVar(&mqttPrefix, "mqtt-prefix", "cs",
		"The MQTT topic prefix that the manager will subscribe to, e.g. cs")
	serveCmd.Flags().StringVar(&mqttGroup, "mqtt-group", "manager",
		"The MQTT group to use for the shared subscription, e.g. manager")
	serveCmd.Flags().StringVarP(&apiAddr, "api-addr", "a", "127.0.0.1:9410",
		"The address that the API server will listen on for connections, e.g. 127.0.0.1:9410")
	serveCmd.Flags().StringSliceVar(&moTrustAnchorCertPEMFiles, "mo-trust-anchor-pem-file", []string{},
		"The set of PEM files containing trusted MO certificates")
	serveCmd.Flags().StringVar(&csoOPCPToken, "cso-opcp-token", "",
		"The token to use when integrating with the CSO OPCP (e.g. Hubject's token)")
	serveCmd.Flags().StringVar(&csoOPCPUrl, "cso-opcp-url", "https://open.plugncharge-test.hubject.com",
		"The Environment URL to integrate with the CSO OPCP (e.g. Hubject's environment)")
	serveCmd.Flags().StringVar(&moOPCPToken, "mo-opcp-token", "",
		"The token to use when integrating with the MO OPCP (e.g. Hubject's token)")
	serveCmd.Flags().StringVar(&moOPCPUrl, "mo-opcp-url", "https://open.plugncharge-test.hubject.com",
		"The Environment URL to integrate with the MO OPCP (e.g. Hubject's environment)")
	serveCmd.Flags().StringVarP(&storageEngine, "storage-engine", "s", "firestore",
		"The storage engine to use for persistence, one of [firestore, inmemory]")
	serveCmd.Flags().StringVar(&gcloudProject, "gcloud-project", "*detect-project-id*",
		"The google cloud project that hosts the firestore instance (if chosen storage-engine)")
	serveCmd.Flags().StringVar(&keyLogFile, "key-log-file", "",
		"File to write TLS key material to in NSS key log format (for debugging)")
}
