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
		"/": "/index.html",
	}))

	// routes
	opsServer.GET("/api/clients", listClients)
	opsServer.POST("/api/clients", newClient)
	opsServer.GET("/api/servers", listServers)
	opsServer.POST("/api/servers", newServer)
	opsServer.Static("/", "www/build")

	logger.Fatal("ops", opsServer.Start(addr))
}
