# Performance

Where this project's performance is measured, and how to read it. Anyone chasing a slowdown starts here, and so do the performance skills when installed.

## Telemetry store

<The system that holds production or staging measurements: Grafana / Prometheus, Datadog, Honeycomb, CloudWatch, a profiling service. Its URL, and how to get access.>

## What to measure

| Metric | Where | Query or dashboard | Budget |
| --- | --- | --- | --- |
| <p95 latency of the API> | <Prometheus> | <`histogram_quantile(0.95, sum by (le, route) (rate(http_request_duration_seconds_bucket[5m])))`> | <200 ms> |
| <memory per worker> | <Grafana dashboard> | <link> | <512 MiB> |

## Profiles

<Where continuous profiles live (Pyroscope, Datadog profiler, pprof endpoints), or "none".>

## Benchmarks

<The benchmark command (`go test -bench=. ./...`, `npm run bench`), where results are kept, and the machine they're comparable on.>
