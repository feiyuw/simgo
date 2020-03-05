package ops

import (
	//"bou.ke/monkey"
	"encoding/json"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"simgo/protocols"
	//"strconv"
	"strings"
	//"sync"
	"testing"
)

func TestServerRESTAPIs(t *testing.T) {
	e := echo.New()

	Convey("show all servers in name order", t, func() {
		serverStorage.Add("hello2", &Server{
			Name:     "hello2",
			Protocol: "grpc",
			Port:     1235,
			Options:  map[string]interface{}{"protos": []string{"helloworld.proto"}},
		})
		serverStorage.Add("world1", &Server{
			Name:     "world1",
			Protocol: "grpc",
			Port:     1236,
			Options:  map[string]interface{}{"protos": []string{"helloworld.proto"}},
		})
		serverStorage.Add("hello1", &Server{
			Name:     "hello1",
			Protocol: "grpc",
			Port:     1234,
			Options:  map[string]interface{}{"protos": []string{"helloworld.proto"}},
		})
		defer serverStorage.Remove("hello1")
		defer serverStorage.Remove("hello2")
		defer serverStorage.Remove("world1")
		req := httptest.NewRequest(http.MethodGet, "/api/v1/servers", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		listServers(c)
		So(rec.Code, ShouldEqual, http.StatusOK)
		servers := []interface{}{}
		json.Unmarshal(rec.Body.Bytes(), &servers)
		So(len(servers), ShouldEqual, 3)
		So(servers[0].(map[string]interface{})["name"], ShouldEqual, "hello1")
		So(servers[1].(map[string]interface{})["name"], ShouldEqual, "hello2")
		So(servers[2].(map[string]interface{})["name"], ShouldEqual, "world1")
	})

	Convey("grpc server e2e test", t, func() {
		// 1. new grpc server
		req := httptest.NewRequest(http.MethodPost, "/api/v1/servers", strings.NewReader(`{"name":"helloworld","port":5000,"protocol":"grpc","options":{"protos":["protocols/helloworld.proto"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		newServer(c)
		So(rec.Code, ShouldEqual, http.StatusOK)

		// 2. add handler

		// 3. send request and verify response
		client, err := protocols.NewGrpcClient("127.0.0.1:5000", []string{"../protocols/helloworld.proto"}, grpc.WithInsecure())
		So(err, ShouldBeNil)
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "you"})
		So(err, ShouldBeNil)
		//So(out.(map[string]interface{})["message"], ShouldEqual, "hello you")
		So(out.(map[string]interface{})["message"], ShouldBeNil)

		// 4. delete server
		req = httptest.NewRequest(http.MethodDelete, "/api/v1/servers?name=helloworld", nil)
		c = e.NewContext(req, rec)
		deleteServer(c)
		So(rec.Code, ShouldEqual, http.StatusOK)
	})
}
