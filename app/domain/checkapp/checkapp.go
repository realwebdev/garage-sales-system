package checkapp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/realwebdev/garage-sales-system/foundation/web"
)

type app struct {
	build string
}

func newApp(build string) *app {
	return &app{
		build: build,
	}
}

// Health represents the health status of the service.
type Health struct {
	Status string `json:"status"`
	Build  string `json:"build"`
	Host   string `json:"host"`
}

// Encode implements the web.Encoder interface.
func (h Health) Encode() ([]byte, string, error) {
	data, err := json.Marshal(h)
	if err != nil {
		return nil, "", err
	}
	return data, "application/json", nil
}

func (a *app) health(ctx context.Context, r *http.Request) web.Encoder {
	host, _ := os.Hostname()

	return Health{
		Status: "OK",
		Build:  a.build,
		Host:   host,
	}
}
