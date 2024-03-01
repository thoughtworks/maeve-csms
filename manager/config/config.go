// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/subnova/slog-exporter/slogtrace"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	mqtt2 "github.com/thoughtworks/maeve-csms/manager/transport/mqtt"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/utils/clock"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ApiSettings struct {
	Addr         string
	ExternalAddr string
	OrgName      string
}

type Config struct {
	Api                              ApiSettings
	Tracer                           oteltrace.Tracer
	TracerProvider                   *trace.TracerProvider
	Storage                          store.Engine
	MsgEmitter                       transport.Emitter
	MsgHandler                       transport.Receiver
	ContractCertValidationService    services.CertificateValidationService
	ContractCertProviderService      services.ContractCertificateProvider
	ChargeStationCertProviderService services.ChargeStationCertificateProvider
	TariffService                    services.TariffService
	OcpiApi                          ocpi.Api
}

func Configure(ctx context.Context, cfg *BaseConfig) (c *Config, err error) {
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	heartbeatInterval, err := time.ParseDuration(cfg.Ocpp.HeartbeatInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to parse heartbeat interval: %s", err)
	}

	c = &Config{
		Api: ApiSettings{
			Addr:         cfg.Api.Addr,
			ExternalAddr: cfg.Api.ExternalAddr,
			OrgName:      cfg.Api.OrgName,
		},
	}

	switch cfg.Observability.LogFormat {
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	case "text":
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	default:
		return nil, fmt.Errorf("unknown log format: %s", cfg.Observability.LogFormat)
	}

	httpClient, err := getHttpClient(cfg.Observability.TlsKeylogFile)
	if err != nil {
		return nil, err
	}

	c.TracerProvider, err = getTracerProvider(ctx, cfg.Observability.OtelCollectorAddr)
	if err != nil {
		return nil, err
	}

	c.Tracer = c.TracerProvider.Tracer("manager")

	c.Storage, err = getStorage(ctx, &cfg.Storage)
	if err != nil {
		return nil, err
	}

	c.ContractCertValidationService, err = getContractCertValidator(&cfg.ContractCertValidator, httpClient)
	if err != nil {
		return nil, err
	}

	c.ContractCertProviderService, err = getContractCertProvider(&cfg.ContractCertProvider, httpClient)
	if err != nil {
		return nil, err
	}

	c.ChargeStationCertProviderService, err = getChargeStationCertProvider(ctx, &cfg.ChargeStationCertProvider, c.Storage, httpClient)
	if err != nil {
		return nil, err
	}

	c.TariffService, err = getTariffService(&cfg.TariffService)
	if err != nil {
		return nil, err
	}

	c.MsgEmitter, err = getMsgEmitter(&cfg.Transport, c.Tracer)
	if err != nil {
		return nil, err
	}

	var routers []transport.Router
	if cfg.Ocpp.Ocpp16Enabled {
		routers = append(routers, ocpp16.NewRouter(c.MsgEmitter,
			clock.RealClock{},
			c.Storage,
			c.ContractCertValidationService,
			c.ChargeStationCertProviderService,
			c.ContractCertProviderService,
			heartbeatInterval,
			schemas.OcppSchemas))
	}
	if cfg.Ocpp.Ocpp201Enabled {
		routers = append(routers, ocpp201.NewRouter(c.MsgEmitter,
			clock.RealClock{},
			c.Storage,
			c.TariffService,
			c.ContractCertValidationService,
			c.ChargeStationCertProviderService,
			c.ContractCertProviderService,
			heartbeatInterval,
			schemas.OcppSchemas))
	}

	c.MsgHandler, err = getMsgHandler(&cfg.Transport, c.MsgEmitter, routers, c.Tracer)
	if err != nil {
		return nil, err
	}

	if cfg.Ocpi != nil {
		c.OcpiApi, err = getOcpiApi(cfg.Ocpi, c.Storage, httpClient)
		if err != nil {
			return nil, err
		}
	}

	return
}

func getOcpiApi(o *OcpiConfig, engine store.Engine, httpClient *http.Client) (ocpi.Api, error) {
	api := ocpi.NewOCPI(engine, httpClient, o.CountryCode, o.PartyId)
	api.SetExternalUrl(o.ExternalURL)
	return api, nil
}

