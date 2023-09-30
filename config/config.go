package config

import (
	"github.com/BurntSushi/toml"
	balancercfg "itisadb/internal/config"
	storagecfg "itisadb/internal/grpc-storage/config"
)

const DefaultConfigPath = "config/default-config.toml"

type Config struct {
	Storage  *storagecfg.Config
	Balancer *balancercfg.Config
}

func New() *Config {
	return &Config{}
}

func (c *Config) FromTOML(f string) error {
	_, err := toml.DecodeFile(f, c)
	if err != nil {
		return err
	}
	return nil
}
