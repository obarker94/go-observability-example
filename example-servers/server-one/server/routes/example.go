package routes

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	otracer "go.opentelemetry.io/otel/trace"
)

func Example(w http.ResponseWriter, r *http.Request, t *trace.TracerProvider) {
	tracer := t.Tracer("ExampleRoute")

	ctx, span := tracer.Start(r.Context(), "Server1 /tracing", otracer.WithSpanKind(otracer.SpanKindServer))
	defer span.End()

	// Call to Server 2
	CallAnotherService(CallAnotherServiceParams{
		Tracer:   tracer,
		Ctx:      ctx,
		W:        w,
		Endpoint: "http://server2:8081/operation",
		SpanName: "Server 2",
	})

	// Call to Server 3
	CallAnotherService(CallAnotherServiceParams{
		Tracer:   tracer,
		Ctx:      ctx,
		W:        w,
		Endpoint: "http://server3:8082/operation",
		SpanName: "Server 3",
	})
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
