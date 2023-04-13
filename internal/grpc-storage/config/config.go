package config

import (
	"flag"
	"os"
)

type DBConfig struct {
	DriverName     string
	DataSourceCred string
}

type Config struct {
	Host        string
	Key         []byte
	DBConfig    *DBConfig
	Balancer    string
	TLoggerType string
	TLoggerDir  string
}

const (
	defaultHost = "127.0.0.1:8080"
)

type Flag struct {
	host        *string
	dsn         *string
	balancer    *string
	tloggerType *string
	tLoggerDir  *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	f.dsn = flag.String("d", "", "-d=connection_string")
	f.balancer = flag.String("connect", "", "-connect=ip:port")
	f.tloggerType = flag.String("tlog_type", "db", "-tlog_type=db")
	f.tLoggerDir = flag.String("tlog_dir", "/", "-tlog_dir=/tmp")
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		f.host = &addr
	}

	if dsn, ok := os.LookupEnv("DATABASE_URI"); ok {
		f.dsn = &dsn
	}

	return &Config{
		Host: *f.host,
		Key:  []byte("CHANGE ME"),
		DBConfig: &DBConfig{
			DriverName:     "mongo",
			DataSourceCred: *f.dsn,
		},
		Balancer:    *f.balancer,
		TLoggerType: *f.tloggerType,
		TLoggerDir:  *f.tLoggerDir,
	}
}
