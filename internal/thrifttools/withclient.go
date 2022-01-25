package thrifttools

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
)

func WithClient(config Config, handler func(thrift.TClient) error) error {
	protocolFactory, err := BuildProtocolFactory(config.Protocol)
	if err != nil {
		return err
	}

	cfg := &thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	transportFactory, err := BuildTransportFactory(config.Buffered, config.Framed, cfg)
	if err != nil {
		fmt.Fprint(os.Stderr, "Invalid transport specified", "\n")
		return err
	}

	var transport thrift.TTransport
	transport = thrift.NewTSocketConf(config.Addr, cfg)
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return err
	}

	defer transport.Close()
	if err := transport.Open(); err != nil {
		return err
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	cc := thrift.NewTStandardClient(iprot, oprot)

	return handler(cc)
}
