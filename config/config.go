package config

import (
	"flag"
	"log"
)

type Config struct {
	TgBotToken        string
	MongodbConnection string
}

func MustLoad() Config {
	tgBotToken := flag.String(
		"tg-bot-token",
		"",
		"token for access to Telegram bot",
	)

	mongodbConnection := flag.String(
		"mongodb-connection",
		"",
		"connection string for MongoDB",
	)

	flag.Parse()

	if *tgBotToken == "" {
		log.Fatal("Telegram token is not specified")
	}

	return Config{
		TgBotToken:        *tgBotToken,
		MongodbConnection: *mongodbConnection,
	}
}
