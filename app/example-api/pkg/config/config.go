package config

import "sync"

type Config struct {
	BaseURL    string
	AuthConfig AuthConfig
}

type AuthConfig struct {
	ConfirmEmail bool
}

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		configInstance = &Config{
			BaseURL: "http://localhost:3000",
			AuthConfig: AuthConfig{
				ConfirmEmail: true, // Default value
			},
		}
	})
	return configInstance
}

func InitConfig(config *Config) {
	once.Do(func() {
		configInstance = config
	})
}
