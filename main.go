package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	providerInitiator, err := NewProviderInitiator(ctx)
	if err != nil {
		log.Fatal(err)
	}
	shutdownTracerProvider, err := providerInitiator.initTracerProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Flush remaining spans and shut down the exporter.
		if err := shutdownTracerProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()
	shutdownMeterProvider, err := providerInitiator.initMeterProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownMeterProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown MeterProvider: %s", err)
		}
	}()

	name := "example"
	tracer := otel.Tracer(name)
	meter := otel.Meter(name)

	attributes := []attribute.KeyValue{
		attribute.String("attribute1", "hoge"),
		attribute.String("attribute2", "fuga"),
	}

	runCounter, err := meter.Int64Counter("runs")
	if err != nil {
		log.Fatal(err)
	}

	ctx, span := tracer.Start(
		ctx,
		"example-trace",
		trace.WithAttributes(attributes...),
	)
	defer span.End()
	for i := 0; i < 10; i++ {
		_, internalSpan := tracer.Start(ctx, fmt.Sprintf("example-%d", i))
		runCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		log.Printf("Do something. (%d / 10)\n", i+1)

		<-time.After(time.Duration(rand.Intn(5)+1) * 100 * time.Millisecond)
		internalSpan.End()
	}

	log.Printf("Completed.")
}
