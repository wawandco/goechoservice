//go:generate thrift --gen go echoservice.thrift
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
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

	var transport thrift.TTransport
	transport = thrift.NewTSocketConf(*addr, cfg)
	transport, err := transportFactory.GetTransport(transport)
	if err != nil {
		panic(err)
	}

	defer transport.Close()
	if err := transport.Open(); err != nil {
		panic(err)
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)

	client := echoservice.NewEchoServiceClient(thrift.NewTStandardClient(iprot, oprot))

	for i := 0; i < 10; i++ {
		input := echoservice.NewTEchoServiceInputDTO()
		input.Message = fmt.Sprintf("Hello World #%d", i+1)

		output, err := client.Echo(context.Background(), input)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Server responded: '%v'\n", output.EchoMessage)
	}
}
