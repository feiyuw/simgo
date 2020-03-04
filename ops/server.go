package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"simgo/logger"
	"simgo/protocols"
	"simgo/storage"
	"sort"
)

var (
	serverStorage storage.Storage
)

func init() {
	var err error

	if serverStorage, err = storage.NewMemoryStorage(); err != nil {
		logger.Fatal("ops/server", "init storage error", err)
	}
}

type Message struct {
	Method string `json:"method"`
	From   string `json:"from"`
	To     string `json:"to"`
	Ts     int    `json:"ts"`
	Body   string `json:"body"`
}

type Server struct {
	Name      string                 `json:"name"`
	Protocol  string                 `json:"protocol"`
	Port      int                    `json:"port"`
	Options   map[string]interface{} `json:"options"`
	Clients   []string               `json:"clients"` // TODO: clients identifier
	RpcServer protocols.RpcServer
	Messages  []*Message
}

func listServers(c echo.Context) error {
	servers, err := serverStorage.FindAll()
	if err != nil {
		return err
	}

	// sort servers with Name order
	sort.Slice(servers, func(idx1, idx2 int) bool {
		return servers[idx1].(*Server).Name < servers[idx2].(*Server).Name
	})

	return c.JSON(http.StatusOK, servers)
}

func newServer(c echo.Context) error {
	server := new(Server)
	if err := c.Bind(server); err != nil {
		return err
	}
	rpcServer, err := protocols.NewRpcServer(server.Protocol, server.Name, server.Port, server.Options)
	if err != nil {
		return err
	}

	server.RpcServer = rpcServer

	if err = serverStorage.Add(server.Name, server); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = server.RpcServer.Start(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, nil)
}

func deleteServer(c echo.Context) error {
	serverName := c.QueryParam("name")
	server, err := serverStorage.FindOne(serverName)
	if err != nil {
		return err
	}
	if err := server.(*Server).RpcServer.Close(); err != nil {
		return err
	}
	if err := serverStorage.Remove(serverName); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func fetchMessages(c echo.Context) error {
	serverName := c.QueryParam("name")
	//limit := c.QueryParam("limit")
	//skip := c.QueryParam("skip")
	server, err := serverStorage.FindOne(serverName)
	if err != nil {
		return err
	}
	messages := server.(*Server).Messages
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, messages)
}
