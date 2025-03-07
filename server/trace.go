package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TracingM(tracer trace.Tracer) gin.HandlerFunc {
	prop := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx := prop.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))

		ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path))
		defer span.End()

		traceId := span.SpanContext().TraceID()
		c.Header("X-Trace-Id", traceId.String())

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
