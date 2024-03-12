// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
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

	return tracerProvider.Tracer("test"), traceExporter
}

// AssertSpan is used to test whether a span contains the expected attributes.
func AssertSpan(t *testing.T, span *tracetest.SpanStub, name string, attributes map[string]any) {
	assert.Equal(t, name, span.Name)
	assert.Len(t, span.Attributes, len(attributes))
	var gotKeys []string
	for _, attr := range span.Attributes {
		gotKeys = append(gotKeys, string(attr.Key))
		want, ok := attributes[string(attr.Key)]
		if !ok {
			t.Errorf("unexpected attribute %s", attr.Key)
		}
		switch want.(type) {
		case string:
			assert.Equal(t, want, attr.Value.AsString())
		case int:
			assert.Equal(t, want, int(attr.Value.AsInt64()))
		case int64:
			assert.Equal(t, want, attr.Value.AsInt64())
		case bool:
			assert.Equal(t, want, attr.Value.AsBool())
		case float64:
			assert.Equal(t, want, attr.Value.AsFloat64())
		default:
			t.Errorf("unsupported attribute type %T", want)
		}
	}
	sort.Strings(gotKeys)
	wantKeys := maps.Keys(attributes)
	sort.Strings(wantKeys)
	assert.Equal(t, wantKeys, gotKeys)
}
