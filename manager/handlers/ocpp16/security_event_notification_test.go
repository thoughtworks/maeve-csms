package ocpp16

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"testing"
	"time"
)

func TestSecurityEventNotificationHandler(t *testing.T) {
	handler := SecurityEventNotificationHandler{}

	now := time.Now().UTC().Format(time.RFC3339)

	traceExporter := tracetest.NewInMemoryExporter()
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSyncer(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()

	func() {
		ctx, span := tracerProvider.Tracer("test").Start(ctx, `test`)
		defer span.End()

		req := &ocpp16.SecurityEventNotificationJson{
			Timestamp: now,
			Type:      "SomeSecurityEvent",
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &ocpp16.SecurityEventNotificationResponseJson{}, resp)
	}()

	require.Len(t, traceExporter.GetSpans(), 1)
	require.Len(t, traceExporter.GetSpans()[0].Attributes, 2)
	for _, attr := range traceExporter.GetSpans()[0].Attributes {
		switch attr.Key {
		case "security_event.timestamp":
			assert.Equal(t, now, attr.Value.AsString())
		case "security_event.type":
			assert.Equal(t, "SomeSecurityEvent", attr.Value.AsString())
		default:
			t.Errorf("unexpected attribute %s", attr.Key)
		}
	}
}
