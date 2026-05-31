package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	App struct {
		Port  int  `yaml:"port"`
		Debug bool `yaml:"debug"`
	} `yaml:"app"`

	SQLite struct {
		Path string `yaml:"path"`
	} `yaml:"sqlite"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
