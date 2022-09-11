package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadFromFile(dir, file string, cfg interface{}) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}
