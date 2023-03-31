package config

import (
	"flag"
	"os"
)

type Config struct {
	Host string
	URI  string
	Key  []byte
}

const (
	defaultHost = "127.0.0.1:8080"
)

type Flag struct {
	host *string
	dsn  *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	f.dsn = flag.String("d", "", "-d=dsn")
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		f.host = &addr
	}

	return &Config{
		Host: *f.host,
		URI:  *f.dsn,
		Key:  []byte("CHANGE ME"),
	}
}
