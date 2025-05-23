package config

import "time"

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	Storage  string        `yaml:"storage" end-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRCP     GRCPConfig
}

type GRCPConfig struct {
	GRPCPort int           `yaml:"grpc_port"`
	Timeout  time.Duration `yaml:"timeout"`
}
