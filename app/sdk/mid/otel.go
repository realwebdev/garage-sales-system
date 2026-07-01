package mid

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/otel"
	"github.com/realwebdev/garage-sales-system/foundation/web"
	"go.opentelemetry.io/otel/trace"
)

func Otel(tracer trace.Tracer) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx = otel.InjectTracing(ctx, tracer)

			return next(ctx, r)
		}

		return h
	}

	return m
}
