package ops

import (
	"bou.ke/monkey"
	"encoding/json"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"simgo/protocols"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestClientRESTAPIs(t *testing.T) {
	e := echo.New()

	Convey("show all clients in id order", t, func() {
		clientStorage.Add("2", &Client{
			Id:       "2",
			Protocol: "grpc",
			Server:   "127.0.0.1:1234",
			Options:  map[string]interface{}{"protos": []string{"helloworld.proto"}},
		})
		clientStorage.Add("1", &Client{
			Id:       "1",
			Protocol: "http",
			Server:   "127.0.0.1:1235",
		})
		clientStorage.Add("3", &Client{
			Id:       "3",
			Protocol: "dubbo",
			Server:   "127.0.0.1:1237",
		})
		defer clientStorage.Remove("1")
		defer clientStorage.Remove("2")
		defer clientStorage.Remove("3")
		req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", strings.NewReader(""))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		listClients(c)
		So(rec.Code, ShouldEqual, http.StatusOK)
		clients := []interface{}{}
		json.Unmarshal(rec.Body.Bytes(), &clients)
		So(len(clients), ShouldEqual, 3)
		So(clients[0].(map[string]interface{})["id"], ShouldEqual, "1")
		So(clients[1].(map[string]interface{})["id"], ShouldEqual, "2")
		So(clients[2].(map[string]interface{})["id"], ShouldEqual, "3")
	})

	Convey("concurrent clients will use different ID", t, func() {
		monkey.Patch(protocols.NewRpcClient, func(protocol string, server string, options map[string]interface{}) (protocols.RpcClient, error) {
			return &protocols.GrpcClient{}, nil
		})

		var wg sync.WaitGroup
		var cnt = 100

		defer func() {
			for idx := 0; idx < cnt; idx++ {
				clientStorage.Remove(strconv.Itoa(idx + 1))
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
				newClient(c)
			}()
		}
		wg.Wait()
		clients, err := clientStorage.FindAll()
		So(err, ShouldBeNil)
		So(len(clients), ShouldEqual, cnt)

		idMap := make(map[string]int, cnt)
		for _, client := range clients {
			key := client.(*Client).Id
			if _, exists := idMap[key]; exists {
				idMap[key]++
			} else {
				idMap[key] = 1
			}
		}
		So(len(idMap), ShouldEqual, cnt)
	})
}
