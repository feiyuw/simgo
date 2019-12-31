package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Client struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Server   string `json:"server"`
}

func listClients(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func newClient(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
