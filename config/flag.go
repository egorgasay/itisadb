package config

import "flag"

type Flag struct {
	grpc       *string
	rest       *string
	tlog       *string
	webAppHost *string
}

var f Flag

const (
	defaultGRPC       = ":8888"
	defaultTLogger    = ""
	defaultRESTHost   = ""
	defaultWebAppHost = ":6070"
)

func init() {
	f.grpc = flag.String("grpc", defaultGRPC, "-grpc=host")
	f.rest = flag.String("rest", defaultRESTHost, "-rest=host")
	f.webAppHost = flag.String("a", defaultWebAppHost, "-a=host")
	f.tlog = flag.String("tlog", defaultTLogger, "-tlog=dir")
}
