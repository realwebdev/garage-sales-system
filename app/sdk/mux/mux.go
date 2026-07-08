package mux

import (
	"embed"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/realwebdev/garage-sales-system/app/sdk/auth"
	"github.com/realwebdev/garage-sales-system/app/sdk/authclient"
	"github.com/realwebdev/garage-sales-system/app/sdk/mid"
	"github.com/realwebdev/garage-sales-system/business/domain/auditbus"
	"github.com/realwebdev/garage-sales-system/business/domain/homebus"
	"github.com/realwebdev/garage-sales-system/business/domain/productbus"
	"github.com/realwebdev/garage-sales-system/business/domain/userbus"
	"github.com/realwebdev/garage-sales-system/business/domain/vproductbus"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
	"github.com/realwebdev/garage-sales-system/foundation/web"
	"go.opentelemetry.io/otel/trace"
)

// StaticSite represents a static site to run.
type StaticSite struct {
	react      bool
	static     embed.FS
	staticDir  string
	staticPath string
}

// Options represent optional parameters.
type Options struct {
	coresOrigin []string
	sites       []StaticSite
}

// WithCors provides configuration options for CORS.
func WithCORS(origin []string) func(opt *Options) {
	return func(opt *Options) {
		opt.coresOrigin = origin
	}
}

// WithFileServer provides configuration options for the file server.
func WithFileServer(react bool, static embed.FS, dir string, path string) func(opt *Options) {
	return func(opts *Options) {
		opts.sites = append(opts.sites, StaticSite{
			react:      react,
			static:     static,
			staticDir:  dir,
			staticPath: path,
		})
	}
}

// SalesConfig contains sales service specific config.
type SalesConfig struct {
	AuthClient authclient.Authenticator
}

// AuthConfig contains auth service specific config.
type AuthConfig struct {
	Auth *auth.Auth
}

type BusConfig struct {
	AuditBus    auditbus.ExtBusiness
	UserBus     userbus.ExtBusiness
	ProductBus  productbus.ExtBusiness
	HomeBus     homebus.ExtBusiness
	VProductBus vproductbus.ExtBusiness
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build       string
	Log         *logger.Logger
	DB          *sqlx.DB
	Tracer      trace.Tracer
	BusConfig   BusConfig
	SalesConfig SalesConfig
	AuthConfig  AuthConfig
}

// RouterAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs an http.Handler with all application routes.
func WebAPI(cfg Config, routeAdder RouteAdder, options ...func(opt *Options)) http.Handler {
	app := web.NewApp(
		cfg.Log.Info,
		cfg.Tracer,
		mid.Otel(cfg.Tracer),
		mid.Logger(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	if len(opts.coresOrigin) > 0 {
		app.EnableCORS(opts.coresOrigin)
	}

	routeAdder.Add(app, cfg)

	for _, site := range opts.sites {
		switch site.react {
		case true:
			app.FileServerReact(site.static, site.staticDir, site.staticPath)

		default:
			app.FileServer(site.static, site.staticDir, site.staticPath)

		}
	}

	return app
}
