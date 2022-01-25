package thrifttools

type Config struct {
	Protocol string
	Framed   bool
	Buffered bool
	Addr     string
}
