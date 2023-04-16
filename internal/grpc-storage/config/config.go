package config

import (
	"flag"
)

type DBConfig struct {
	DriverName     string
	DataSourceCred string
}

type Config struct {
	Host       string
	Key        []byte
	DBConfig   *DBConfig
	Balancer   string
	TLoggerDir string
	DSN        string
}

const (
	defaultHost = "127.0.0.1:8080"
)

type Flag struct {
	host       *string
	dsn        *string
	balancer   *string
	tLoggerDir *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	f.dsn = flag.String("d", "", "-d=mysql_connection_string")
	f.balancer = flag.String("connect", "", "-connect=ip:port")
	f.tLoggerDir = flag.String("tlog_dir", "/", "-tlog_dir=/tmp")
}

func New() *Config {
	flag.Parse()

	return &Config{
		Host: *f.host,
		Key:  []byte("CHANGE ME"),
		DBConfig: &DBConfig{
			DriverName:     "mongo",
			DataSourceCred: *f.dsn,
		},
		Balancer:   *f.balancer,
		DSN:        *f.dsn,
		TLoggerDir: *f.tLoggerDir,
	}
}
