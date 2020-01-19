package ops

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"simgo/logger"
	"simgo/protocols"
	"simgo/storage"
	"sort"
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
	Id        string                 `json:"id"`
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

	// sort clients with ID order
	sort.Slice(clients, func(idx1, idx2 int) bool {
		return clients[idx1].(*Client).Id < clients[idx2].(*Client).Id
	})

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
	client.Id = strconv.FormatUint(atomic.AddUint64(&nextClientID, 1), 10)
	clientStorage.Add(client.Id, client)
	return c.JSON(http.StatusOK, nil)
}

func deleteClient(c echo.Context) error {
	clientId := c.QueryParam("id")
	client, err := clientStorage.FindOne(clientId)
	if err != nil {
		return err
	}
	if err := client.(*Client).RpcClient.Close(); err != nil {
		return err
	}
	if err := clientStorage.Remove(clientId); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

type rpcRequest struct {
	ClientID string `json:"clientId"`
	Method   string `json:"method"`
	Data     string `json:"data"`
}

func invokeClientRPC(c echo.Context) error {
	req := new(rpcRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	client, err := clientStorage.FindOne(req.ClientID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "client not found!")
	}

	resp, err := client.(*Client).RpcClient.InvokeRPC(req.Method, req.Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
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
