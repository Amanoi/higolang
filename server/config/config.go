package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Scheduler SchedulerConfig `mapstructure:"scheduler"`
	Seed      SeedConfig      `mapstructure:"seed"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig holds database settings.
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // hours
}

// SchedulerConfig holds feed scheduler settings.
type SchedulerConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	DefaultInterval  int  `mapstructure:"default_interval"` // minutes
}

// SeedConfig holds seed data settings.
type SeedConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

var Cfg Config

// Load reads configuration from config.yaml and environment variables.
func Load() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return Cfg
}
