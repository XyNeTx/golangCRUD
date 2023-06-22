// configs/configs.go
package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	APIURL        string `mapstructure:"API_URL"`
	LineToken     string `mapstructure:"API_LINE_TOKEN"`
}

var AppConfig Config

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return err
	}

	return nil
}

// InitConfig initializes the global config variable.
func InitConfig(path string) {
	err := LoadConfig(path)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
}
