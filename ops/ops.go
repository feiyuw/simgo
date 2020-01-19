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
	opsServer.GET("/api/v1/clients", listClients)
	opsServer.POST("/api/v1/clients", newClient)
	opsServer.GET("/api/v1/grpc/services", listGrpcServices)
	opsServer.GET("/api/v1/grpc/methods", listGrpcMethods)
	opsServer.GET("/api/v1/servers", listServers)
	opsServer.POST("/api/v1/servers", newServer)
	opsServer.POST("/api/v1/files", uploadFile)
	opsServer.DELETE("/api/v1/files", removeFile)
	opsServer.Static("/", "www/build")

	logger.Fatal("ops", opsServer.Start(addr))
}
