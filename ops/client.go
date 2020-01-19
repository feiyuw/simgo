package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"simgo/logger"
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
	Id       int                    `json:"id"`
	Protocol string                 `json:"protocol"`
	Server   string                 `json:"server"`
	Options  map[string]interface{} `json:"options"`
}

func listClients(c echo.Context) error {
	clientStorage.Add("0", &Client{Id: 0, Protocol: "grpc", Server: "127.0.0.1:1777", Options: map[string]interface{}{"protos": []string{"hello.proto", "echo.proto"}}})
	clientStorage.Add("1", &Client{Id: 1, Protocol: "http", Server: "127.0.0.1:8080"})
	clientStorage.Add("2", &Client{Id: 2, Protocol: "dubbo", Server: "127.0.0.1:3001"})

	allClients, err := clientStorage.FindAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, allClients)
}

func newClient(c echo.Context) error {
	client := new(Client)
	if err := c.Bind(client); err != nil {
		return err
	}
	clientStorage.Add(strconv.FormatUint(atomic.AddUint64(&nextClientID, 1), 10), client)
	return c.JSON(http.StatusOK, nil)
}
