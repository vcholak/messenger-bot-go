package main

import (
	"log"

	tgClient "github.com/vcholak/messenger-bot/clients/telegram"
	"github.com/vcholak/messenger-bot/config"
	event_consumer "github.com/vcholak/messenger-bot/consumer/event-consumer"
	"github.com/vcholak/messenger-bot/events/telegram"
	"github.com/vcholak/messenger-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = ".file-storage"
	batchSize   = 100
)

func main() {
	cfg := config.MustLoad()

	storage := files.New(storagePath)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken),
		storage,
	)

	log.Print("Bot service is started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("Bot service is stopped", err)
	}
}
