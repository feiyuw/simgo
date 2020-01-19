package ops

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"simgo/logger"
	"simgo/protocols"
	"simgo/storage"
	"strconv"
	"sync/atomic"
)

var (
	clientStorage storage.Storage
	nextClientID  uint64
)

func init() {
	var err error

	if clientStorage, err = storage.NewMemoryStorage(); err != nil {
		logger.Fatal("ops/client", "init storage error", err)
	}
}

type Client struct {
	Id        uint64                 `json:"id"`
	Protocol  string                 `json:"protocol"`
	Server    string                 `json:"server"`
	Options   map[string]interface{} `json:"options"`
	RpcClient protocols.RpcClient
}

func listClients(c echo.Context) error {
	clients, err := clientStorage.FindAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, clients)
}

func newClient(c echo.Context) error {
	client := new(Client)
	if err := c.Bind(client); err != nil {
		return err
	}

	rpcClient, err := protocols.NewRpcClient(client.Protocol, client.Server, client.Options)
	if err != nil {
		return err
	}
	client.RpcClient = rpcClient
	client.Id = atomic.AddUint64(&nextClientID, 1)
	clientStorage.Add(strconv.FormatUint(client.Id, 10), client)
	return c.JSON(http.StatusOK, nil)
}

func listGrpcServices(c echo.Context) error {
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

func listGrpcMethods(c echo.Context) error {
	return nil
}
