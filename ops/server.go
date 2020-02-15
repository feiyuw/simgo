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

type Server struct {
	Name      string                 `json:"name"`
	Protocol  string                 `json:"protocol"`
	Port      int                    `json:"port"`
	Options   map[string]interface{} `json:"options"`
	Clients   []string               `json:"clients"` // TODO: clients identifier
	RpcServer protocols.RpcServer
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
	serverStorage.Add(server.Name, server)
	return c.JSON(http.StatusOK, nil)
}

func deleteServer(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func updateServer(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
