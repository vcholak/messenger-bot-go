package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	TgBotToken        string `mapstructure:"TG_BOT_TOKEN"`
	MongodbConnection string `mapstructure:"MONGO_DB_CONNECTION"`
}

// MustLoad reads config from ./app.env file.
func MustLoad() (cfg *Config) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file: ", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
	return
}
