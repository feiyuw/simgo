package client

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/feiyuw/simgo/protocols"
	"github.com/feiyuw/simgo/utils"
)

func Query(c echo.Context) error {
	clients, err := clientStorage.FindAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, clients)
}

func New(c echo.Context) error {
	client := new(Client)
	if err := c.Bind(client); err != nil {
		return err
	}

	rpcClient, err := protocols.NewRpcClient(client.Protocol, client.Server, client.Options)
	if err != nil {
		return err
	}
	client.RpcClient = rpcClient
	clientStorage.Add(client)
	return c.JSON(http.StatusOK, nil)
}

func Delete(c echo.Context) error {
	clientId, err := utils.AtoUint64(c.QueryParam("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "incorrect clientId")
	}
	if err := clientStorage.Remove(clientId); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

type rpcRequest struct {
	ClientID uint64 `json:"clientId"`
	Method   string `json:"method"`
	Data     string `json:"data"`
}

func Invoke(c echo.Context) error {
	req := new(rpcRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	client, err := clientStorage.FindOne(req.ClientID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "client not found!")
	}

	resp, err := client.RpcClient.InvokeRPC(req.Method, req.Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
