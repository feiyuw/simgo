package server

import (
	"errors"
	"net/http"

	"github.com/feiyuw/simgo/protocols"
	"github.com/feiyuw/simgo/utils"

	"github.com/labstack/echo/v4"
)

func ListGrpcMethods(c echo.Context) error {
	serverId, err := utils.AtoUint64(c.QueryParam("serverId"))
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
