package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	TransactionLoggerConfig TransactionLoggerConfig `toml:"TransactionLogger"`
	NetworkConfig           NetworkConfig           `toml:"Network"`
	EncryptionConfig        EncryptionConfig        `toml:"Encryption"`
	WebAppConfig            WebAppConfig            `toml:"WebApp"`
	Balancer                BalancerConfig          `toml:"Balancer"`
}

type TransactionLoggerConfig struct {
	On              bool   `toml:"On"`
	BackupDirectory string `toml:"BackupDirectory"`
}

type NetworkConfig struct {
	GRPC string `toml:"GRPC"`
	REST string `toml:"FastHTTP"`
}

type EncryptionConfig struct {
	On  bool   `toml:"On"`
	Key string `toml:"Key"`
}

type WebAppConfig struct {
	On   bool   `toml:"On"`
	Host string `toml:"Host"`
}

type BalancerConfig struct {
	On      bool     `toml:"On"`
	Servers []string `toml:"Servers"`
}

func New() (*Config, error) {
	flag.Parse()

	cfg := &Config{}
	_, err := toml.DecodeFile("config/default-config.toml", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}
