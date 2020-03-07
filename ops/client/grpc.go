package client

import (
	"errors"
	"net/http"

	"simgo/protocols"

	"github.com/labstack/echo/v4"
)

func ListGrpcServices(c echo.Context) error {
	clientId := c.QueryParam("clientId")
	client, err := clientStorage.FindOne(clientId)
	if err != nil {
		return err
	}
	if client.(*Client).Protocol != "grpc" {
		return errors.New("invalid protocol")
	}
	services, err := client.(*Client).RpcClient.(*protocols.GrpcClient).ListServices()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, services)
}

func ListGrpcMethods(c echo.Context) error {
	clientId := c.QueryParam("clientId")
	service := c.QueryParam("service")
	client, err := clientStorage.FindOne(clientId)
	if err != nil {
		return err
	}
	if client.(*Client).Protocol != "grpc" {
		return errors.New("invalid protocol")
	}
	methods, err := client.(*Client).RpcClient.(*protocols.GrpcClient).ListMethods(service)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, methods)
}