func getHttpClient(keylogFile string) (*http.Client, error) {
	var httpTransport http.RoundTripper

	if keylogFile != "" {
		slog.Warn("***** TLS key logging enabled *****")

		//#nosec G304 - only files specified by the person running the application will be used
		keyLog, err := os.OpenFile(keylogFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, fmt.Errorf("opening key log file: %v", err)
		}

		baseTransport := http.DefaultTransport.(*http.Transport).Clone()
		baseTransport.TLSClientConfig = &tls.Config{
			KeyLogWriter: keyLog,
			MinVersion:   tls.VersionTLS12,
		}

		httpTransport = otelhttp.NewTransport(baseTransport)
	} else {
		httpTransport = otelhttp.NewTransport(http.DefaultTransport)
	}

	return &http.Client{Transport: httpTransport}, nil
}

func getStorage(ctx context.Context, cfg *StorageConfig) (engine store.Engine, err error) {
	switch cfg.Type {
	case "firestore":
		engine, err = firestore.NewStore(ctx, cfg.FirestoreStorage.ProjectId, clock.RealClock{})
		if err != nil {
			return nil, fmt.Errorf("create firestore storage: %w", err)
		}
	case "in_memory":
		engine = inmemory.NewStore(clock.RealClock{})
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}

	return
}

func getContractCertValidator(cfg *ContractCertValidatorConfig, httpClient *http.Client) (contractCertValidator services.CertificateValidationService, err error) {
	switch cfg.Type {
	case "ocsp":
		var rootCertificateProvider services.RootCertificateProviderService
		rootCertificateProvider, err = getRootCertProvider(&cfg.Ocsp.RootCertProvider, httpClient)
		if err != nil {
			return nil, fmt.Errorf("create root certificate provider: %w", err)
		}

		contractCertValidator, err = &services.OnlineCertificateValidationService{
			RootCertificateProvider: rootCertificateProvider,
			MaxOCSPAttempts:         cfg.Ocsp.MaxAttempts,
			HttpClient:              httpClient,
		}, nil
	default:
		return nil, fmt.Errorf("unknown contract certificate validator type: %s", cfg.Type)
	}

	return
}

func getRootCertProvider(cfg *RootCertProviderConfig, httpClient *http.Client) (rootCertificateProvider services.RootCertificateProviderService, err error) {
	switch cfg.Type {
	case "file":
		rootCertificateProvider = services.FileRootCertificateProviderService{
			FilePaths: cfg.File.FileNames,
		}
	case "opcp":
		ttl := 24 * time.Hour

		if cfg.Opcp.Ttl != "" {
			ttl, err = time.ParseDuration(cfg.Opcp.Ttl)
			if err != nil {
				return nil, fmt.Errorf("failed to parse root certificate provider TTL: %w", err)
			}
		}

		httpTokenService, err := getHttpTokenService(&cfg.Opcp.HttpAuth, httpClient)
		if err != nil {
			return nil, fmt.Errorf("create http auth service: %w", err)
		}

		rootCertificateProvider = services.NewCachingRootCertificateProviderService(
			services.OpcpRootCertificateProviderService{
				BaseURL:      cfg.Opcp.Url,
				TokenService: httpTokenService,
				HttpClient:   httpClient,
			}, ttl, clock.RealClock{})
	case "composite":
		providers := make([]services.RootCertificateProviderService, len(cfg.Composite.Providers))
		for index, providerCfg := range cfg.Composite.Providers {
			providerCfg := providerCfg
			providers[index], err = getRootCertProvider(&providerCfg, httpClient)
			if err != nil {
				return nil, fmt.Errorf("creating composite root certificate provider %d: %v", index, err)
			}
		}

		rootCertificateProvider = services.CompositeRootCertificateProviderService{
			Providers: providers,
		}

	default:
		return nil, fmt.Errorf("unknown root certificate provider type: %s", cfg.Type)
	}

	return
}

func getContractCertProvider(cfg *ContractCertProviderConfig, httpClient *http.Client) (evCertificateProvider services.ContractCertificateProvider, err error) {
	switch cfg.Type {
	case "opcp":
		httpTokenService, err := getHttpTokenService(&cfg.Opcp.HttpAuth, httpClient)
		if err != nil {
			return nil, fmt.Errorf("create http auth service: %w", err)
		}

		evCertificateProvider = &services.OpcpContractCertificateProvider{
			BaseURL:          cfg.Opcp.Url,
			HttpTokenService: httpTokenService,
			HttpClient:       httpClient,
		}
	case "default":
		evCertificateProvider = &services.DefaultContractCertificateProvider{}
	default:
		return nil, fmt.Errorf("unknown contract certificate provider type: %s", cfg.Type)
	}

	return
}

