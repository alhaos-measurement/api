package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address string `yaml:"address"`
	DB      DB     `yaml:"db"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// New read config from file and return config instance
func New(filename string) (*Config, error) {
	config := Config{}
	err := cleanenv.ReadConfig(filename, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
