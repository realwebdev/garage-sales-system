package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/ardanlabs/service/foundation/otel"
	"github.com/realwebdev/garage-sales-system/app/sdk/errs"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
	"github.com/realwebdev/garage-sales-system/foundation/web"
)

func Errors(log *logger.Logger) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			resp := next(ctx, r)

			err := checkIsError(resp)
			if err == nil {
				return resp
			}

			_, span := otel.AddSpan(ctx, "app.sdk.mid.error")
			span.RecordError(err)
			defer span.End()

			var appErr *errs.Error
			if !errors.As(err, &appErr) {
				appErr = errs.Errorf(errs.Internal, "Internal Server Error")
			}
		

		log.Error(ctx, "handled error during request")
	}
}
