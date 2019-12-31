package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	Name     string   `json:"name"`
	Protocol string   `json:"protocol"`
	Clients  []string `json:"clients"` // TODO: clients identifier
}

func listServers(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func newServer(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
