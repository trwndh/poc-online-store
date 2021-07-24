package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operationName := "HTTP " + r.Method + " " + r.URL.Path
		tracer := opentracing.GlobalTracer()
		serverSpanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span, traceCtx := opentracing.StartSpanFromContextWithTracer(r.Context(), tracer, operationName, ext.RPCServerOption(serverSpanCtx), opentracing.FollowsFrom(serverSpanCtx))
		defer span.Finish()

		span.SetTag("Header", r.Header)
		span.SetTag("Uber-Trace-Id", r.Header.Get("Uber-Trace-Id"))
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.Path)

		next.ServeHTTP(w, r.WithContext(traceCtx))
	})
}
