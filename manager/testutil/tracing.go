// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/maps"
	"sort"
	"testing"
)

// GetTracer is used in tests to get a trace.Tracer implementation that captures
// the spans generated during a test.
func GetTracer() (trace.Tracer, *tracetest.InMemoryExporter) {
	traceExporter := tracetest.NewInMemoryExporter()
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSyncer(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Tracer("test"), traceExporter
}

// AssertSpan is used to test whether a span contains the expected attributes.
func AssertSpan(t *testing.T, span *tracetest.SpanStub, name string, attributes map[string]any) {
	assert.Equal(t, name, span.Name)
	AssertAttributes(t, span.Attributes, attributes)
}

func AssertAttributes(t *testing.T, got []attribute.KeyValue, want map[string]any) {
	assert.Len(t, got, len(want))
	var gotKeys []string
	for _, attr := range got {
		gotKeys = append(gotKeys, string(attr.Key))
		w, ok := want[string(attr.Key)]
		if !ok {
			t.Errorf("unexpected attribute %s", attr.Key)
		} else {
			switch w := w.(type) {
			case string:
				assert.Equal(t, w, attr.Value.AsString())
			case int:
				assert.Equal(t, w, int(attr.Value.AsInt64()))
			case int64:
				assert.Equal(t, w, attr.Value.AsInt64())
			case bool:
				assert.Equal(t, w, attr.Value.AsBool())
			case float64:
				assert.Equal(t, w, attr.Value.AsFloat64())
			case func(value attribute.Value) bool:
				assert.True(t, w(attr.Value), "attribute %s does not match custom function", attr.Key)
			default:
				t.Errorf("unsupported attribute type %T", w)
			}
		}
	}
	sort.Strings(gotKeys)
	wantKeys := maps.Keys(want)
	sort.Strings(wantKeys)
	assert.Equal(t, wantKeys, gotKeys)
}
