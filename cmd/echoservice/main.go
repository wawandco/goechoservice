//go:generate thrift --gen go echoservice.thrift
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
	"motus/goechoservice/service"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
)

var (
	protocol = flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed   = flag.Bool("framed", false, "Use framed transport")
	buffered = flag.Bool("buffered", false, "Use buffered transport")
	addr     = flag.String("addr", "localhost:9090", "Address to listen to")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
	}

	flag.Parse()

	var protocolFactory thrift.TProtocolFactory
	switch *protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactoryConf(nil)
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactoryConf(nil)
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryConf(nil)
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		flag.Usage()
		os.Exit(1)
	}

	var transportFactory thrift.TTransportFactory
	cfg := &thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if *framed {
		transportFactory = thrift.NewTFramedTransportFactoryConf(transportFactory, cfg)
	}

	if err := runServer(transportFactory, protocolFactory, *addr); err != nil {
		fmt.Println("error running server:", err)
	}
}

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string) error {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return err
	}

	processor := echoservice.NewEchoServiceProcessor(service.Echo{})
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	return server.Serve()
}
