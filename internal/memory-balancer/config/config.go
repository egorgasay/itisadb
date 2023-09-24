package config

import (
	"flag"
	"os"
)

type Config struct {
	GRPC string
	REST string
	URI  string
	Key  []byte
}

const (
	defaultGRPC = ":8888"
)

type Flag struct {
	grpc *string
	rest *string
	dsn  *string
}

var f Flag

func init() {
	f.grpc = flag.String("grpc", defaultGRPC, "-grpc=host")
	f.rest = flag.String("rest", "", "-rest=host")
	f.dsn = flag.String("d", "", "-d=dsn")
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		f.grpc = &addr
	}

	return &Config{
		GRPC: *f.grpc,
		REST: *f.rest,
		URI:  *f.dsn,
		Key:  []byte("CHANGE ME"),
	}
}
