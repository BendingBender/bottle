package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	keyPort = "PORT"
	keyHost = "HOST"
	portMin = 1024
	portMax = 65535
)

type Config struct {
	Port string
	Host string
}

type validateFunc func(string) error

// Get returns a configuration struct, or error if it can't.
func Get() (Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault(keyPort, "7889")
	viper.SetDefault(keyHost, "go.bottle.remotehack.space")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("reading config failed: %s", err)
		os.Exit(1)
	}

	validations := map[string]validateFunc{
		keyPort: isPort,
		keyHost: isNotEmpty,
	}

	for k, v := range validations {
		if err := v(k); err != nil {
			return Config{}, fmt.Errorf("config.Get: %w", err)
		}
	}

	return Config{
		Port: viper.GetString(keyPort),
		Host: viper.GetString(keyHost),
	}, nil
}

func isPort(key string) error {
	n := viper.GetInt(key)
	if n < portMin || n > portMax {
		return fmt.Errorf("configured port is not in permissible range of %d - %d, got %d", portMin, portMax, n)
	}

	return nil
}

func isNotEmpty(key string) error {
	h := viper.GetString(key)
	if h == "" {
		return errors.New("configured host is empty")
	}

	return nil
}
