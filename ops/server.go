package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"simgo/logger"
	"sort"
	//"simgo/protocols"
	"simgo/storage"
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
	Name     string   `json:"name"`
	Protocol string   `json:"protocol"`
	Clients  []string `json:"clients"` // TODO: clients identifier
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
	return c.JSON(http.StatusOK, nil)
}
