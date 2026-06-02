package config

import (
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	App struct {
		Port  int  `yaml:"port"`
		Debug bool `yaml:"debug"`
	} `yaml:"app"`

	Server struct {
		ReadTimeout       time.Duration `yaml:"read_timeout"`
		WriteTimeout      time.Duration `yaml:"write_timeout"`
		IdleTimeout       time.Duration `yaml:"idle_timeout"`
		ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
		MaxHeaderBytes    int           `yaml:"max_header_bytes"`
	} `yaml:"server"`

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
