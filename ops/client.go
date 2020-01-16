package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Client struct {
	Id       int                    `json:"id"`
	Protocol string                 `json:"protocol"`
	Server   string                 `json:"server"`
	Options  map[string]interface{} `json:"options"`
}

func listClients(c echo.Context) error {
	mockClients := []*Client{
		&Client{Id: 0, Protocol: "grpc", Server: "127.0.0.1:1777", Options: map[string]interface{}{"protos": []string{"hello.proto", "echo.proto"}}},
		&Client{Id: 1, Protocol: "http", Server: "127.0.0.1:8080"},
		&Client{Id: 2, Protocol: "dubbo", Server: "127.0.0.1:3001"},
	}
	return c.JSON(http.StatusOK, mockClients)
}

func newClient(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
