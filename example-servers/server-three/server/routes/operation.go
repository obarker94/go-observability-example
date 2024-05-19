package routes

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otracer "go.opentelemetry.io/otel/trace"
)

func Operation(w http.ResponseWriter, r *http.Request, t *trace.TracerProvider) {
	tracer := t.Tracer("modC")
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := tracer.Start(ctx, "Server3 /operation", otracer.WithSpanKind(otracer.SpanKindServer))

	defer span.End()

	// Simulate failure 30% of the time
	if randomFailure(0.3) {
		defer span.End()

		w.WriteHeader(http.StatusInternalServerError)

		span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...)
		span.SetStatus(codes.Error, "Internal Server Error")
		return
	}

	// Call to Server 2
	CallAnotherService(CallAnotherServiceParams{
		Tracer:   tracer,
		Ctx:      ctx,
		W:        w,
		Endpoint: "http://server2:8081/operation",
		SpanName: "Server 2",
	})

	w.Write([]byte("Operation complete"))

	span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...)
	span.SetStatus(codes.Ok, "Operation complete")

	// create a new internal span that says do something
	_, secondarySpan := tracer.Start(ctx, "Do something", otracer.WithSpanKind(otracer.SpanKindInternal))
	defer secondarySpan.End()

	w.Write([]byte("Operation from Server 3"))
}

func randomFailure(rate float64) bool {
	return rand.Float64() < rate
}

type CallAnotherServiceParams struct {
	Tracer   otracer.Tracer
	Ctx      context.Context
	W        http.ResponseWriter
	Endpoint string
	SpanName string
}

func CallAnotherService(params CallAnotherServiceParams) context.Context {
	client := &http.Client{}

	req, _ := http.NewRequestWithContext(params.Ctx, "GET", params.Endpoint, nil)
	ctx, span := params.Tracer.Start(params.Ctx, fmt.Sprintf("Call %s", params.SpanName), otracer.WithSpanKind(otracer.SpanKindClient))
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		msg := fmt.Sprintf("Failed to call %s", params.SpanName)
		span.SetStatus(codes.Error, msg)
		http.Error(params.W, msg, http.StatusInternalServerError)
		span.End()
		return params.Ctx
	}
	defer resp.Body.Close()

	span.SetStatus(codes.Ok, fmt.Sprintf("%s call success", params.SpanName))
	span.End()

	return ctx
}
