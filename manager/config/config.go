// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/subnova/slog-exporter/slogtrace"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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
	Addr string
}

type MqttSettings struct {
	Urls              []*url.URL
	Prefix            string
	Group             string
	ConnectTimeout    time.Duration
	ConnectRetryDelay time.Duration
	KeepAliveInterval time.Duration
}

type Config struct {
	Api                              ApiSettings
	Mqtt                             MqttSettings
	TracerProvider                   *trace.TracerProvider
	Storage                          store.Engine
	ContractCertValidationService    services.CertificateValidationService
	ContractCertProviderService      services.EvCertificateProvider
	ChargeStationCertProviderService services.CertificateSignerService
	TariffService                    services.TariffService
}

func Configure(ctx context.Context, cfg *BaseConfig) (c *Config, err error) {
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

	c = &Config{
		Api: ApiSettings{
			Addr: cfg.Api.Addr,
		},
		Mqtt: MqttSettings{
			Urls:              mqttUrls,
			Prefix:            cfg.Mqtt.Prefix,
			Group:             cfg.Mqtt.Group,
			ConnectTimeout:    mqttConnectTimeout,
			ConnectRetryDelay: mqttConnectRetryDelay,
			KeepAliveInterval: mqttKeepAliveInterval,
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

	c.Storage, err = getStorage(ctx, &cfg.Storage)
	if err != nil {
		return nil, err
	}

	c.ContractCertValidationService, err = getContractCertValidator(ctx, &cfg.ContractCertValidator, httpClient)
	if err != nil {
		return nil, err
	}

	c.ContractCertProviderService, err = getContractCertProvider(&cfg.ContractCertProvider, httpClient)
	if err != nil {
		return nil, err
	}

	c.ChargeStationCertProviderService, err = getChargeStationCertProvider(&cfg.ChargeStationCertProvider, httpClient)
	if err != nil {
		return nil, err
	}

	c.TariffService, err = getTariffService(&cfg.TariffService)
	if err != nil {
		return nil, err
	}

	return
}

func getHttpClient(keylogFile string) (*http.Client, error) {
	var transport http.RoundTripper

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

		transport = otelhttp.NewTransport(baseTransport)
	} else {
		transport = otelhttp.NewTransport(http.DefaultTransport)
	}

	return &http.Client{Transport: transport}, nil
}

func getStorage(ctx context.Context, cfg *StorageConfig) (engine store.Engine, err error) {
	switch cfg.Type {
	case "firestore":
		engine, err = firestore.NewStore(ctx, cfg.FirestoreStorage.ProjectId)
		if err != nil {
			return nil, fmt.Errorf("create firestore storage: %w", err)
		}
	case "in_memory":
		engine = inmemory.NewStore()
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}

	return
}

func getContractCertValidator(ctx context.Context, cfg *ContractCertValidatorConfig, httpClient *http.Client) (contractCertValidator services.CertificateValidationService, err error) {
	switch cfg.Type {
	case "ocsp":
		var rootCertificateProvider services.RootCertificateProviderService
		rootCertificateProvider, err = getRootCertificateProvider(&cfg.Ocsp.RootCertProvider, httpClient)
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

func getRootCertificateProvider(cfg *RootCertProviderConfig, httpClient *http.Client) (rootCertificateProvider services.RootCertificateProviderService, err error) {
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

	default:
		return nil, fmt.Errorf("unknown root certificate provider type: %s", cfg.Type)
	}

	return
}

func getContractCertProvider(cfg *ContractCertProviderConfig, httpClient *http.Client) (evCertificateProvider services.EvCertificateProvider, err error) {
	switch cfg.Type {
	case "opcp":
		httpTokenService, err := getHttpTokenService(&cfg.Opcp.HttpAuth, httpClient)
		if err != nil {
			return nil, fmt.Errorf("create http auth service: %w", err)
		}

		evCertificateProvider = &services.OpcpEvCertificateProvider{
			BaseURL:          cfg.Opcp.Url,
			HttpTokenService: httpTokenService,
			HttpClient:       httpClient,
		}
	case "default":
		evCertificateProvider = &services.DefaultEvCertificateProvider{}
	default:
		return nil, fmt.Errorf("unknown contract certificate provider type: %s", cfg.Type)
	}

	return
}

func getChargeStationCertProvider(cfg *ChargeStationCertProviderConfig, httpClient *http.Client) (chargeStationCertProvider services.CertificateSignerService, err error) {
	switch cfg.Type {
	case "opcp":
		httpTokenService, err := getHttpTokenService(&cfg.Opcp.HttpAuth, httpClient)
		if err != nil {
			return nil, fmt.Errorf("create http auth service: %w", err)
		}

		chargeStationCertProvider = &services.OpcpCpoCertificateSignerService{
			BaseURL:          cfg.Opcp.Url,
			ISOVersion:       services.ISO15118V2,
			HttpTokenService: httpTokenService,
			HttpClient:       httpClient,
		}
	case "default":
		chargeStationCertProvider = &services.DefaultCpoCertificateSignerService{}
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

func getTariffService(cfg *TariffServiceConfig) (tariffService services.TariffService, err error) {
	switch cfg.Type {
	case "kwh":
		tariffService = services.BasicKwhTariffService{}
	default:
		return nil, fmt.Errorf("unknown tariff service type: %s", cfg.Type)
	}

	return
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