func getChargeStationCertProvider(ctx context.Context, cfg *ChargeStationCertProviderConfig, engine store.Engine, httpClient *http.Client) (chargeStationCertProvider services.ChargeStationCertificateProvider, err error) {
	switch cfg.Type {
	case "opcp":
		httpTokenService, err := getHttpTokenService(&cfg.Opcp.HttpAuth, httpClient)
		if err != nil {
			return nil, fmt.Errorf("create http auth service: %w", err)
		}

		chargeStationCertProvider = &services.OpcpChargeStationCertificateProvider{
			BaseURL:          cfg.Opcp.Url,
			ISOVersion:       services.ISO15118V2,
			HttpTokenService: httpTokenService,
			HttpClient:       httpClient,
		}
	case "local":
		certificateSource, err := getLocalSource(cfg.Local.CertificateSource)
		if err != nil {
			return nil, fmt.Errorf("create local source: %w", err)
		}
		privateKeySource, err := getLocalSource(cfg.Local.PrivateKeySource)
		if err != nil {
			return nil, fmt.Errorf("create private key source: %w", err)
		}

		chargeStationCertProvider = &services.LocalChargeStationCertificateProvider{
			Store:             engine,
			CertificateReader: certificateSource,
			PrivateKeyReader:  privateKeySource,
		}
	case "delegating":
		var v2gChargeStationCertProvider services.ChargeStationCertificateProvider
		v2gChargeStationCertProvider, err = getChargeStationCertProvider(ctx, cfg.Delegating.V2G, engine, httpClient)
		if err != nil {
			return
		}
		var csoChargeStationCertProvider services.ChargeStationCertificateProvider
		csoChargeStationCertProvider, err = getChargeStationCertProvider(ctx, cfg.Delegating.CSO, engine, httpClient)
		if err != nil {
			return
		}

		chargeStationCertProvider = &services.DelegatingChargeStationCertificateProvider{
			V2GChargeStationCertificateProvider: v2gChargeStationCertProvider,
			CSOChargeStationCertificateProvider: csoChargeStationCertProvider,
		}
	case "default":
		chargeStationCertProvider = &services.DefaultChargeStationCertificateProvider{}
	default:
		return nil, fmt.Errorf("unknown charge station certificate provider type: %s", cfg.Type)
	}

	return
}

func getHttpTokenService(cfg *HttpAuthConfig, httpClient *http.Client) (httpTokenService services.HttpTokenService, err error) {
	switch cfg.Type {
	case "env_token":
		httpTokenService, err = services.NewEnvHttpTokenService(cfg.EnvToken.EnvVar)
	case "fixed_token":
		httpTokenService = services.NewFixedHttpTokenService(cfg.FixedToken.Token)
	case "oauth2_token":
		var clientSecret string
		if cfg.OAuth2Token.ClientSecret != nil {
			clientSecret = *cfg.OAuth2Token.ClientSecret
		} else if cfg.OAuth2Token.ClientSecretEnvVar != nil {
			clientSecret = os.Getenv(*cfg.OAuth2Token.ClientSecretEnvVar)
		} else {
			return nil, fmt.Errorf("client_secret or client_secret_env_var must be provided")
		}
		httpTokenService = services.NewOAuth2HttpTokenService(cfg.OAuth2Token.Url, cfg.OAuth2Token.ClientId, clientSecret, httpClient, clock.RealClock{})
	case "hubject_test_token":
		ttl := 24 * time.Hour
		if cfg.HubjectTestToken.Ttl != "" {
			ttl, err = time.ParseDuration(cfg.HubjectTestToken.Ttl)
			if err != nil {
				return nil, fmt.Errorf("parse hubject test token ttl: %w", err)
			}
		}
		httpTokenService = services.NewCachingHttpTokenService(
			services.NewHubjectTestHttpTokenService(cfg.HubjectTestToken.Url, httpClient), ttl, clock.RealClock{})
	default:
		return nil, fmt.Errorf("unknown http auth service type: %s", cfg.Type)
	}

	return
}

