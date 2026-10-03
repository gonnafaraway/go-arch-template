package trace

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Client interface {
	Start(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span)
	Shutdown(ctx context.Context) error
}

type otelClient struct {
	tracerProvider *tracesdk.TracerProvider
	tracer         oteltrace.Tracer
}

func NewClient(serviceName string) (Client, error) {
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

	return &otelClient{
		tracerProvider: tp,
		tracer:         tp.Tracer(serviceName),
	}, nil
}

func NewNoopClient() Client {
	return &noopClient{}
}

func newExporter() (tracesdk.SpanExporter, error) {
	otlpEndpoint := "localhost:4317"
	if endpoint := os.Getenv("OTLP_ENDPOINT"); endpoint != "" {
		otlpEndpoint = endpoint
	}

	return otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithInsecure(),
	)
}

func (c *otelClient) Start(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return c.tracer.Start(ctx, spanName, opts...)
}

func (c *otelClient) Shutdown(ctx context.Context) error {
	return c.tracerProvider.Shutdown(ctx)
}

type noopClient struct{}

func (c *noopClient) Start(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return ctx, oteltrace.SpanFromContext(ctx)
}

func (c *noopClient) Shutdown(ctx context.Context) error {
	return nil
}
