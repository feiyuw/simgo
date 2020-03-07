package server

import (
	"errors"
	"net/http"

	"simgo/protocols"

	"github.com/labstack/echo/v4"
)

func ListGrpcMethods(c echo.Context) error {
	serverName := c.QueryParam("name")
	server, err := serverStorage.FindOne(serverName)
	if err != nil {
		return err
	}

	if server.(*Server).Protocol != "grpc" {
		return errors.New("incorrect protocol")
	}

	methods, err := server.(*Server).RpcServer.(*protocols.GrpcServer).ListMethods()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, methods)
}
