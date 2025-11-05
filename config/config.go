package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppName           string        `mapstructure:"APP_NAME"`
	AppEnv            string        `mapstructure:"APP_ENV"`
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	ServerPort        int           `mapstructure:"SERVER_PORT"`
	BasicAuthUsername string        `mapstructure:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword string        `mapstructure:"BASIC_AUTH_PASSWORD"`
	RequestTimeout    time.Duration `mapstructure:"REQUEST_TIMEOUT"`
}

func LoadConfig(configFilePath string) (*AppConfig, error) {
	// Configure Viper to read the specified ..env.yaml file
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("env")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error reading config file: %s", err)
	}

	// Optionally, automatically override values with environment variables if they exist
	viper.AutomaticEnv()

	// Create an instance of the AppConfig struct
	var config AppConfig

	// Unmarshal the configuration values into the struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %s", err)
	}

	return &config, nil
}
