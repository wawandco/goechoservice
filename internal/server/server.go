package server

import (
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
	"motus/goechoservice/internal/service"
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"
)

type protocolGetter interface {
	GetProtocol(thrift.TTransport) thrift.TProtocol
}

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/thrift", handlerForProtocol(thrift.NewTBinaryProtocolFactoryConf(nil)))
	mux.HandleFunc("/thrift/json", handlerForProtocol(thrift.NewTJSONProtocolFactory()))

	return mux
}

func handlerForProtocol(pb protocolGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "polo\n")
		case http.MethodPost:
			tt := thrift.NewStreamTransport(r.Body, w)
			protocol := pb.GetProtocol(tt)

			processor := echoservice.NewEchoServiceProcessor(service.Echo{})
			_, err := processor.Process(r.Context(), protocol, protocol)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
