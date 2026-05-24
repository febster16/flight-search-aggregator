package config

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP HTTPConfig `mapstructure:"http"`
}

func FromPath(ctx context.Context, path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config file path is empty")
	}

	log.Printf("Reading config: %s", path)

	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	viper.SetConfigType("yaml")

	buffer := bytes.NewBuffer(bts)

	if readErr := viper.ReadConfig(buffer); readErr != nil {
		return nil, fmt.Errorf("read config file data: %w", readErr)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	log.Printf("Finished reading config: %s", path)

	return &config, nil
}
