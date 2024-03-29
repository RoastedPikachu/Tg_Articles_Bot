package main

import (
	"context"
	"flag"
	"log"
	"os"
	tgClient "read-adviser-bot/clients/telegram"
	"read-adviser-bot/consumer/eventConsumer"
	"read-adviser-bot/events/telegram"
	"read-adviser-bot/storage/sqlite"
)

func main() {
	os.Setenv("TG_BOT_HOST", "api.telegram.org")
	os.Setenv("GET_UPDATES_METHOD", "getUpdates")
	os.Setenv("SEND_MESSAGE_METHOD", "sendMessage")

	//s := files.New("filesStorage")
	s, err := sqlite.New("data/sqlite/storage.db")
	if err != nil {
		log.Fatalf("Не смог подключится к хранилищу: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("Не смог инициализировать хранилище: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(os.Getenv("TG_BOT_HOST"), mustToken()),
		s,
	)

	log.Print("Cервис запущен")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal("Сервис остановлен", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "Токен для доступа от телеграм бота")

	flag.Parse()

	if *token == "" {
		log.Fatalln("token is not specified")
	}

	return *token
}
