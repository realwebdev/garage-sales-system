package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/app/sdk/auth"
	"github.com/realwebdev/garage-sales-system/app/sdk/authclient"
	"github.com/realwebdev/garage-sales-system/app/sdk/errs"
	"github.com/realwebdev/garage-sales-system/foundation/web"
)

// Authenticate is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authenticate(client authclient.Authenticator) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			resp, err := client.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctx = setUserID(ctx, resp.UserID)
			ctx = setClaims(ctx, resp.Claims)

			return next(ctx, r)
		}

		return h
	}

	return m
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			authorizationHeader := r.Header.Get("authorization")
			ctx, err := HandleAuthentication(ctx, ath, authorizationHeader)
			if err != nil {
				return err
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}

func HandleAuthentication(ctx context.Context, ath *auth.Auth, authorizationHeader string) (context.Context, *errs.Error) {
	claims, err := ath.Authenticate(ctx, authorizationHeader)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return ctx, errs.Errorf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Errorf(errs.Unauthenticated, "parsing subject: %s", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}
