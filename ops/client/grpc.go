package client

import (
	"net/http"
	"strconv"

	"simgo/protocols"

	"github.com/labstack/echo/v4"
)

func ListGrpcServices(c echo.Context) error {
	clientId, err := strconv.ParseUint(c.QueryParam("clientId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "incorrect clientId")
	}
	client, err := clientStorage.FindOne(clientId)
	if err != nil {
		return err
	}
	if client.Protocol != "grpc" {
		return c.String(http.StatusBadRequest, "invalid protocol")
	}
	services, err := client.RpcClient.(*protocols.GrpcClient).ListServices()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, services)
}

func ListGrpcMethods(c echo.Context) error {
	clientId, err := strconv.ParseUint(c.QueryParam("clientId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "incorrect clientId")
	}
	service := c.QueryParam("service")
	client, err := clientStorage.FindOne(clientId)
	if err != nil {
		return err
	}
	if client.Protocol != "grpc" {
		return c.String(http.StatusBadRequest, "invalid protocol")
	}
	methods, err := client.RpcClient.(*protocols.GrpcClient).ListMethods(service)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, methods)
}
