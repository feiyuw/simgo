package ops

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"simgo/logger"
)

var (
	opsServer = echo.New()
)

func Start(addr string) {
	opsServer.Use(middleware.Logger())
	opsServer.Use(middleware.Recover())
	opsServer.Pre(middleware.Rewrite(map[string]string{
		"/": "/static/index.html",
	}))

	// routes
	opsServer.Static("/static", "www")
	logger.Fatal("ops", opsServer.Start(addr))
}
