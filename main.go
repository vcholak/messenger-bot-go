package main

import (
	"context"
	"flag"
	"log"
	"time"

	tgClient "github.com/vcholak/messenger-bot/clients/telegram"
	event_consumer "github.com/vcholak/messenger-bot/consumer/event-consumer"
	"github.com/vcholak/messenger-bot/events/telegram"

	// "github.com/vcholak/messenger-bot/storage/files"
	"github.com/vcholak/messenger-bot/storage/sqlite"
)

const (
	tgBotHost = "api.telegram.org"
	// storagePath = ".file-storage"
	storagePath = ".sqlite/storage.db"
	batchSize   = 100
)

func main() {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// storage := files.New(storagePath)
	storage, err := sqlite.New(storagePath, cancelFunc)
	if err != nil {
		log.Fatal("Can't connect to SQlite storage: ", err)
	}

	if err := storage.Init(ctx); err != nil {
		log.Fatal("Can't init the storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		storage,
	)

	log.Print("Bot service is started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("Bot service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to Telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Telegram token is not specified")
	}

	return *token
}
