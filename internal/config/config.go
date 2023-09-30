package config

import (
	"flag"
)

type Config struct {
	GRPC      string
	REST      string
	WebApp    string
	IsTLogger bool `toml:"TransactionLoggerDir"`
	Key       []byte
}

const (
	defaultGRPC       = ":8888"
	defaultTLogger    = "transaction-logger"
	defaultWebAppHost = ":6070"
)

type Flag struct {
	grpc       *string
	rest       *string
	tlog       *string
	webAppHost *string
}

var f Flag

func init() {
	f.grpc = flag.String("grpc", defaultGRPC, "-grpc=host")
	f.rest = flag.String("rest", "", "-rest=host")
	f.webAppHost = flag.String("a", defaultWebAppHost, "-a=host")
	f.tlog = flag.String("tlog", defaultTLogger, "-tlog=dir")
}

func New() *Config {
	flag.Parse()

	return &Config{
		GRPC:      *f.grpc,
		REST:      *f.rest,
		IsTLogger: *f.tlog != "",
		WebApp:    *f.webAppHost,
		Key:       []byte("CHANGE ME"),
	}
}
