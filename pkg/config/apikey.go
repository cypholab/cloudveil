package config

type ApiKey struct {
	Name string `yaml:"name"`
	Auth any    `yaml:"auth"`
}
