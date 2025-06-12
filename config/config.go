package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      *Server   `mapstructure:"app"`
		Database *Database `mapstructure:"database"`
	}

	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SslMode  string `mapstructure:"sslmode"`
		TimeZone string `mapstructure:"timezone"`
	}

	Server struct {
		Port        int    `mapstructure:"port"`
		Environment string `mapstructure:"environment"`
		ShortUrl    string `mapstructure:"short_url"`
	}
)

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		config = &Config{} // Initialize the config struct
		err := viper.Unmarshal(config)
		if err != nil {
			log.Fatalf("Error unmarshaling config: %v", err)
		}

		log.Println("Configuration loaded successfully:", config)
	})

	return config
}

// GetString retrieves string value from config
func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}