func getLocalSource(cfg *LocalSourceConfig) (source services.LocalSource, err error) {
	switch cfg.Type {
	case "file":
		source = services.FileSource{
			FileName: cfg.File,
		}
	case "google_cloud_secret":
		source = services.GoogleSecretSource{
			SecretName: cfg.GoogleCloudSecret,
		}
	default:
		return nil, fmt.Errorf("unknown local source type: %s", cfg.Type)
	}

	return
}

func getTariffService(cfg *TariffServiceConfig) (tariffService services.TariffService, err error) {
	switch cfg.Type {
	case "kwh":
		tariffService = services.BasicKwhTariffService{}
	default:
		return nil, fmt.Errorf("unknown tariff service type: %s", cfg.Type)
	}

	return
}

func getMsgEmitter(cfg *TransportConfig, tracer oteltrace.Tracer) (transport.Emitter, error) {
	switch cfg.Type {
	case "mqtt":
		var mqttUrls []*url.URL
		for _, urlStr := range cfg.Mqtt.Urls {
			u, err := url.Parse(urlStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mqtt url: %w", err)
			}
			mqttUrls = append(mqttUrls, u)
		}

		mqttConnectTimeout, err := time.ParseDuration(cfg.Mqtt.ConnectTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect timeout: %w", err)
		}

		mqttConnectRetryDelay, err := time.ParseDuration(cfg.Mqtt.ConnectRetryDelay)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect retry delay: %w", err)
		}

		mqttKeepAliveInterval, err := time.ParseDuration(cfg.Mqtt.KeepAliveInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt keep alive interval: %w", err)
		}

		mqttEmitter := mqtt2.NewEmitter(
			mqtt2.WithMqttBrokerUrls[mqtt2.Emitter](mqttUrls),
			mqtt2.WithMqttPrefix[mqtt2.Emitter](cfg.Mqtt.Prefix),
			mqtt2.WithMqttConnectSettings[mqtt2.Emitter](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval),
			mqtt2.WithOtelTracer[mqtt2.Emitter](tracer))

		return mqttEmitter, nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}
}

func getMsgHandler(cfg *TransportConfig, emitter transport.Emitter, routers []transport.Router, tracer oteltrace.Tracer) (transport.Receiver, error) {
	switch cfg.Type {
	case "mqtt":
		var mqttUrls []*url.URL
		for _, urlStr := range cfg.Mqtt.Urls {
			u, err := url.Parse(urlStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mqtt url: %w", err)
			}
			mqttUrls = append(mqttUrls, u)
		}

		mqttConnectTimeout, err := time.ParseDuration(cfg.Mqtt.ConnectTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect timeout: %w", err)
		}

		mqttConnectRetryDelay, err := time.ParseDuration(cfg.Mqtt.ConnectRetryDelay)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect retry delay: %w", err)
		}

		mqttKeepAliveInterval, err := time.ParseDuration(cfg.Mqtt.KeepAliveInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt keep alive interval: %w", err)
		}

		opts := []mqtt2.Opt[mqtt2.Receiver]{
			mqtt2.WithMqttBrokerUrls[mqtt2.Receiver](mqttUrls),
			mqtt2.WithMqttPrefix[mqtt2.Receiver](cfg.Mqtt.Prefix),
			mqtt2.WithMqttConnectSettings[mqtt2.Receiver](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval),
			mqtt2.WithMqttGroup(cfg.Mqtt.Group),
			mqtt2.WithEmitter(emitter),
			mqtt2.WithOtelTracer[mqtt2.Receiver](tracer),
		}

		for _, router := range routers {
			opts = append(opts, mqtt2.WithRouter(router))
		}

		mqttHandler := mqtt2.NewReceiver(opts...)

		return mqttHandler, nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}
}

func getTracerProvider(ctx context.Context, collectorAddr string) (*trace.TracerProvider, error) {
	var err error
	var res *resource.Resource
	var traceExporter trace.SpanExporter

	if collectorAddr != "" {
		res, err = resource.New(ctx,
			resource.WithDetectors(gcp.NewDetector()),
			resource.WithTelemetrySDK(),
			resource.WithAttributes(
				// the service name used to display traces in backends
				semconv.ServiceName("maeve-csms-manager"),
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
		traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}
	} else {
		res, err = resource.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource: %w", err)
		}

		traceExporter, err = slogtrace.New()
		if err != nil {
			return nil, err
		}
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

	return tracerProvider, nil
}
