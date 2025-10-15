package config

import "sync"

type Config struct {
	BaseURL     string
	EmailConfig EmailConfig
	AuthConfig  AuthConfig
}

type AuthConfig struct {
	ConfirmEmail bool
}

type EmailConfig struct {
	From string
	Name string
}

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		configInstance = &Config{
			BaseURL: "http://localhost:3000",
			EmailConfig: EmailConfig{
				From: "noreply@example.com",
				Name: "Example",
			},
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
