package config

import (
	"flag"
	"reflect"
)

type DBConfig struct {
	DriverName     string
	DataSourceCred string
}

type Config struct {
	Host      string
	Balancer  string `toml:"Balancer"`
	IsTLogger bool   `toml:"TransactionLoggerDir"`
}

const (
	defaultHost = "127.0.0.1:8080"
)

type Flag struct {
	host      *string
	dsn       *string
	balancer  *string
	isTLogger *bool
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
}

func New(cfg *Config) *Config {
	flag.Parse()

	return &Config{
		Host:      chooseLeftIfSet[string](f.host, &cfg.Host),
		Balancer:  chooseLeftIfSet[string](f.balancer, &cfg.Balancer),
		IsTLogger: chooseLeftIfSet[bool](f.isTLogger, &cfg.IsTLogger),
	}
}

func chooseLeftIfSet[C any](l, r *C) C {
	left := reflect.ValueOf(l)
	right := reflect.ValueOf(r)

	if left.IsNil() || left.Elem().IsZero() {
		return right.Elem().Interface().(C)
	}
	return left.Elem().Interface().(C)
}
