package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	portMin = 1024
	portMax = 65535
)

type Config struct {
	Port string
}

type validateFunc func(string) error

// Get returns a configuration struct, or error if it can't.
func Get() (Config, error) {
	viper.AutomaticEnv()

	validations := map[string]validateFunc{
		"PORT": isPort,
	}

	for k, v := range validations {
		if err := v(k); err != nil {
			return Config{}, fmt.Errorf("config.Get: %w", err)
		}
	}

	return Config{
		Port: viper.GetString("PORT"),
	}, nil
}

func isPort(key string) error {
	n := viper.GetInt(key)
	if n < portMin || n > portMax {
		return fmt.Errorf("configured port is not in permissible range of %d - %d, got %d", portMin, portMax, n)
	}

	return nil
}
