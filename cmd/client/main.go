package main

import (
	"context"
	"flag"
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
	"motus/goechoservice/internal/thrifttools"
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
		fmt.Fprint(os.Stderr, "Usage of the client", ":\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
	}

	config := thrifttools.Config{
		Protocol: *protocol,
		Framed:   *framed,
		Buffered: *buffered,
		Addr:     *addr,
	}

	err := thrifttools.WithClient(config, func(cc thrift.TClient) error {
		for i := 0; i < 10; i++ {
			client := echoservice.NewEchoServiceClient(cc)
			input := echoservice.NewTEchoServiceInputDTO()
			input.Message = fmt.Sprintf("Hello World #%d", i+1)

			output, err := client.Echo(context.Background(), input)
			if err != nil {
				return err
			}

			fmt.Printf("Server responded: '%v'\n", output.EchoMessage)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
