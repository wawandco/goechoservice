package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
	"motus/goechoservice/internal/thrifttools"
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
		fmt.Fprint(os.Stderr, "Usage of the echoservice", ":\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
	}

	flag.Parse()

	protocolFactory, err := thrifttools.BuildProtocolFactory(*protocol)
	if err != nil {
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		flag.Usage()
		os.Exit(1)
	}

	cfg := &thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	transportFactory, err := thrifttools.BuildTransportFactory(*buffered, *framed, cfg)
	if err != nil {
		fmt.Fprint(os.Stderr, "Invalid transport specified", "\n")
		os.Exit(1)
	}

	transport, err := thrift.NewTServerSocket(*addr)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error building transport: "+err.Error(), "\n")
		os.Exit(1)
	}

	processor := echoservice.NewEchoServiceProcessor(service.Echo{})
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", *addr)
	if err := server.Serve(); err != nil {
		fmt.Println("error running server:", err)
	}
}
