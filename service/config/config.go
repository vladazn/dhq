package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func New(path string) (*Config, error) {
	conf := Config{}
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	err = yaml.Unmarshal(dat, &conf)

	return &conf, err
}

type Config struct {
	Redis   RedisConfigs   `yaml:"redis"`
	Grpc    GrpcConfigs    `yaml:"grpc"`
	MariaDb MariaDBConfigs `yaml:"mariadb"`
}

type MariaDBConfigs struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Db   string `yaml:"db"`
}

type RedisConfigs struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Db   int    `yaml:"db"`
}

type GrpcConfigs struct {
	Host string
	Port int
}
