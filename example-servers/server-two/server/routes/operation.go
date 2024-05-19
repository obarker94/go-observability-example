package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otracer "go.opentelemetry.io/otel/trace"
)

func Operation(w http.ResponseWriter, r *http.Request, t *trace.TracerProvider) {
	tracer := t.Tracer("modB")
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := tracer.Start(ctx, "Server2 /operation", otracer.WithSpanKind(otracer.SpanKindServer))

	defer span.End()

	// Simulate failure 30% of the time
	if randomFailure(0.3) {
		defer span.End()

		w.WriteHeader(http.StatusInternalServerError)

		span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...)
		span.SetStatus(codes.Error, "Internal Server Error")
		return
	}

	fmt.Println("Operation is being performed")
	time.Sleep(200 * time.Millisecond)

	w.Write([]byte("Operation complete"))

	span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...)
	span.SetStatus(codes.Ok, "Operation complete")

	// create a new internal span that says do something
	_, secondarySpan := tracer.Start(ctx, "Do something", otracer.WithSpanKind(otracer.SpanKindInternal))
	defer secondarySpan.End()

	w.Write([]byte("Operation from Server 2"))
}

func randomFailure(rate float64) bool {
	return rand.Float64() < rate
}
