package ops

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"simgo/logger"
	"simgo/protocols"
	"simgo/storage"
	"sort"
	"strconv"
	"time"
)

const (
	MSGSIZE = 1000
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
	Method    string `json:"method"`
	Direction string `json:"direction"`
	From      string `json:"from"`
	To        string `json:"to"`
	Ts        int64  `json:"ts"`
	Body      string `json:"body"`
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

	server.Messages = make([]*Message, 0, MSGSIZE) // TODO: handle message save and load in storage
	server.RpcServer.AddListener(func(mtd, direction, from, to, body string) error {
		msg := &Message{
			Method:    mtd,
			Direction: direction,
			From:      from,
			To:        to,
			Ts:        time.Now().UnixNano() / time.Hour.Milliseconds(),
			Body:      body,
		}
		// TODO: add rlock
		if len(server.Messages) == MSGSIZE {
			copy(server.Messages[1:], server.Messages[0:MSGSIZE-1])
			server.Messages[0] = msg
		} else {
			server.Messages = append([]*Message{msg}, server.Messages...)
		}
		return nil
	})

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
	var (
		limit, skip int
		err         error
	)

	serverName := c.QueryParam("name")
	if qLimit := c.QueryParam("limit"); qLimit != "" {
		limit, err = strconv.Atoi(qLimit)
	}
	if err != nil {
		limit = 30
	}
	if qSkip := c.QueryParam("skip"); qSkip != "" {
		skip, err = strconv.Atoi(qSkip)
	}
	if err != nil {
		skip = 0
	}

	server, err := serverStorage.FindOne(serverName)
	if err != nil {
		return err
	}

	messages := make([]*Message, 0, limit)
	copy(messages, server.(*Server).Messages[skip:skip+limit])
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, messages)
}
