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
	Host        string
	Balancer    string `toml:"Balancer"`
	TLoggerDir  string `toml:"TransactionLoggerDir"`
	TLoggerDSN  string `toml:"TransactionLoggerDSN"`
	DSN         string `toml:"DSN"`
	TLoggerType string `toml:"TransactionLogger"`
}

const (
	defaultHost = "127.0.0.1:8080"
)

type Flag struct {
	host        *string
	dsn         *string
	balancer    *string
	tLoggerDir  *string
	tLoggerDSN  *string
	tLoggerType *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	//f.dsn = flag.String("d", "", "-d=mongodb_connection_string")
	//f.balancer = flag.String("connect", "", "-connect=ip:port")
	//f.tLoggerDir = flag.String("tlog_dir", "/", "-tlog_dir=/tmp")
	//f.tLoggerDSN = flag.String("tlog_dsn", "", "-tlog_dsn=mysql_connection_string")
	//f.tLoggerType = flag.String("tlog_type", "Off", "-tlog_type=docker_db || db || file")
}

func New(cfg *Config) *Config {
	flag.Parse()

	return &Config{
		Host:     chooseLeftIfSet[string](f.host, &cfg.Host),
		Balancer: chooseLeftIfSet[string](f.balancer, &cfg.Balancer),
		DSN:      chooseLeftIfSet[string](f.dsn, &cfg.DSN),

		TLoggerDir:  chooseLeftIfSet[string](f.tLoggerDir, &cfg.TLoggerDir),
		TLoggerType: chooseLeftIfSet[string](f.tLoggerType, &cfg.TLoggerType),
		TLoggerDSN:  chooseLeftIfSet[string](f.tLoggerDSN, &cfg.TLoggerDSN),
	}
}

func chooseLeftIfSet[C any](l, r *C) C {
	left := reflect.ValueOf(l)
	right := reflect.ValueOf(r)

	if left.IsNil() || left.IsZero() {
		return right.Elem().Interface().(C)
	}
	return left.Elem().Interface().(C)
}
