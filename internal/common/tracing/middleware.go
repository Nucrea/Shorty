package tracing

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const TraceIdHeader = "X-Trace-Id"
const TraceIdCtxKey = "traceId"

func NewMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	prop := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		traceId := ""
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
			c.Set(TraceIdCtxKey, traceId)
		}()

		ctx := prop.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))

		ctx, span := tracer.Start(ctx, "request")
		defer span.End()

		traceId = span.SpanContext().TraceID().String()
		c.Header(TraceIdHeader, traceId)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
