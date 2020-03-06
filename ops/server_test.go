package ops

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"simgo/protocols"
	"strings"
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

	Convey("unary grpc server e2e test", t, func() {
		// 1. new grpc server
		req := httptest.NewRequest(http.MethodPost, "/api/v1/servers", strings.NewReader(`{"name":"server_e2e","port":5000,"protocol":"grpc","options":{"protos":["../protocols/helloworld.proto"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := newServer(c)
		So(err, ShouldBeNil)
		So(rec.Code, ShouldEqual, http.StatusOK)
		servers, err := serverStorage.FindAll()
		So(err, ShouldBeNil)
		So(len(servers), ShouldEqual, 1)

		// 2. add handler
		req = httptest.NewRequest(http.MethodPost, "/api/v1/servers/handlers", strings.NewReader(`{"name":"server_e2e","method":"helloworld.Greeter.SayHello","type":"raw","content":"{\"message\":\"hello you\"}"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c = e.NewContext(req, rec)
		err = addMethodHandler(c)
		So(err, ShouldBeNil)
		So(rec.Code, ShouldEqual, http.StatusOK)

		// 3. send request and verify response
		client, err := protocols.NewGrpcClient("127.0.0.1:5000", []string{"../protocols/helloworld.proto"}, grpc.WithInsecure())
		So(err, ShouldBeNil)
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "you"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "hello you")

		// 4. fetch messages
		req = httptest.NewRequest(http.MethodGet, "/api/v1/servers/messages?name=server_e2e", nil)
		c = e.NewContext(req, rec)
		err = fetchMessages(c)
		So(err, ShouldBeNil)
		So(rec.Code, ShouldEqual, http.StatusOK)
		So(rec.Body.String(), ShouldNotEqual, "")
		server, err := serverStorage.FindOne("server_e2e")
		So(err, ShouldBeNil)
		So(len(server.(*Server).Messages), ShouldEqual, 2)

		// 5. delete server
		req = httptest.NewRequest(http.MethodDelete, "/api/v1/servers?name=server_e2e", nil)
		c = e.NewContext(req, rec)
		err = deleteServer(c)
		So(err, ShouldBeNil)
		So(rec.Code, ShouldEqual, http.StatusOK)
	})
}
