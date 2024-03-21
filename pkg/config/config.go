package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigFile = "config.yaml"

type ApiKey struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type Config struct {
	ApiKeys []ApiKey `yaml:"api_keys"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}

	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read yaml file: #%v", err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml file: #%v", err)
	}

	return cfg, nil
}
