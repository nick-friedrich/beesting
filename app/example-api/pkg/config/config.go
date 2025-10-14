package config

import "sync"

type Config struct {
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
