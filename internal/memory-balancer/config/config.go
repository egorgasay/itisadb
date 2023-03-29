package config

import (
	"flag"
	"os"
)

type Config struct {
	Host string
	Key  []byte
}

const (
	defaultHost = "127.0.0.1:8080"
)

type Flag struct {
	host *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		f.host = &addr
	}

	return &Config{
		Host: *f.host,
		Key:  []byte("CHANGE ME"),
	}
}
