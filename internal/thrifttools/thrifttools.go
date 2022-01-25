package thrifttools

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

func BuildProtocolFactory(protocol string) (thrift.TProtocolFactory, error) {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactoryConf(nil)
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactoryConf(nil)
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryConf(nil)
	default:
		return nil, fmt.Errorf("Invalid protocol specified: %s", protocol)
	}

	return protocolFactory, nil
}

func BuildTransportFactory(buffered, framed bool, cfg *thrift.TConfiguration) (thrift.TTransportFactory, error) {
	var transportFactory thrift.TTransportFactory
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactoryConf(transportFactory, cfg)
	}

	return transportFactory, nil
}
