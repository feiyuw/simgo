package server

import (
	"errors"
	"net/http"
	"strconv"

	"simgo/protocols"

	"github.com/labstack/echo/v4"
)

func ListGrpcMethods(c echo.Context) error {
	serverId, err := strconv.ParseUint(c.QueryParam("serverId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "incorrect server ID")
	}
	server, err := serverStorage.FindOne(serverId)
	if err != nil {
		return err
	}

	if server.Protocol != "grpc" {
		return errors.New("incorrect protocol")
	}

	methods, err := server.RpcServer.(*protocols.GrpcServer).ListMethods()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, methods)
}
