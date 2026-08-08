// Package checkapp maintains the app layer api for the check domain.
package checkapp

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/jmoiron/sqlx"
	"github.com/realwebdev/garage-sales-system/app/sdk/errs"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
	"github.com/realwebdev/garage-sales-system/foundation/web"
)

type app struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func newApp(build string, log *logger.Logger, db *sqlx.DB) *app {
	return &app{
		build: build,
		log:   log,
		db:    db,
	}
}

// readiness checks if the databse is ready and if not will return a 500 status.
// Don not respond by just returning an error because further up in the call
// stack it will interpret that as non-trusted error.
func (a *app) readiness(ctx context.Context, r *http.Request) web.Encoder {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := sqldb.StatusCheck(ctx, a.db); err != nil {
		a.log.Info(ctx, "readiness failure", "ERROR", err)
		return errs.New(errs.Internal, err)
	}

	return nil
}

// liveness returns simple status info of the service is alive. If the
// app is deployed to a kubernetes cluster, it will also return pod, node, and
// namespace details via Downward API. The kubernetes environment variables
// need to be set within you Pod/Deployment manifest.
func (a *app) liveness(ctx context.Context, r *http.Request) web.Encoder {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailbale"
	}

	info := Info{
		Status:     "up",
		Build:      a.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	return info
}
