package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"simgo/protocols"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/labstack/echo/v4"
	"github.com/robertkrimen/otto"
	"google.golang.org/grpc"
)

func Query(c echo.Context) error {
	servers, err := serverStorage.FindAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, servers)
}

func New(c echo.Context) error {
	server := new(Server)
	if err := c.Bind(server); err != nil {
		return err
	}
	rpcServer, err := protocols.NewRpcServer(server.Protocol, server.Name, server.Port, server.Options)
	if err != nil {
		return err
	}

	server.RpcServer = rpcServer
	server.MethodHandlers = map[string]*MethodHandler{}

	serverId, err := serverStorage.Add(server)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = server.RpcServer.Start(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	server.RpcServer.AddListener(newMessageRecorder(server))

	return c.JSON(http.StatusOK, map[string]uint64{"id": serverId})
}

func Delete(c echo.Context) error {
	serverId, err := strconv.ParseUint(c.QueryParam("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "incorrect server ID")
	}
	if err := serverStorage.Remove(serverId); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

type MethodHandler struct {
	ServerID uint64 `json:"serverId"`
	Method   string `json:"method"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

func ListMethodHandlers(c echo.Context) error {
	serverId, err := strconv.ParseUint(c.QueryParam("serverId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "incorrect server ID")
	}
	server, err := serverStorage.FindOne(serverId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, server.MethodHandlers)
}

func AddMethodHandler(c echo.Context) error {
	handler := new(MethodHandler)
	if err := c.Bind(handler); err != nil {
		return err
	}

	server, err := serverStorage.FindOne(handler.ServerID)
	if err != nil {
		return err
	}

	if _, exists := server.MethodHandlers[handler.Method]; exists {
		return errors.New("method handler exists")
	}
	switch server.Protocol {
	case "grpc":
		err = server.RpcServer.(*protocols.GrpcServer).SetMethodHandler(handler.Method, func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			switch handler.Type {
			case "raw":
				resp := make(map[string]interface{})
				if err := json.Unmarshal([]byte(handler.Content), &resp); err != nil {
					return err
				}
				for k, v := range resp {
					out.SetFieldByName(k, v)
				}
			case "javascript":
				vm := otto.New()
				vm.Set("ctx", map[string]interface{}{
					"in":     in,
					"out":    out,
					"stream": stream,
				})
				vm.Run(handler.Content)
			}

			return nil
		})
		if err != nil {
			return err
		}
		server.MethodHandlers[handler.Method] = handler
	}

	return c.JSON(http.StatusOK, nil)
}

func DeleteMethodHandler(c echo.Context) error {
	serverId, err := strconv.ParseUint(c.QueryParam("serverId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "incorrect server ID")
	}
	server, err := serverStorage.FindOne(serverId)
	if err != nil {
		return err
	}

	mtd := c.QueryParam("method")

	if _, exists := server.MethodHandlers[mtd]; exists {
		switch server.Protocol {
		case "grpc":
			if err := server.RpcServer.(*protocols.GrpcServer).RemoveMethodHandler(mtd); err != nil {
				return err
			}
			delete(server.MethodHandlers, mtd)
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func FetchMessages(c echo.Context) error {
	var (
		limit, skip int
	)

	serverId, err := strconv.ParseUint(c.QueryParam("serverId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "incorrect server ID")
	}
	server, err := serverStorage.FindOne(serverId)
	if err != nil {
		return err
	}

	if qLimit := c.QueryParam("limit"); qLimit != "" {
		limit, err = strconv.Atoi(qLimit)
	}
	if err != nil || limit <= 0 {
		limit = 30
	}
	if qSkip := c.QueryParam("skip"); qSkip != "" {
		skip, err = strconv.Atoi(qSkip)
	}
	if err != nil || skip < 0 {
		skip = 0
	}

	return c.JSON(http.StatusOK, queryMessages(server, skip, limit))
}
