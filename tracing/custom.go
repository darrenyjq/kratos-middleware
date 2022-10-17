package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var tp trace.TracerProvider

func InitTracerProvider(url string, samplingRate float64, serviceName, env string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp = tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(samplingRate))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String(keyEnv, env),
		)),
	)
	return nil
}

func AddTraceSpan(ctx context.Context, spanName string, attribs []attribute.KeyValue, err error) {
	if ctx == nil || tp == nil {
		return
	}
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer(defaultTracerName)
	kind := trace.SpanKindServer
	ctx, span := tracer.Start(ctx,
		spanName,
		trace.WithAttributes(attribs...),
		trace.WithSpanKind(kind),
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
