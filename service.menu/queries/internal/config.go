package internal

import (
	"github.com/spf13/viper"
)

type Config struct {
	EventStoreConnectionString string `mapstructure:"EVENT_STORE_CONNECTION_STRING"`
	PostgresConnectionString   string `mapstructure:"POSTGRES_CONNECTION_STRING"`
	ResourcePath               string `mapstructure:"RESOURCE_PATH"`
	ResourceHost               string `mapstructure:"RESOURCE_HOST"`
}

func LoadConfig(path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	return
}
