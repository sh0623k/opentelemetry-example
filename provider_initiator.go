package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type ProviderInitiator struct {
	connection *grpc.ClientConn
	resource   *resource.Resource
}

func NewProviderInitiator(ctx context.Context) (*ProviderInitiator, error) {
	connection, err := grpc.NewClient(
		"localhost:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	resource, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String("example-service")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return &ProviderInitiator{
		connection: connection,
		resource:   resource,
	}, nil
}

func (i *ProviderInitiator) initTracerProvider(ctx context.Context) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(i.connection))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(i.resource),
		// A batch span processor aggregates spans before export.
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExporter)),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func (i *ProviderInitiator) initMeterProvider(ctx context.Context) (func(context.Context) error, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(i.connection))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(i.resource),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}
