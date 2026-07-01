package checkapp

import "github.com/realwebdev/garage-sales-system/foundation/web"

// Routes registers the routes for the checkapp.
func Routes(appWeb *web.App, build string) {
	a := newApp(build)

	appWeb.HandlerFunc("GET", "/health", a.health)
}
