# Go Observability with OTEL

Many examples of observability setups often miss showing how OTEL can be used
directly in the code and show a functioning obersvability stack with just a simulator.

Whilst simulators are better suited to ensure your stack is working correctly,
it can be difficult to translate all the functionality observed into real code.

Thus, I've also included three minimal Go servers which setup traces, propgate
across HTTP calls and set client / server kinds in order to enable to the
Tempo service graph.

# The Observability Stack

- Open Telemetry - span / metric collector that apps write to.
- Prometheus - time series database that scrapes the OTEL server.
- Tempo - distributed tracing service.
- PromTail - log collector / agent
- Loki - log aggregation
- Grafana - visualiser for traces, logs and anything else you wish to add.
