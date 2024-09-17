# OpenTelemetry Example

## Starting docker containers

```bash
docker compose up -d
```

## Running the code

```bash
go run .
```

## Jaeger UI

The Jaeger UI is available at [http://localhost:16686](http://localhost:16686).

## Prometheus UI

The Prometheus UI is available at [http://localhost:9090](http://localhost:9090).

The number of runs is displayed as `test_runs_total`.

## Shut down

```bash
docker compose down
```

## References

- https://github.com/open-telemetry/opentelemetry-go/tree/main/example/otel-collector
- https://opentelemetry.io/docs/collector/configuration/#basics
- https://opentelemetry.io/docs/languages/go/getting-started/
- https://opentelemetry.io/docs/languages/go/exporters/
- https://opentelemetry.io/docs/languages/go/exporters/#otlp
- https://www.jaegertracing.io/docs/1.60/getting-started/
