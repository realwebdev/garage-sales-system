package checkapp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
	"github.com/realwebdev/garage-sales-system/foundation/web"
)

type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes registers the routes for the checkapp.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.Build, cfg.Log, cfg.DB)

	app.HandlerFunc(http.MethodGet, version, "/readiness", api.readiness)
	app.HandlerFunc(http.MethodGet, version, "/readiness", api.liveness)
}
