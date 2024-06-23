package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Host        string `mapstructure:"HOST"`
	Port        string `mapstructure:"PORT"`

	DBUsername    string `mapstructure:"DATABASE_USER"`
	DBPassword    string `mapstructure:"DATABASE_PASSWORD"`
	DBHost        string `mapstructure:"DATABASE_HOST"`
	DBPort        string `mapstructure:"DATABASE_PORT"`
	DBName        string `mapstructure:"DATABASE_NAME"`
	MigrationPath string `mapstructure:"MIGRATION_PATH"`
	DBRecreate    bool   `mapstructure:"DATABASE_RECREATE"`
}

func LoadConfig(name string, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	return
}
