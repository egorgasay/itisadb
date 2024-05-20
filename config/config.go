package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	TransactionLogger TransactionLoggerConfig `toml:"TransactionLogger"`
	Network           NetworkConfig           `toml:"Network"`
	Encryption        EncryptionConfig        `toml:"Encryption"`
	WebApp            WebAppConfig            `toml:"WebApp"`
	Balancer          BalancerConfig          `toml:"Balancer"`
	Security          SecurityConfig          `toml:"Security"`
	Logging           LoggingConfig           `toml:"Logging"`
}

type TransactionLoggerConfig struct {
	On              bool          `toml:"On"`
	BackupDirectory string        `toml:"BackupDirectory"`
	SyncBufferTime  time.Duration `toml:"SyncBufferTime"`
}

type NetworkConfig struct {
	GRPC string `toml:"GRPC"`
	REST string `toml:"FastHTTP"`
}

type EncryptionConfig struct {
	Key string `toml:"Key"`
}

type WebAppConfig struct {
	On   bool   `toml:"On"`
	Host string `toml:"Host"`
	Logs bool   `toml:"Logs"`
}

type BalancerConfig struct {
	On           bool     `toml:"On"`
	BalancerOnly bool     `toml:"BalancerOnly"`
	Servers      []string `toml:"Servers"`
}

type SecurityConfig struct {
	MandatoryAuthorization bool `toml:"MandatoryAuthorization"`
	MandatoryEncryption    bool `toml:"MandatoryEncryption"`
}

type LoggingConfig struct {
	Level string `toml:"Level"`
}

var _configFlag = flag.String("config", "", "Specify the path to the config file")
var _configServersFlag = flag.String("config-servers", "", "Specify the path to the config file")

var _noSecurity = SecurityConfig{
	MandatoryAuthorization: false,
	MandatoryEncryption:    false,
}

const _defaultPathToConfig = "config/config.toml"
const _defaultPathToServers = "config/config-servers.toml"

func getPathToConfig() string {
	var pathToConfig = _defaultPathToConfig
	if *_configFlag != "" {
		pathToConfig = *_configFlag
	}

	return pathToConfig
}

func getPathToServers() string {
	var pathToServers = _defaultPathToServers
	if *_configServersFlag != "" {
		pathToServers = *_configServersFlag
	}

	return pathToServers
}

func New() (*Config, error) {
	flag.Parse()

	cfg := &Config{}

	_, err := toml.DecodeFile(getPathToServers(), &cfg.Balancer)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	_, err = toml.DecodeFile(getPathToConfig(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}

type ServersConfig struct {
	Servers []string `toml:"Servers"`
}

func UpdateServers(servers []string) error {
	f, err := os.OpenFile(getPathToServers(), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("can't open servers file to insert new")
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(&ServersConfig{Servers: servers}); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	return nil
}
