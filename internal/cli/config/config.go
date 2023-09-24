package config

import (
	"flag"
)

type Config struct {
	Balancer string
	Host     string
}

const (
	defaultBalancer = "127.0.0.1:8888"
	defaultHost     = "127.0.0.1:8087"
)

type Flag struct {
	balancer *string
	host     *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	f.balancer = flag.String("b", defaultBalancer, "-b=host")
}

func New() *Config {
	flag.Parse()

	return &Config{
		Balancer: *f.balancer,
		Host:     *f.host,
	}
}
