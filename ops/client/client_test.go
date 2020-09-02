package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"bou.ke/monkey"
	"github.com/feiyuw/simgo/protocols"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

type mockClient struct {
}

func (mc *mockClient) InvokeRPC(mtd string, data interface{}) (interface{}, error) {
	return nil, nil
}

func (mc *mockClient) Close() error {
	return nil
}

func TestClientRESTAPIs(t *testing.T) {
	e := echo.New()

	Convey("show all clients in id order", t, func() {
		id1, _ := clientStorage.Add(&Client{
			Protocol: "grpc",
			Server:   "127.0.0.1:1234",
			Options:  map[string]interface{}{"protos": []string{"helloworld.proto"}},
		})
		id2, _ := clientStorage.Add(&Client{
			Protocol: "http",
			Server:   "127.0.0.1:1235",
		})
		id3, _ := clientStorage.Add(&Client{
			Protocol: "dubbo",
			Server:   "127.0.0.1:1237",
		})
		defer clientStorage.Remove(id1)
		defer clientStorage.Remove(id2)
		defer clientStorage.Remove(id3)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := Query(c)
		So(err, ShouldBeNil)
		So(rec.Code, ShouldEqual, http.StatusOK)
		clients := []*Client{}
		json.Unmarshal(rec.Body.Bytes(), &clients)
		So(len(clients), ShouldEqual, 3)
		So(clients[0].Id, ShouldEqual, id1)
		So(clients[1].Id, ShouldEqual, id2)
		So(clients[2].Id, ShouldEqual, id3)
	})

	Convey("concurrent clients will use different ID", t, func() {
		monkey.Patch(protocols.NewRpcClient, func(protocol string, server string, options map[string]interface{}) (protocols.RpcClient, error) {
			return &protocols.GrpcClient{}, nil
		})
		defer monkey.Unpatch(protocols.NewRpcClient)

		var wg sync.WaitGroup
		var cnt = 100
		nextClientID = 0 // reset clientId to 0

		defer func() {
			for idx := 0; idx < cnt; idx++ {
				clientStorage.Remove(uint64(idx + 1))
			}
		}()

		for idx := 0; idx < cnt; idx++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", strings.NewReader(`{"server":"127.0.0.1:4001","protocol":"grpc","options":{"protos":["upload/003886855.helloworld.proto"]}}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				New(c)
			}()
		}
		wg.Wait()
		clients, err := clientStorage.FindAll()
		So(err, ShouldBeNil)
		So(len(clients), ShouldEqual, cnt)

		idMap := make(map[uint64]int, cnt)
		for _, client := range clients {
			key := client.Id
			if _, exists := idMap[key]; exists {
				idMap[key]++
			} else {
				idMap[key] = 1
			}
		}
		So(len(idMap), ShouldEqual, cnt)
	})

	Convey("remove one client", t, func() {
		id4, _ := clientStorage.Add(&Client{
			Protocol:  "dubbo",
			Server:    "127.0.0.1:1237",
			RpcClient: &mockClient{},
		})
		defer clientStorage.Remove(id4)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/clients?id="+strconv.FormatUint(id4, 10), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		Delete(c)
		So(rec.Code, ShouldEqual, http.StatusOK)
		So(rec.Body.String(), ShouldEqual, "null\n")
	})
}
