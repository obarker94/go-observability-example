package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/sdk/trace"
)

type Route func(http.ResponseWriter, *http.Request, *trace.TracerProvider)

func (m *Middleware) Public(next Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, m.TracerProvider)
	}
}
