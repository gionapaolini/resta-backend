package internal

import "github.com/spf13/viper"

type Config struct {
	EventStoreConnectionString string `mapstructure:"EVENT_STORE_CONNECTION_STRING"`
	ResourcePath               string `mapstructure:"RESOURCE_PATH"`
}

func LoadConfig(path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	return
}
