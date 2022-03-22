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
	Service ServiceConfig `yaml:"service"`
	Api     ApiConfig     `yaml:"api"`
}

type ServiceConfig struct {
	Host string `yaml:"host"`
}

type ApiConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
