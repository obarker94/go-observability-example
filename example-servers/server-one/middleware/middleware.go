package middleware

import (
	"go.opentelemetry.io/otel/sdk/trace"
)

type Middleware struct {
	TracerProvider *trace.TracerProvider
}

func New(serviceName string) *Middleware {
	tracerProvider := InitTracer(serviceName, "otel-collector:4317")
	return &Middleware{TracerProvider: tracerProvider}
}
