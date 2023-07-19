package config

import (
	"flag"
	"log"
)

type Config struct {
	TgBotToken string
}

func MustLoad() Config {
	tgBotToken := flag.String(
		"tg-bot-token",
		"",
		"token for access to Telegram bot",
	)

	flag.Parse()

	if *tgBotToken == "" {
		log.Fatal("Telegram token is not specified")
	}

	return Config{
		TgBotToken: *tgBotToken,
	}
}
