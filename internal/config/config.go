package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type InfluxDBConfig struct {
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

type Config struct {
	ListenAddress  string         `yaml:"listen_address"`
	AgentAuthToken string         `yaml:"agent_auth_token"`
	InfluxDB       InfluxDBConfig `yaml:"influxdb"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
