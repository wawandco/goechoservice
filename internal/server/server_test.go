package server_test

import (
	"context"
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
	"motus/goechoservice/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
)

func TestServerRoutes(t *testing.T) {
	s := server.New()

	tcases := []struct {
		name    string
		method  string
		path    string
		code    int
		content string
	}{
		{
			name:    "GET /thrift",
			method:  http.MethodGet,
			path:    "/thrift",
			code:    http.StatusOK,
			content: "polo\n",
		},

		{
			name:    "GET /thrift/json",
			method:  http.MethodGet,
			path:    "/thrift",
			code:    http.StatusOK,
			content: "polo\n",
		},

		{
			name:    "GET /",
			method:  http.MethodGet,
			path:    "/",
			code:    http.StatusNotFound,
			content: "",
		},

		{
			name:    "PUT /thrift",
			method:  http.MethodPut,
			path:    "/thrift",
			code:    http.StatusMethodNotAllowed,
			content: "",
		},

		{
			name:    "PUT /thrift/json",
			method:  http.MethodPut,
			path:    "/thrift",
			code:    http.StatusMethodNotAllowed,
			content: "",
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			r, err := http.NewRequest(tcase.method, tcase.path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			w := httptest.NewRecorder()
			s.ServeHTTP(w, r)

			if w.Code != tcase.code {
				t.Fatalf("Expected status code %d, got %d", tcase.code, w.Code)
			}

			if tcase.content == "" {
				return
			}

			if !strings.Contains(w.Body.String(), tcase.content) {
				t.Fatalf("Expected body to contain %q, got %q", tcase.content, w.Body.String())
			}
		})
	}
}

func TestServerThriftCalls(t *testing.T) {
	ts := httptest.NewServer(server.New())
	t.Cleanup(ts.Close)

	tcases := []struct {
		name            string
		path            string
		protocolFactory thrift.TProtocolFactory
	}{
		{
			name:            "binary call",
			path:            "/thrift",
			protocolFactory: thrift.NewTBinaryProtocolFactoryDefault(),
		},

		{
			name:            "json call",
			path:            "/thrift/json",
			protocolFactory: thrift.NewTJSONProtocolFactory(),
		},
	}

	for _, tcase := range tcases {
		tc, err := thrift.NewTHttpClient(ts.URL + tcase.path)
		if err != nil {
			t.Fatalf("Could not create client: %v", err)
		}

		client := echoservice.NewEchoServiceClientFactory(tc, tcase.protocolFactory)
		if err := tc.Open(); err != nil {
			t.Fatalf("Could not open client: %v", err)
		}

		// Defering close of the transport
		defer func() {
			if err := tc.Close(); err != nil {
				panic(err)
			}
		}()

		out, err := client.Echo(context.Background(), &echoservice.TEchoServiceInputDTO{
			Message: fmt.Sprintf("Hello from %v", tcase.name),
		})

		if err != nil {
			t.Fatalf("Could not call Echo: %v", err)
		}

		if out.EchoMessage != fmt.Sprintf("ACK from server: 'Hello from %v'", tcase.name) {
			t.Fatalf("Expected ACK from server, got '%s'", out.EchoMessage)
		}
	}
}
