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
	// clients
	opsServer.GET("/api/v1/clients", listClients)
	opsServer.POST("/api/v1/clients", newClient)
	opsServer.DELETE("/api/v1/clients", deleteClient)
	opsServer.POST("/api/v1/clients/invoke", invokeClientRPC)
	// servers
	opsServer.GET("/api/v1/servers", listServers)
	opsServer.POST("/api/v1/servers", newServer)
	opsServer.DELETE("/api/v1/servers", deleteServer)
	opsServer.GET("/api/v1/servers/messages", fetchMessages)
	opsServer.GET("/api/v1/servers/handlers", listMethodHandlers)
	opsServer.POST("/api/v1/servers/handlers", addMethodHandler)
	opsServer.DELETE("/api/v1/servers/handlers", deleteMethodHandler)
	opsServer.GET("/api/v1/servers/grpc/methods", listServerGrpcMethods)
	// other
	opsServer.POST("/api/v1/files", uploadFile)
	opsServer.DELETE("/api/v1/files", removeFile)
	opsServer.Static("/", "www/build")
	// grpc specified
	opsServer.GET("/api/v1/clients/grpc/services", listGrpcServices)
	opsServer.GET("/api/v1/clients/grpc/methods", listGrpcMethods)

	logger.Fatal("ops", opsServer.Start(addr))
}
