package ops

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"simgo/logger"
	"simgo/ops/client"
	"simgo/ops/server"
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
	// clients
	opsServer.GET("/api/v1/clients", client.Query)
	opsServer.POST("/api/v1/clients", client.New)
	opsServer.DELETE("/api/v1/clients", client.Delete)
	opsServer.POST("/api/v1/clients/invoke", client.Invoke)
	opsServer.GET("/api/v1/clients/grpc/services", client.ListGrpcServices)
	opsServer.GET("/api/v1/clients/grpc/methods", client.ListGrpcMethods)
	// servers
	opsServer.GET("/api/v1/servers", server.Query)
	opsServer.POST("/api/v1/servers", server.New)
	opsServer.DELETE("/api/v1/servers", server.Delete)
	opsServer.GET("/api/v1/servers/messages", server.FetchMessages)
	opsServer.GET("/api/v1/servers/handlers", server.ListMethodHandlers)
	opsServer.POST("/api/v1/servers/handlers", server.AddMethodHandler)
	opsServer.DELETE("/api/v1/servers/handlers", server.DeleteMethodHandler)
	opsServer.GET("/api/v1/servers/grpc/methods", server.ListGrpcMethods)
	// other
	opsServer.POST("/api/v1/files", uploadFile)
	opsServer.DELETE("/api/v1/files", removeFile)
	opsServer.Static("/", "www/build")

	logger.Fatal("ops", opsServer.Start(addr))
}
