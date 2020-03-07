package client

import (
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

func Query(c echo.Context) error {
	clients, err := clientStorage.FindAll()
	if err != nil {
		return err
	}

	// sort clients with ID order
	sort.Slice(clients, func(idx1, idx2 int) bool {
		intIdx1, err1 := strconv.Atoi(clients[idx1].(*Client).Id)
		intIdx2, err2 := strconv.Atoi(clients[idx2].(*Client).Id)
		if err1 != nil || err2 != nil {
			return clients[idx1].(*Client).Id < clients[idx2].(*Client).Id
		}
		return intIdx1 < intIdx2
	})

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
	client.Id = strconv.FormatUint(atomic.AddUint64(&nextClientID, 1), 10)
	clientStorage.Add(client.Id, client)
	return c.JSON(http.StatusOK, nil)
}

func Delete(c echo.Context) error {
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

func Invoke(c echo.Context) error {
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
