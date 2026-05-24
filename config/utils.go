package config

import (
	"context"
	"log"
	"os"
)

func LoadConfig(ctx context.Context, environmentFlag string) *Config {
	configPath := getConfigPath(environmentFlag)
	if configPath == "" {
		log.Printf("Environment type not supported, value {%v}", environmentFlag)
		os.Exit(1)
	}

	conf, err := FromPath(ctx, configPath)
	if err != nil {
		log.Printf("failed to load configs: %v", err)
		os.Exit(1)
	}

	return conf
}

func getConfigPath(environmentFlag string) string {
	switch Environment(environmentFlag) {
	case Production:
		return "./config-production.yml"
	case Staging:
		return "./config-staging.yml"
	default:
		return "./config.yml" // unreachable code
	}
}
