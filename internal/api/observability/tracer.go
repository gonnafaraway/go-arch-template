package observability

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	Shutdown(ctx context.Context) error
}

type otelTracer struct {
	tracerProvider *tracesdk.TracerProvider
	tracer         trace.Tracer
}

func NewTracer(serviceName string) (Tracer, error) {
	exporter, err := newExporter()
	if err != nil {
		return nil, err
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &otelTracer{
		tracerProvider: tp,
		tracer:         tp.Tracer(serviceName),
	}, nil
}

func newExporter() (tracesdk.SpanExporter, error) {
	// Пробуем OTLP, если не работает - используем Jaeger
	otlpEndpoint := "localhost:4317"
	if endpoint := os.Getenv("OTLP_ENDPOINT"); endpoint != "" {
		otlpEndpoint = endpoint
	}

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err == nil {
		return exporter, nil
	}

	// Fallback на Jaeger
	jaegerEndpoint := "http://localhost:14268/api/traces"
	if endpoint := os.Getenv("JAEGER_ENDPOINT"); endpoint != "" {
		jaegerEndpoint = endpoint
	}

	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
}

func (t *otelTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

func (t *otelTracer) Shutdown(ctx context.Context) error {
	return t.tracerProvider.Shutdown(ctx)
}

// Noop tracer для случаев когда трассировка недоступна
type noopTracer struct{}

func NewNoopTracer() Tracer {
	return &noopTracer{}
}

func (t *noopTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, trace.SpanFromContext(ctx)
}

func (t *noopTracer) Shutdown(ctx context.Context) error {
	return nil
}

