package config

type HTTPConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}
